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
		return fmt.Errorf("filesSQLite.Clean() is allowed only when ENV=TEST, but env = %s", env)
	}

	err := sqllib.InitDB(db, filepath.Join(filelib.CurrentPath(), "create.sql"))
	return errors.Wrap(err, onCreate)
}

const sqlClean = "DELETE FROM files"

const onClean = "on files_sqlite.Clean():"

func (op *filesSQLite) Clean() error {
	env := strings.ToUpper(os.Getenv("ENV"))
	if env != "TEST" {
		return fmt.Errorf("filesSQLite.Clean() is allowed only when ENV=TEST, but env = %s", env)
	}

	_, err := op.stmClean.Exec()
	return errors.Wrap(err, onClean)
}
