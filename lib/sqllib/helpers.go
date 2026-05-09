package sqllib

import (
	"database/sql"
	"github.com/pavlo67/base_go/lib/errors"
)

type SqlStmt struct {
	Stmt **sql.Stmt
	Sql  string
}

func Prepare(dbh *sql.DB, sqlQuery string, stmt **sql.Stmt) error {
	var err error

	*stmt, err = dbh.Prepare(sqlQuery)
	if err != nil {
		return errors.Wrapf(err, "can't dbh.Prepare(%s)", sqlQuery)
	}

	return nil
}
