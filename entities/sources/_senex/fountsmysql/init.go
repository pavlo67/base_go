package fountsmysql

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/clients"
	"github.com/pavlo67/punctum/basis/config"
	"github.com/pavlo67/punctum/basis/program"
	"github.com/pavlo67/punctum/interfaces/fixturer"
	"github.com/pavlo67/punctum/interfaces/founts"
	"github.com/pavlo67/punctum/interfaces/groups"
	"github.com/pavlo67/punctum/interfaces/starter"
)

const fountTableDefault = "fount"
const fountTagsTableDefault = "fount_tags"
const fountStatTableDefault = "fount_stat"
const scannerStatTableDefault = "scanner_stat"

// Starter ...
func Starter() starter.Operator {
	return &fountsComponent{}
}

type fountsComponent struct {
	mysqlConfig clients.MySQLConfig
	conf        config.PunctumConfig
	index       config.ServerComponentsIndex
}

var Tables map[string]string

func getTables(params map[string]string) []fixturer.MySQLTable {
	fountTable := params["fountTable"]
	if fountTable == "" {
		fountTable = fountTableDefault
	}
	fountTagsTable := params["fountTagsTable"]
	if fountTagsTable == "" {
		fountTagsTable = fountTagsTableDefault
	}
	fountStatTable := params["fountStatTable"]
	if fountStatTable == "" {
		fountStatTable = fountStatTableDefault
	}
	scannerStatTable := params["scannerStatTable"]
	if scannerStatTable == "" {
		scannerStatTable = scannerStatTableDefault
	}

	return []fixturer.MySQLTable{
		{"fountTable", fountTable},
		{"fountTagsTable", fountTagsTable},
		{"fountStatTable", fountStatTable},
		{"scannerStatTable", scannerStatTable},
	}
}

func (fl *fountsComponent) Name() string {
	return founts.InterfaceKey
}

func (sc *fountsComponent) Check(conf config.PunctumConfig, componentKeys config.ComponentKeys, indexPath string) ([]starter.Info, error) {
	var errs []error
	sc.conf = conf
	sc.mysqlConfig, errs = conf.MySQL(componentKeys["mysql"], errs)
	sc.index, errs = config.ComponentIndex(indexPath, basis.CurrentPath(), errs)
	if len(errs) > 0 {
		return nil, basis.JoinErrors(errs)
	}

	return fixturer.CheckMySQLTables(sc.mysqlConfig, basis.CurrentPath(), getTables(sc.index.Params))
}

func (sc *fountsComponent) Setup(conf config.PunctumConfig, componentKeys config.ComponentKeys, indexPath string, data map[string]string) error {
	var errs []error
	sc.conf = conf
	sc.mysqlConfig, errs = conf.MySQL(componentKeys["mysql"], errs)
	sc.index, errs = config.ComponentIndex(indexPath, basis.CurrentPath(), errs)
	if len(errs) > 0 {
		return basis.JoinErrors(errs)
	}

	return fixturer.SetupMySQLTables(sc.mysqlConfig, basis.CurrentPath(), getTables(sc.index.Params), data["reinit"] != "", data["reindex"] != "")
}

func (sc *fountsComponent) Init() error {

	gr, ok := program.Interface(groups.InterfaceKey).(groups.Operator)
	if !ok {
		return errors.New("no group interface found for scannermysql")
	}

	fount, err := NewMySQLFount(
		gr,
		sc.mysqlConfig,
		sc.index.Params["fountTable"],
		sc.index.Params["fountTagsTable"],
		sc.index.Params["fountStatTable"],
		sc.index.Params["scannerStatTable"],
		nil,
		//controller.Managers{rights.Create: IS, rights.Change: IS, rights.View: IS, rights.Delete: IS},
	)
	if err != nil {
		return errors.Wrap(err, "can't init founts")
	}
	err = program.JoinInterface(fount, "founts")
	if err != nil {
		return errors.Wrap(err, "can't join founts")
	}

	return nil
}
