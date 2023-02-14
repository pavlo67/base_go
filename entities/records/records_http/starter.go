package records_http

import (
	"fmt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/server/server_http"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/data/entities/records"
)

func Starter() starter.Operator {
	return &recordsHTTPStarter{}
}

var l logger.Operator
var _ starter.Operator = &recordsHTTPStarter{}

type recordsHTTPStarter struct {
	pagesConfig server_http.Config
	restConfig  server_http.Config

	interfaceKey joiner.InterfaceKey
}

func (ahs *recordsHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ahs *recordsHTTPStarter) Prepare(cfg *config.Config, options common.Map) error {
	var ok bool

	if ahs.pagesConfig, ok = options["pages_config"].(server_http.Config); !ok {
		return fmt.Errorf(`no server_http.Config in options["pages_config"]`)
	}
	if ahs.restConfig, ok = options["rest_config"].(server_http.Config); !ok {
		return fmt.Errorf(`no server_http.Config in options["rest_config"]`)
	}

	ahs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(records.InterfaceKey)))

	return nil
}

func (ahs *recordsHTTPStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.OperatorV2 with key %s", logger.InterfaceKey)
	}

	recordsOp, err := New(ahs.pagesConfig, ahs.restConfig)
	if err != nil {
		return errors.Wrap(err, "can't init *recordsHTTP{} as records.OperatorV2")
	}

	if err = joinerOp.Join(recordsOp, ahs.interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join *recordsHTTP{} as records.OperatorV2 with key '%s'", ahs.interfaceKey)
	}

	return nil
}
