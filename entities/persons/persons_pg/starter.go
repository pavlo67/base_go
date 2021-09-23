package persons_pg

//import (
//	"github.com/pkg/errors"
//
//	"github.com/pavlo67/common/common"
//	"github.com/pavlo67/common/common/config"
//	"github.com/pavlo67/common/common/crud"
//	"github.com/pavlo67/common/common/joiner"
//	"github.com/pavlo67/common/common/logger"
//	"github.com/pavlo67/common/common/starter"
//	"github.com/pavlo67/workshop/components/users"
//	"github.com/pavlo67/workshop/components/tagger"
//)
//
//func Starter() starter.Operator {
//	return &usersPgStarter{}
//}
//
//var l logger.Operator
//var _ starter.Operator = &usersPgStarter{}
//
//type usersPgStarter struct {
//	config config.Access
//	table  string
//
//	interfaceKey joiner.InterfaceKey
//	cleanerKey   joiner.InterfaceKey
//
//	noTagger bool
//}
//
//func (dp *usersPgStarter) Name() string {
//	return logger.GetCallInfo().PackageName
//}
//
//func (dp *usersPgStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
//	l = lCommon
//
//	var cfgPG config.Access
//	err := cfg.Value("pg", &cfgPG)
//	if err != nil {
//		return nil, err
//	}
//
//	dp.config = cfgPG
//	dp.table, _ = options.String("table")
//	dp.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(users.InterfaceKey)))
//	dp.cleanerKey = joiner.InterfaceKey(options.StringDefault("cleaner_key", string(users.CleanerInterfaceKey)))
//
//	dp.noTagger = options.IsTrue("no_tagger")
//
//	// sqllib.CheckTables
//
//	return nil, nil
//}
//
//func (dp *usersPgStarter) Setup() error {
//	return nil
//
//	//return sqllib.SetupTables(
//	//	sm.mysqlConfig,
//	//	sm.index.MySQL,
//	//	[]config.Table{{Key: "table", Title: sm.table}},
//	//)
//}
//
//func (dp *usersPgStarter) Run(joinerOp joiner.Operator) error {
//	var ok bool
//	var taggerOp tagger.Operator
//	var taggercleanerOp crud.Cleaner
//
//	if !dp.noTagger {
//		taggerOp, ok = joinerOp.Interface(tagger.InterfaceKey).(tagger.Operator)
//		if !ok {
//			return errors.Errorf("no tagger.Operator with key %s", tagger.InterfaceKey)
//		}
//
//		taggercleanerOp, ok = joinerOp.Interface(tagger.CleanerInterfaceKey).(crud.Cleaner)
//		if !ok {
//			return errors.Errorf("no tagger.Cleaner with key %s", tagger.InterfaceKey)
//		}
//	}
//
//	usersOp, userscleanerOp, err := New(dp.config, dp.table, dp.interfaceKey, taggerOp, taggercleanerOp)
//	if err != nil {
//		return errors.Wrap(err, "can't init *usersPG as users.ActorKey")
//	}
//
//	err = joinerOp.Join(usersOp, dp.interfaceKey)
//	if err != nil {
//		return errors.Wrapf(err, "can't join *usersPG as users.ActorKey with key '%s'", dp.interfaceKey)
//	}
//
//	err = joinerOp.Join(userscleanerOp, dp.cleanerKey)
//	if err != nil {
//		return errors.Wrapf(err, "can't join *usersPG as crud.Cleaner with key '%s'", dp.cleanerKey)
//	}
//
//	return nil
//}
