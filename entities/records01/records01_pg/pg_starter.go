package records01_pg

import (
	"database/sql"
	"fmt"

	"github.com/pavlo67/common/common/rbac"

	"github.com/pavlo67/data/entities/records01"

	"github.com/pavlo67/common/common/db/db_pg"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"
)

func Starter() starter.Operator {
	return &records01PgStarter{}
}

var l logger.Operator
var _ starter.Operator = &records01PgStarter{}

type records01PgStarter struct {
	roles rbac.Roles

	dbGetKey joiner.InterfaceKey
	dbSetKey joiner.InterfaceKey

	table string

	interfaceKey joiner.InterfaceKey
	crudKey      joiner.InterfaceKey
	cleanerKey   joiner.InterfaceKey
}

func (r01ps *records01PgStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (r01ps *records01PgStarter) Prepare(cfg *config.Config, options common.Map) error {
	r01ps.roles, _ = options["roles"].(rbac.Roles)

	r01ps.dbGetKey = joiner.InterfaceKey(options.StringDefault("db_get_key", string(db_pg.InterfaceKey)))
	r01ps.dbSetKey = joiner.InterfaceKey(options.StringDefault("db_set_key", ""))

	r01ps.table = options.StringDefault("table", records01.CollectionDefault)

	r01ps.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(records01.InterfaceKey)))
	r01ps.crudKey = joiner.InterfaceKey(options.StringDefault("crud_key", ""))
	r01ps.cleanerKey = joiner.InterfaceKey(options.StringDefault("cleaner_key", string(records01.InterfaceCleanerKey)))

	return nil
}

func (r01ps *records01PgStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	dbGet, _ := joinerOp.Interface(r01ps.dbGetKey).(*sql.DB)
	if dbGet == nil {
		return fmt.Errorf("no *sql.DB with key %s", r01ps.dbGetKey)
	}
	var dbSet *sql.DB
	if r01ps.dbSetKey != "" {
		dbSet, _ = joinerOp.Interface(r01ps.dbSetKey).(*sql.DB)
		if dbSet == nil {
			return fmt.Errorf("no *sql.DB with key %s", r01ps.dbSetKey)
		}
	}

	recordsOp, recordsCleanerOp, err := New(dbGet, dbSet, r01ps.table)
	if err != nil {
		return errors.Wrap(err, "can't init *recordsStub{} as records.Operator")
	}

	if err = joinerOp.Join(recordsOp, r01ps.interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join *recordsStub{} as records.Operator with key '%s'", r01ps.interfaceKey)
	}

	if r01ps.crudKey != "" {
		if crudOp, err := records01.OperatorCRUD(recordsOp, r01ps.roles); err != nil {
			return err
		} else if err = joinerOp.Join(crudOp, r01ps.crudKey); err != nil {
			return errors.Wrapf(err, "can't join *records.OperatorCRUD as crud.Operator with key '%s'", r01ps.crudKey)
		}
	}

	if err = joinerOp.Join(recordsCleanerOp, r01ps.cleanerKey); err != nil {
		return errors.Wrapf(err, "can't join *recordsStub{} as db.Cleaner with key '%s'", r01ps.cleanerKey)
	}

	return nil
}
