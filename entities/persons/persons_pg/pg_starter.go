package persons_pg

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/db/db_pg"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/rbac"
	"github.com/pavlo67/common/common/starter"
	"github.com/pavlo67/data/entities/persons"

	"github.com/pavlo67/data/components/crud"
)

func Starter() starter.Operator {
	return &personsPgStarter{}
}

var l logger.Operator
var _ starter.Operator = &personsPgStarter{}

type personsPgStarter struct {
	roles rbac.Roles

	dbGetKey joiner.InterfaceKey
	dbSetKey joiner.InterfaceKey

	domain string
	table  string

	interfaceKey joiner.InterfaceKey
	crudKey      joiner.InterfaceKey
	cleanerKey   joiner.InterfaceKey
}

func (pps *personsPgStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (pps *personsPgStarter) Prepare(cfg *config.Config, options common.Map) error {
	pps.roles, _ = options["roles"].(rbac.Roles)

	pps.dbGetKey = joiner.InterfaceKey(options.StringDefault("db_get_key", string(db_pg.InterfaceKey)))
	pps.dbSetKey = joiner.InterfaceKey(options.StringDefault("db_set_key", ""))

	pps.domain = strings.TrimSpace(options.StringDefault("domain", ""))
	if pps.domain == "" {
		return fmt.Errorf("empty domain")
	}

	pps.table = options.StringDefault("table", persons.CollectionDefault)

	interfaceKey, interfaceCRUDKey, cleanerKey := persons.InterfaceKey, persons.InterfaceCRUDKey, joiner.InterfaceKey("")
	if pps.roles.Has(crud.RoleTester) {
		interfaceKey, interfaceCRUDKey, cleanerKey = persons.InterfaceTestKey, persons.InterfaceCRUDTestKey, persons.InterfaceCleanerKey
	}

	pps.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(interfaceKey)))
	pps.crudKey = joiner.InterfaceKey(options.StringDefault("crud_key", string(interfaceCRUDKey)))
	pps.cleanerKey = joiner.InterfaceKey(options.StringDefault("cleaner_key", string(cleanerKey)))

	return nil
}

func (pps *personsPgStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	dbGet, _ := joinerOp.Interface(pps.dbGetKey).(*sql.DB)
	if dbGet == nil {
		return fmt.Errorf("no *sql.DB with key %s", pps.dbGetKey)
	}
	var dbSet *sql.DB
	if pps.dbSetKey != "" {
		dbSet, _ = joinerOp.Interface(pps.dbSetKey).(*sql.DB)
		if dbSet == nil {
			return fmt.Errorf("no *sql.DB with key %s", pps.dbSetKey)
		}
	}

	personsOp, personsCleanerOp, err := New(dbGet, dbSet, pps.domain, pps.table)
	if err != nil {
		return errors.Wrap(err, "can't init *personsStub{} as persons.Operator")
	}

	if err = joinerOp.Join(personsOp, pps.interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join *personsPg{} as persons.Operator with key '%s'", pps.interfaceKey)
	}

	if pps.crudKey != "" {
		if crudOp, err := persons.OperatorCRUD(personsOp, pps.roles); err != nil {
			return err
		} else if err = joinerOp.Join(crudOp, pps.crudKey); err != nil {
			return errors.Wrapf(err, "can't join *persons.OperatorCRUD as crud.Operator with key '%s'", pps.crudKey)
		}
	}

	if pps.cleanerKey != "" {
		if err = joinerOp.Join(personsCleanerOp, pps.cleanerKey); err != nil {
			return errors.Wrapf(err, "can't join *personsPg{} as db.Cleaner with key '%s'", pps.cleanerKey)
		}
	}

	return nil
}
