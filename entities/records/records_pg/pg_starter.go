package records_pg

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pavlo67/data/entities"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/db/db_pg"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/rbac"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/data/entities/records"
)

func Starter() starter.Operator {
	return &recordsPgStarter{}
}

var l logger.Operator
var _ starter.Operator = &recordsPgStarter{}

type recordsPgStarter struct {
	roles rbac.Roles

	dbGetKey joiner.InterfaceKey
	dbSetKey joiner.InterfaceKey

	domain string
	table  string

	interfaceKey joiner.InterfaceKey
	crudKey      joiner.InterfaceKey
	cleanerKey   joiner.InterfaceKey
}

func (rps *recordsPgStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (rps *recordsPgStarter) Prepare(cfg *config.Config, options common.Map) error {
	rps.roles, _ = options["roles"].(rbac.Roles)

	rps.dbGetKey = joiner.InterfaceKey(options.StringDefault("db_get_key", string(db_pg.InterfaceKey)))
	rps.dbSetKey = joiner.InterfaceKey(options.StringDefault("db_set_key", ""))

	rps.table = options.StringDefault("table", records.CollectionDefault)
	rps.domain = strings.TrimSpace(options.StringDefault("domain", ""))
	if rps.domain == "" {
		return fmt.Errorf("empty domain")
	}

	interfaceKey, interfaceCRUDKey, cleanerKey := records.InterfaceKey, records.InterfaceCRUDKey, joiner.InterfaceKey("")
	if rps.roles.Has(entities.RoleTester) {
		interfaceKey, interfaceCRUDKey, cleanerKey = records.InterfaceTestKey, records.InterfaceCRUDTestKey, records.InterfaceCleanerKey
	}

	rps.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(interfaceKey)))
	rps.crudKey = joiner.InterfaceKey(options.StringDefault("crud_key", string(interfaceCRUDKey)))
	rps.cleanerKey = joiner.InterfaceKey(options.StringDefault("cleaner_key", string(cleanerKey)))

	return nil
}

func (rps *recordsPgStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.OperatorCRUD with key %s", logger.InterfaceKey)
	}

	dbGet, _ := joinerOp.Interface(rps.dbGetKey).(*sql.DB)
	if dbGet == nil {
		return fmt.Errorf("no *sql.DB with key %s", rps.dbGetKey)
	}
	var dbSet *sql.DB
	if rps.dbSetKey != "" {
		dbSet, _ = joinerOp.Interface(rps.dbSetKey).(*sql.DB)
		if dbSet == nil {
			return fmt.Errorf("no *sql.DB with key %s", rps.dbSetKey)
		}
	}

	recordsOp, recordsCleanerOp, err := New(dbGet, dbSet, rps.domain, rps.table)
	if err != nil {
		return errors.Wrap(err, "can't init *recordsStub{} as records.OperatorCRUD")
	}

	if err = joinerOp.Join(recordsOp, rps.interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join *recordsPg{} as records.OperatorCRUD with key '%s'", rps.interfaceKey)
	}

	if rps.crudKey != "" {
		if crudOp, err := records.OperatorCRUD(recordsOp, rps.roles); err != nil {
			return err
		} else if err = joinerOp.Join(crudOp, rps.crudKey); err != nil {
			return errors.Wrapf(err, "can't join *records.OperatorCRUD as crud.OperatorCRUD with key '%s'", rps.crudKey)
		}
	}

	if rps.cleanerKey != "" {
		if err = joinerOp.Join(recordsCleanerOp, rps.cleanerKey); err != nil {
			return errors.Wrapf(err, "can't join *recordsPg{} as db.Cleaner with key '%s'", rps.cleanerKey)
		}
	}

	return nil
}
