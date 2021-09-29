package records01_pg

import (
	"database/sql"
	"fmt"

	"github.com/pavlo67/common/common/db/db_pg"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/data/entities/records01"
)

func Starter() starter.Operator {
	return &records01PgStarter{}
}

var l logger.Operator
var _ starter.Operator = &records01PgStarter{}

type records01PgStarter struct {
	dbGetKey joiner.InterfaceKey
	dbSetKey joiner.InterfaceKey

	table string

	interfaceKey joiner.InterfaceKey
	cleanerKey   joiner.InterfaceKey
}

func (p01ps *records01PgStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (p01ps *records01PgStarter) Prepare(cfg *config.Config, options common.Map) error {
	p01ps.dbGetKey = joiner.InterfaceKey(options.StringDefault("db_get_key", string(db_pg.InterfaceKey)))
	p01ps.dbSetKey = joiner.InterfaceKey(options.StringDefault("db_set_key", ""))

	p01ps.table = options.StringDefault("table", records01.CollectionDefault)

	p01ps.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(records01.InterfaceKey)))
	p01ps.cleanerKey = joiner.InterfaceKey(options.StringDefault("cleaner_key", string(records01.InterfaceCleanerKey)))

	return nil
}

func (p01ps *records01PgStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	dbGet, _ := joinerOp.Interface(p01ps.dbGetKey).(*sql.DB)
	if dbGet == nil {
		return fmt.Errorf("no *sql.DB with key %s", p01ps.dbGetKey)
	}
	var dbSet *sql.DB
	if p01ps.dbSetKey != "" {
		dbSet, _ = joinerOp.Interface(p01ps.dbSetKey).(*sql.DB)
		if dbSet == nil {
			return fmt.Errorf("no *sql.DB with key %s", p01ps.dbSetKey)
		}
	}

	recordsOp, recordsCleanerOp, err := New(dbGet, dbSet, p01ps.table)
	if err != nil {
		return errors.Wrap(err, "can't init *recordsStub{} as records.Operator")
	}

	if err = joinerOp.Join(recordsOp, p01ps.interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join *recordsStub{} as records.Operator with key '%s'", p01ps.interfaceKey)
	}

	if err = joinerOp.Join(recordsCleanerOp, p01ps.cleanerKey); err != nil {
		return errors.Wrapf(err, "can't join *recordsStub{} as db.Cleaner with key '%s'", p01ps.cleanerKey)
	}

	return nil
}
