package db_sqlite

//import (
//	"fmt"
//	"os"
//
//	"github.com/pavlo67/common/common"
//	"github.com/pavlo67/common/common/config"
//	"github.com/pavlo67/common/common/errors"
//	"github.com/pavlo67/common/common/joiner"
//	"github.com/pavlo67/common/common/logger"
//	"github.com/pavlo67/common/common/sqllib/sqllib_sqlite"
//	"github.com/pavlo67/common/common/starter"
//)
//
//const InterfaceKey joiner.InterfaceKey = "db_sqlite"
//
//func Starter() starter.Operator {
//	return &connectSQLiteStarter{}
//}
//
//var l logger.Operator
//var _ starter.Operator = &connectSQLiteStarter{}
//
//type connectSQLiteStarter struct {
//	cfgSQLite config.Access
//
//	interfaceKey joiner.InterfaceKey
//}
//
//func (css *connectSQLiteStarter) Name() string {
//	return logger.GetCallInfo().PackageName
//}
//
//const onRun = "on connectSQLiteStarter.Run()"
//
//func (css *connectSQLiteStarter) Run(env *config.Envs, options common.Map, joinerOp joiner.Operator, l_ logger.OperatorJ) error {
//	l = l_
//
//	if err := env.Value(options.StringDefault("db_key", "db_sqlite"), &css.cfgSQLite); err != nil {
//		return err
//	}
//
//	css.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))
//
//	if os.Getenv("SHOW_CONNECTS") != "" {
//		l.Infof("CONNECTING TO SQLITE: %#v", css.cfgSQLite)
//	}
//
//	db, err := sqllib_sqlite.Connect(css.cfgSQLite)
//	if err != nil || db == nil {
//		return errors.CommonError(err, fmt.Sprintf(onRun+": got %#v", db))
//	}
//
//	if err = joinerOp.Join(db, css.interfaceKey); err != nil {
//		return errors.CommonError(err, fmt.Sprintf("can't join *sql.DB with key '%s'", css.interfaceKey))
//	}
//
//	return nil
//}
