package flowmysql

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/clients"
	"github.com/pavlo67/punctum/basis/config"
	"github.com/pavlo67/punctum/basis/program"
	"github.com/pavlo67/punctum/interfaces/fixturer"
	"github.com/pavlo67/punctum/interfaces/flow"
	"github.com/pavlo67/punctum/interfaces/groups"
	"github.com/pavlo67/punctum/interfaces/starter"
)

const flowTableDefault = "flow"

// Starter ...
func Starter() starter.Operator {
	return &flowComponent{}
}

type flowComponent struct {
	mysqlConfig clients.MySQLConfig
	conf        config.PunctumConfig
	index       config.ServerComponentsIndex
}

func getTables(params map[string]string) []fixturer.MySQLTable {

	flowTable := params["flowTable"]
	if flowTable == "" {
		flowTable = flowTableDefault
	}
	return []fixturer.MySQLTable{
		{"flowTable", flowTable},
	}
}

func (fl *flowComponent) Name() string {
	return flow.InterfaceKey
}

func (fl *flowComponent) Check(conf config.PunctumConfig, componentKeys config.ComponentKeys, indexPath string) ([]starter.Info, error) {

	var errs []error
	fl.conf = conf
	fl.mysqlConfig, errs = conf.MySQL(componentKeys["mysql"], errs)
	fl.index, errs = config.ComponentIndex(indexPath, basis.CurrentPath(), errs)
	if len(errs) > 0 {
		return nil, basis.JoinErrors(errs)
	}
	return fixturer.CheckMySQLTables(fl.mysqlConfig, basis.CurrentPath(), getTables(fl.index.Params))
}

func (fl *flowComponent) Setup(conf config.PunctumConfig, componentKeys config.ComponentKeys, indexPath string, data map[string]string) error {

	var errs []error
	fl.conf = conf
	fl.mysqlConfig, errs = conf.MySQL(componentKeys["mysql"], errs)
	fl.index, errs = config.ComponentIndex(indexPath, basis.CurrentPath(), errs)
	if len(errs) > 0 {
		return basis.JoinErrors(errs)
	}
	return fixturer.SetupMySQLTables(fl.mysqlConfig, basis.CurrentPath(), getTables(fl.index.Params), data["reinit"] != "", data["reindex"] != "")
}

func (fl *flowComponent) Init() error {

	gr, ok := program.Interface(groups.InterfaceKey).(groups.Operator)
	if !ok {
		return errors.New("no group interface found for scannermysql")
	}
	flowOp, err := NewFlowMySQL(gr, fl.mysqlConfig, fl.index.Params["flowTable"], nil)
	if err != nil {
		return errors.Wrap(err, "can't init flow ")
	}
	err = program.JoinInterface(flowOp, flow.InterfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join flow ")
	}
	return nil
}
