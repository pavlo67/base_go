package sources_pg

// import (
// 	"github.com/pkg/errors"

// 	"github.com/pavlo67/common/common"
// 	"github.com/pavlo67/common/common/config"
// 	"github.com/pavlo67/common/common/joiner"
// 	"github.com/pavlo67/common/common/logger"
// 	"github.com/pavlo67/common/common/starter"

// 	"github.com/pavlo67/workshop/components/packs"
// )

// func Starter() starter.ActorKey {
// 	return &packsPgStarter{}
// }

// var l logger.ActorKey
// var _ starter.ActorKey = &packsPgStarter{}

// type packsPgStarter struct {
// 	config       config.Access
// 	table        string
// 	interfaceKey joiner.InterfaceKey
// }

// func (ps *packsPgStarter) Name() string {
// 	return logger.GetCallInfo().PackageName
// }

// func (ps *packsPgStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.ActorKey, options common.Map) ([]common.Map, error) {
// 	l = lCommon

// 	var cfgPg config.Access
// 	err := cfg.Value("pg", &cfgPg)
// 	if err != nil {
// 		return nil, err
// 	}

// 	ps.config = cfgPg
// 	ps.table, _ = options.String("table")
// 	ps.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(packs.InterfaceKey)))

// 	return nil, nil
// }

// func (ps *packsPgStarter) Setup() error {
// 	return nil
// }

// func (ps *packsPgStarter) Run(joinerOp joiner.ActorKey) error {
// 	packsOp, _, err := New(ps.config, ps.table, ps.interfaceKey)
// 	if err != nil {
// 		return errors.Wrap(err, "can't init packs.ActorKey")
// 	}

// 	err = joinerOp.Join(packsOp, ps.interfaceKey)
// 	if err != nil {
// 		return errors.Wrapf(err, "can't join &packsPg as packs.ActorKey with key '%s'", ps.interfaceKey)
// 	}

// 	return nil
// }
