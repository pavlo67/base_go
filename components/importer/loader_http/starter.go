package loader_http

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/workshop/constructions/filesloader"
)

func Starter() starter.Operator {
	return &loaderHTTPStarter{}
}

var l logger.Operator
var _ starter.Operator = &loaderHTTPStarter{}

type loaderHTTPStarter struct {
	interfaceKey joiner.InterfaceKey
	pathToStore  string

	// TODO: use proxies
}

// ------------------------------------------------

func (fl *loaderHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (fl *loaderHTTPStarter) Init(cfgCommon, cfg *config.Config, lCommon logger.Operator, options common.Map) ([]common.Map, error) {
	l = lCommon

	fl.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(filesloader.InterfaceKey)))
	fl.pathToStore = options.StringDefault("path_to_store", "./")

	return nil, nil
}

func (fl *loaderHTTPStarter) Setup() error {
	return nil
}

func (fl *loaderHTTPStarter) Run(joinerOp joiner.Operator) error {
	flOp, _, err := New(fl.pathToStore)
	if err != nil {
		return errors.Wrap(err, "can't init filesloader.ActorKey")
	}

	err = joinerOp.Join(flOp, fl.interfaceKey)
	if err != nil {
		return errors.Wrapf(err, "can't join *filesloaderHTTP as filesloader.ActorKey with key '%s'", fl.interfaceKey)
	}

	return nil
}
