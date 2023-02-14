package data_pg

import (
	"strings"

	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/identity"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"
	"github.com/pavlo67/workshop/components/data"
	"github.com/pavlo67/workshop/components/tagger"
)

func Starter() starter.Operator {
	return &dataPgStarter{}
}

var l logger.Operator
var _ starter.Operator = &dataPgStarter{}

type dataPgStarter struct {
	config config.Access
	domain identity.Domain
	table  string

	interfaceKey joiner.InterfaceKey
	cleanerKey   joiner.InterfaceKey

	noTagger bool
}

func (dp *dataPgStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (dp *dataPgStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon

	var cfgPG config.Access
	err := cfg.Value("pg", &cfgPG)
	if err != nil {
		return nil, err
	}

	dp.domain = identity.Domain(strings.TrimSpace(cfg.ServiceName()))
	if dp.domain == "" {
		return nil, errors.New("no service name")
	}
	dp.config = cfgPG
	dp.table, _ = options.String("table")
	dp.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(data.InterfaceKey)))
	dp.cleanerKey = joiner.InterfaceKey(options.StringDefault("cleaner_key", string(data.CleanerInterfaceKey)))

	dp.noTagger = options.IsTrue("no_tagger")

	// sqllib.CheckTables

	return nil, nil
}

func (dp *dataPgStarter) Setup() error {
	return nil

	//return sqllib.SetupTables(
	//	sm.mysqlConfig,
	//	sm.index.MySQL,
	//	[]config.Table{{Key: "table", Title: sm.table}},
	//)
}

func (dp *dataPgStarter) Run(joinerOp joiner.Operator) error {
	var ok bool
	var taggerOp tagger.Operator
	var taggercleanerOp db.Cleaner

	if !dp.noTagger {
		taggerOp, ok = joinerOp.Interface(tagger.InterfaceKey).(tagger.Operator)
		if !ok {
			return errors.Errorf("no tagger.Operator with key %s", tagger.InterfaceKey)
		}

		taggercleanerOp, ok = joinerOp.Interface(tagger.CleanerInterfaceKey).(db.Cleaner)
		if !ok {
			return errors.Errorf("no tagger.Cleaner with key %s", tagger.InterfaceKey)
		}
	}

	dataOp, datacleanerOp, err := New(dp.config, dp.domain, dp.table, dp.interfaceKey, taggerOp, taggercleanerOp)
	if err != nil {
		return errors.Wrap(err, "can't init *dataPG as data.Operator")
	}

	err = joinerOp.Join(dataOp, dp.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join *dataPG as data.Operator with key '%s'", dp.interfaceKey)
	}

	err = joinerOp.Join(datacleanerOp, dp.cleanerKey)
	if err != nil {
		return errors.Wrapf(err, "can't join *dataPG as db.Cleaner with key '%s'", dp.cleanerKey)
	}

	return nil
}
