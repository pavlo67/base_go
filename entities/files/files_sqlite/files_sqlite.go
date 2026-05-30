package files_sqlite

import (
	"database/sql"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/pavlo67/base_go/lib/db"
	"github.com/pavlo67/base_go/lib/errors"
	"github.com/pavlo67/base_go/lib/logger"
	"github.com/pavlo67/base_go/lib/sqllib"

	"github.com/pavlo67/base_go/entities/files"
)

var _ files.Operator = &filesSQLite{}
var _ db.Operator = &filesSQLite{}

type filesSQLite struct {
	db                                                                 *sql.DB
	stmSave, stmRead, stmRemove, stmRemoveRecursive, stmList, stmClean *sql.Stmt
}

var l logger.Operator

const onNew = "on files_sqlite.Init():"

func New(dsn string, create bool, l_ logger.Operator) (files.Operator, db.Operator, error) {
	if l_ == nil {
		return nil, nil, errors.New("", onNew+" l_ == nil")
	}
	l = l_

	sqlDB, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, nil, errors.Wrap(err, onNew)
	}

	op := &filesSQLite{db: sqlDB}

	if create {
		if err = op.Create(sqlDB); err != nil {
			return nil, nil, errors.Wrap(err, onNew)
		}
	}

	sqlStmts := []sqllib.SqlStmt{
		{&op.stmSave, sqlSave},
		{&op.stmRead, sqlRead},
		{&op.stmRemove, sqlRemove},
		{&op.stmRemoveRecursive, sqlRemoveRecursive},
		{&op.stmList, sqlList},
		{&op.stmClean, sqlClean},
	}

	for _, sqlStmt := range sqlStmts {
		if err := sqllib.PrepareQuery(sqlDB, sqlStmt.Sql, sqlStmt.Stmt); err != nil {
			return nil, nil, errors.Wrap(err, onNew)
		}
	}

	return op, op, nil
}

// files.Operator -----------------------------------------------------------------------------

const sqlSave = `
	INSERT INTO files (
		path, is_dir, size, ctime, mtime, crc, mime_type, created_at, updated_at
	)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	ON CONFLICT(path) DO UPDATE SET
		is_dir     = excluded.is_dir,
		size       = excluded.size,
		ctime      = excluded.ctime,
		mtime      = excluded.mtime,
		crc        = excluded.crc,
		mime_type  = excluded.mime_type,
		updated_at = excluded.updated_at
`
const onSave = "on files_sqlite.Save():"

func (op *filesSQLite) Save(data files.Data) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)

	_, err := op.stmSave.Exec(
		data.Path,
		boolToInt(data.IsDir),
		data.Size,
		timeToDB(data.CTime),
		timeToDB(data.MTime),
		data.CRC,
		data.MimeType,
		now,
		now,
	)

	return errors.Wrap(err, onSave)
}

const sqlRead = `
	SELECT path, is_dir, size, ctime, mtime, crc, mime_type, created_at, updated_at
	FROM files
	WHERE path = ?
`
const onRead = "on files_sqlite.Read():"

func (op *filesSQLite) Read(path string) (*files.Item, error) {
	row := op.stmRead.QueryRow(path)

	item, err := scanItem(row)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, onRead)
	}

	return item, nil
}

const sqlRemove = `DELETE FROM files WHERE path = ?`
const sqlRemoveRecursive = `DELETE FROM files WHERE path = ? OR path LIKE ?`

const onRemove = "on files_sqlite.Remove():"

func (op *filesSQLite) Remove(path string, forceRecursion bool) error {
	if !forceRecursion {
		_, err := op.stmRemove.Exec(path)
		return errors.Wrap(err, onRemove)
	}

	_, err := op.stmRemoveRecursive.Exec(path, strings.TrimRight(path, "/")+"/%")
	return errors.Wrap(err, onRemove)
}

const sqlList = `
	SELECT path, is_dir, size, ctime, mtime, crc, mime_type, created_at, updated_at
	FROM files
	WHERE path LIKE ?
	ORDER BY path
`

const onList = "on files_sqlite.List():"

func (op *filesSQLite) List(path string, depth int) ([]files.Item, error) {
	if depth < 0 {
		return nil, errors.New("", onList+" depth must be >= 0")
	}

	path = strings.TrimRight(path, "/")
	rows, err := op.stmList.Query(strings.TrimRight(path, "/") + `/%`)
	if err != nil {
		return nil, errors.Wrap(err, onList)
	}
	defer func() {
		if err = rows.Close(); err != nil {
			l.Errorf(onList+": on rows.Close(): %v", err)
		}
	}()

	var items []files.Item

	for rows.Next() {
		item, err := scanItem(rows)
		if err != nil {
			return nil, errors.Wrap(err, onList)
		}

		if depth == 0 || itemDepth(path, item.Path) <= depth {
			items = append(items, *item)
		}
	}

	return items, errors.Wrap(rows.Err(), onList)
}

func itemDepth(basePath, itemPath string) int {
	rel := strings.TrimPrefix(strings.TrimPrefix(itemPath, basePath), "/")
	if rel == "" {
		return 0
	}
	return len(strings.Split(filepath.ToSlash(rel), "/"))
}
