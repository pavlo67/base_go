package files_sqlite

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pavlo67/base_go/lib/errors"
	"github.com/pavlo67/base_go/lib/filelib"
	"github.com/pavlo67/base_go/lib/sqllib"
)

// db.Operator --------------------------------------------------------------------------------

const onCreate = "on files_sqlite.Create():"

func (op *filesSQLite) Create(db *sql.DB) error {
	env := strings.ToUpper(os.Getenv("ENV"))
	if env != "TEST" {
		return fmt.Errorf(onCreate+" allowed only when ENV=TEST, but env = %s", env)
	}

	err := sqllib.InitDB(db, filepath.Join(filelib.CurrentPath(), "create.sql"))
	return errors.Wrap(err, onCreate)
}

const sqlClean = "DELETE FROM files WHERE path = ? OR path LIKE ?"

const onClean = "on files_sqlite.Clean():"

func (op *filesSQLite) Clean(opts interface{}) error {
	env := strings.ToUpper(os.Getenv("ENV"))
	if env != "TEST" {
		return fmt.Errorf(onClean+" allowed only when ENV=TEST, but env = %s", env)
	}

	path, ok := opts.(string)
	if !ok {
		return fmt.Errorf(onClean+" expects directory path as string, got %T", opts)
	}
	if path == "" {
		return fmt.Errorf(onClean + " expects non-empty directory path")
	}

	path = strings.TrimRight(path, "/")
	pathLike := path + "/%"
	if path == "" {
		path = "/"
		pathLike = "/%"
	}

	_, err := op.stmClean.Exec(path, pathLike)
	return errors.Wrap(err, onClean)
}
