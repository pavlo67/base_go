package files_sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/pavlo67/base_go/entities/files"
	"github.com/pavlo67/base_go/lib/db"
	"github.com/pavlo67/base_go/lib/errors"
	"github.com/pavlo67/base_go/lib/logger"
)

var _ files.Operator = &filesSQLite{}
var _ db.Operator = &filesSQLite{}

type filesSQLite struct {
	db *sql.DB
}

var l logger.Operator

const onNew = "on files_sqlite.Init():"

func New(dsn string, l_ logger.Operator) (files.Operator, db.Operator, error) {
	if l_ == nil {
		return nil, nil, errors.New("", "l_ == nil")
	}
	l = l_

	sqlDB, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, nil, errors.Wrap(err, onNew)
	}

	op := &filesSQLite{db: sqlDB}

	return op, op, nil
}

// db.Operator --------------------------------------------------------------------------------

const onCreate = "on files_sqlite.Create():"

func (op *filesSQLite) Create() error {
	env := strings.ToUpper(os.Getenv("ENV"))
	if env != "TEST" {
		return fmt.Errorf("filesSQLite.Clean() is allowed only when ENV=TEST, but env = %s", env)
	}

	_, err := op.db.Exec(`
		CREATE TABLE IF NOT EXISTS files (
			path        TEXT PRIMARY KEY,
			is_dir      INTEGER NOT NULL,
			size        INTEGER NOT NULL,
			ctime       TEXT,
			mtime       TEXT,
			mime_type   TEXT NOT NULL,
			created_at  TEXT NOT NULL,
			updated_at  TEXT NOT NULL
		)
	`)

	if err != nil {
		errClose := sqlDB.Close()
		if errClose != nil {
			l.Errorf("%s: %s ", onNew, errClose.Error())
		}
		return nil, nil, errors.Wrap(err, onNew)
	}
}

const onClean = "on files_sqlite.Clean():"

func (op *filesSQLite) Clean() error {
	env := strings.ToUpper(os.Getenv("ENV"))
	if env != "TEST" {
		return fmt.Errorf("filesSQLite.Clean() is allowed only when ENV=TEST, but env = %s", env)
	}

	_, err := op.db.Exec(`DELETE FROM files`)
	return errors.Wrap(err, onClean)
}

// files.Operator -----------------------------------------------------------------------------

const onSave = "on files_sqlite.Save():"

func (op *filesSQLite) Save(file files.File) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)

	_, err := op.db.Exec(`
		INSERT INTO files (
			path, is_dir, size, ctime, mtime, mime_type, created_at, updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(path) DO UPDATE SET
			is_dir     = excluded.is_dir,
			size       = excluded.size,
			сtime      = excluded.сtime,
			mtime      = excluded.mtime,
			mime_type  = excluded.mime_type,
			updated_at = excluded.updated_at
	`,
		file.Path,
		boolToInt(file.IsDir),
		file.Size,
		timeToDB(file.CTime),
		timeToDB(file.MTime),
		file.MimeType,
		now,
		now,
	)

	return errors.Wrap(err, onSave)
}

const onRead = "on files_sqlite.Read():"

func (op *filesSQLite) Read(path string) (*files.Item, error) {
	row := op.db.QueryRow(`
		SELECT path, is_dir, size, ctime, mtime, mime_type, created_at, updated_at
		FROM files
		WHERE path = ?
	`, path)

	item, err := scanItem(row)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors.Wrap(err, onRead)
	}

	return item, nil
}

const onRemove = "on files_sqlite.Remove():"

func (op *filesSQLite) Remove(path string) error {
	_, err := op.db.Exec(`DELETE FROM files WHERE path = ?`, path)
	return errors.Wrap(err, onRemove)
}

const onList = "on files_sqlite.List():"

func (op *filesSQLite) List(path string, depth int) ([]files.Item, error) {
	rows, err := op.db.Query(`
		SELECT path, is_dir, size, ctime, mtime, mime_type, created_at, updated_at
		FROM files
		WHERE path LIKE ?
		ORDER BY path
	`, strings.TrimRight(path, "/")+`/%`)
	if err != nil {
		return nil, errors.Wrap(err, onList)
	}
	defer rows.Close()

	var items []files.Item

	for rows.Next() {
		item, err := scanItem(rows)
		if err != nil {
			return nil, errors.Wrap(err, onList)
		}

		if filepath.Dir(item.Path) == path {
			items = append(items, *item)
		}
	}

	return items, errors.Wrap(rows.Err(), onList)
}
