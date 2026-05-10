package sqllib

import (
	"database/sql"
	"os"

	"github.com/pavlo67/base_go/lib/errors"
)

const onInitDB = "on sqllib.InitDB()"

func InitDB(db *sql.DB, initStmntPath string) error {
	if db == nil {
		return errors.New("", onInitDB+": db == nil")
	}

	sqlBytes, err := os.ReadFile(initStmntPath)
	if err != nil {
		return errors.Newf("", onInitDB+": os.ReadFile(%s) error: %v", initStmntPath, err)
	}

	if err := db.Ping(); err != nil {
		return errors.Newf("", onInitDB+": db.Ping() error: %v", err)
	}

	if _, err := db.Exec(string(sqlBytes)); err != nil {
		return errors.Newf("", onInitDB+": db.Exec(%s) error: %v", string(sqlBytes), err)
	}

	return nil
}

type SqlStmt struct {
	Stmt **sql.Stmt
	Sql  string
}

func PrepareQuery(dbh *sql.DB, sqlQuery string, stmt **sql.Stmt) error {
	var err error

	*stmt, err = dbh.Prepare(sqlQuery)
	if err != nil {
		return errors.Wrapf(err, "can't dbh.Prepare(%s)", sqlQuery)
	}

	return nil
}
