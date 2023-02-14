package records_sqlite

import (
	"database/sql"
	"fmt"

	"github.com/pavlo67/common/common/db/db_sqlite"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/data/entities/records"
)

func Starter() starter.Operator {
	return &recordsSQLiteStarter{}
}

var l logger.Operator
var _ starter.Operator = &recordsSQLiteStarter{}

type recordsSQLiteStarter struct {
	table string

	dbKey joiner.InterfaceKey

	interfaceKey joiner.InterfaceKey
	cleanerKey   joiner.InterfaceKey
}

func (rss *recordsSQLiteStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (rss *recordsSQLiteStarter) Prepare(cfg *config.Config, options common.Map) error {

	rss.table, _ = options.String("table")
	rss.dbKey = joiner.InterfaceKey(options.StringDefault("db_key", string(db_sqlite.InterfaceKey)))

	rss.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(records.InterfaceKey)))
	rss.cleanerKey = joiner.InterfaceKey(options.StringDefault("cleaner_key", string(records.InterfaceCleanerKey)))

	// sqllib.CheckTables

	return nil
}

func (rss *recordsSQLiteStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	db, _ := joinerOp.Interface(rss.dbKey).(*sql.DB)
	if db == nil {
		return fmt.Errorf("no *sql.DB with key %s", rss.dbKey)
	}
	recordsOp, recordsCleanerOp, err := New(db, rss.table)
	if err != nil {
		return errors.CommonError(err, "can't init records.Operator")
	}

	if err = joinerOp.Join(recordsOp, rss.interfaceKey); err != nil {
		return errors.CommonError(err, fmt.Sprintf("can't join *recordsSQLite as records.Operator with key '%s'", rss.interfaceKey))
	}

	if err = joinerOp.Join(recordsCleanerOp, rss.cleanerKey); err != nil {
		return errors.CommonError(err, fmt.Sprintf("can't join *recordsSQLite as db.Cleaner with key '%s'", rss.cleanerKey))
	}

	return nil
}
