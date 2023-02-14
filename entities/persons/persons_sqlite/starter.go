package persons_sqlite

import (
	"database/sql"
	"fmt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/db/db_sqlite"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/data/entities/persons"
)

func Starter() starter.Operator {
	return &personsSQLiteStarter{}
}

var l logger.Operator
var _ starter.Operator = &personsSQLiteStarter{}

type personsSQLiteStarter struct {
	connectKey joiner.InterfaceKey
	table      string

	interfaceKey joiner.InterfaceKey
	cleanerKey   joiner.InterfaceKey
}

func (rss *personsSQLiteStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (rss *personsSQLiteStarter) Prepare(cfg *config.Config, options common.Map) error {

	rss.table, _ = options.String("table")
	rss.connectKey = joiner.InterfaceKey(options.StringDefault("connect_key", string(db_sqlite.InterfaceKey)))
	rss.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(persons.InterfaceKey)))
	rss.cleanerKey = joiner.InterfaceKey(options.StringDefault("cleaner_key", string(persons.InterfaceCleanerKey)))

	// sqllib.CheckTables

	return nil
}

func (rss *personsSQLiteStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	db, _ := joinerOp.Interface(rss.connectKey).(*sql.DB)
	if db == nil {
		return fmt.Errorf("no *sql.DB with key %s", rss.connectKey)
	}

	personsOp, personsCleanerOp, err := New(db, rss.table)
	if err != nil {
		return errors.CommonError(err, "can't init personsSQLite as persons.Operator")
	}

	if err = joinerOp.Join(personsOp, rss.interfaceKey); err != nil {
		return errors.CommonError(err, fmt.Sprintf("can't join *personsSQLite as persons.Operator with key '%s'", rss.interfaceKey))
	}

	if err = joinerOp.Join(personsCleanerOp, rss.cleanerKey); err != nil {
		return errors.CommonError(err, fmt.Sprintf("can't join *personsSQLite as db.Cleaner with key '%s'", rss.cleanerKey))
	}

	return nil
}
