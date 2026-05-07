package files_sqlite

import (
	"database/sql"
	"fmt"
	errors_new "github.com/pavlo67/data/add_new/errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/data/entities/files"
)

var _ files.Operator = &filesSQLite{}

type filesSQLite struct {
	db *sql.DB
}

var l logger.Operator

const onNew = "on files_sqlite.Init():"

func New(dsn string, l_ logger.Operator) (files.Operator, db.Cleaner, error) {
	if l_ == nil {
		return nil, nil, errors.New("l_ == nil")
	}
	l = l_

	sqlDB, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, nil, errors_new.Wrap(err, onNew)
	}

	op := &filesSQLite{db: sqlDB}

	_, err = op.db.Exec(`
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
		return nil, nil, errors_new.Wrap(err, onNew)
	}

	return op, op, nil
}

const onClean = "on files_sqlite.Clean():"

func (f *filesSQLite) Clean() error {
	env := strings.ToUpper(os.Getenv("ENV"))
	if env != "TEST" {
		return fmt.Errorf("filesSQLite.Clean() is allowed only when ENV=TEST, but env = %s", env)
	}

	_, err := f.db.Exec(`DELETE FROM files`)
	return errors_new.Wrap(err, onClean)
}

const onSave = "on files_sqlite.Save():"

func (f *filesSQLite) Save(file files.File) error {
	now := time.Now().UTC().Format(time.RFC3339Nano)

	_, err := f.db.Exec(`
		INSERT INTO files (
			path, is_dir, size, ctime, mtime, mime_type, created_at, updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(path) DO UPDATE SET
			is_dir     = excluded.is_dir,
			size       = excluded.size,
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

	return errors_new.Wrap(err, onSave)
}

const onRead = "on files_sqlite.Read():"

func (f *filesSQLite) Read(path string) (*files.Item, error) {
	row := f.db.QueryRow(`
		SELECT path, is_dir, size, ctime, mtime, mime_type, created_at, updated_at
		FROM files
		WHERE path = ?
	`, path)

	item, err := scanItem(row)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, errors_new.Wrap(err, onRead)
	}

	return item, nil
}

const onRemove = "on files_sqlite.Remove():"

func (f *filesSQLite) Remove(path string) error {
	_, err := f.db.Exec(`DELETE FROM files WHERE path = ?`, path)
	return errors_new.Wrap(err, onRemove)
}

const onList = "on files_sqlite.List():"

func (f *filesSQLite) List(path string, depth int) ([]files.Item, error) {
	rows, err := f.db.Query(`
		SELECT path, is_dir, size, ctime, mtime, mime_type, created_at, updated_at
		FROM files
		WHERE path LIKE ?
		ORDER BY path
	`, strings.TrimRight(path, "/")+`/%`)
	if err != nil {
		return nil, errors_new.Wrap(err, onList)
	}
	defer rows.Close()

	var items []files.Item

	for rows.Next() {
		item, err := scanItem(rows)
		if err != nil {
			return nil, errors_new.Wrap(err, onList)
		}

		if filepath.Dir(item.Path) == path {
			items = append(items, *item)
		}
	}

	return items, errors_new.Wrap(rows.Err(), onList)
}
