package persons01_pg

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
	"github.com/pavlo67/data/entities/persons01"
)

func Starter() starter.Operator {
	return &persons01PgStarter{}
}

var l logger.Operator
var _ starter.Operator = &persons01PgStarter{}

type persons01PgStarter struct {
	dbGetKey joiner.InterfaceKey
	dbSetKey joiner.InterfaceKey

	table string

	interfaceKey joiner.InterfaceKey
	crudKey      joiner.InterfaceKey
	cleanerKey   joiner.InterfaceKey
}

func (p01ps *persons01PgStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (p01ps *persons01PgStarter) Prepare(cfg *config.Config, options common.Map) error {
	p01ps.dbGetKey = joiner.InterfaceKey(options.StringDefault("db_get_key", string(db_pg.InterfaceKey)))
	p01ps.dbSetKey = joiner.InterfaceKey(options.StringDefault("db_set_key", ""))

	p01ps.table = options.StringDefault("table", persons01.CollectionDefault)

	p01ps.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(persons01.InterfaceKey)))
	p01ps.crudKey = joiner.InterfaceKey(options.StringDefault("crud_key", ""))
	p01ps.cleanerKey = joiner.InterfaceKey(options.StringDefault("cleaner_key", string(persons01.InterfaceCleanerKey)))

	return nil
}

func (p01ps *persons01PgStarter) Run(joinerOp joiner.Operator) error {
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

	personsOp, personsCleanerOp, err := New(dbGet, dbSet, p01ps.table)
	if err != nil {
		return errors.Wrap(err, "can't init *personsStub{} as persons.Operator")
	}

	if err = joinerOp.Join(personsOp, p01ps.interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join *personsStub{} as persons.Operator with key '%s'", p01ps.interfaceKey)
	}

	if p01ps.crudKey != "" {
		if crudOp, err := persons01.OperatorCRUD(personsOp); err != nil {
			return err
		} else if err = joinerOp.Join(crudOp, p01ps.crudKey); err != nil {
			return errors.Wrapf(err, "can't join *persons.OperatorCRUD as crud.Operator with key '%s'", p01ps.crudKey)
		}
	}

	if err = joinerOp.Join(personsCleanerOp, p01ps.cleanerKey); err != nil {
		return errors.Wrapf(err, "can't join *personsStub{} as db.Cleaner with key '%s'", p01ps.cleanerKey)
	}

	return nil
}
