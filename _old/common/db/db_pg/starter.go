package db_pg

import (
	"fmt"
	"os"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/sqllib/sqllib_pg"
	"github.com/pavlo67/common/common/starter"
)

const InterfaceKey joiner.InterfaceKey = "db_pg"

func Starter() starter.Operator {
	return &connectPgStarter{}
}

var l logger.Operator
var _ starter.Operator = &connectPgStarter{}

type connectPgStarter struct {
	cfgPg config.Access

	interfaceKey joiner.InterfaceKey
}

func (cps *connectPgStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (cps *connectPgStarter) Run(cfg *config.Envs, options common.Map, joinerOp joiner.Operator, l_ logger.OperatorJ) error {
	l = l_

	if err := cfg.Value(options.StringDefault("db_key", "db_pg"), &cps.cfgPg); err != nil {
		return err
	}

	cps.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	if os.Getenv("SHOW_CONNECTS") != "" {
		l.Infof("CONNECTING TO PG: %#v", cps.cfgPg)
	}

	db, err := sqllib_pg.Connect(cps.cfgPg)
	if err != nil || db == nil {
		return errors.CommonError(err, fmt.Sprintf("got %#v", db))
	}

	if err = joinerOp.Join(db, cps.interfaceKey); err != nil {
		return errors.CommonError(err, fmt.Sprintf("can't join *sql.DB with key '%s'", cps.interfaceKey))
	}

	return nil
}
