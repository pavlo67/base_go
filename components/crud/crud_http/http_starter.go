package crud_http

import (
	"fmt"

	"github.com/pavlo67/common/common/crud/crud_server_http"
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/server/server_http"
	"github.com/pavlo67/common/common/starter"
)

func Starter() starter.Operator {
	return &crudHTTPStarter{}
}

var l logger.Operator
var _ starter.Operator = &crudHTTPStarter{}

type crudHTTPStarter struct {
	serverConfig server_http.Config
	interfaceKey joiner.InterfaceKey
}

func (ahs *crudHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ahs *crudHTTPStarter) Prepare(cfg *config.Config, options common.Map) error {

	var access config.Access
	if err := cfg.Value("crud_http", &access); err != nil {
		return err
	}

	// TODO!!! pass for each server separately
	prefix := options.StringDefault("prefix", "")

	var ok bool
	if ahs.serverConfig, ok = options["server_config"].(server_http.Config); !ok {
		return errors.New("no server config for crudHTTPStarter")
	}

	ahs.serverConfig.CompleteDirectly(crud_server_http.Endpoints, access.Host, access.Port, prefix)

	ahs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(crud.InterfaceKey)))

	return nil
}

//} else if endpointsPtr, ok := options["endpoints"].(*server_http.Endpoints); ok {
//	ihs.endpoints = *endpointsPtr

func (ahs *crudHTTPStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	crudOp, err := New(ahs.serverConfig)
	if err != nil {
		return errors.Wrap(err, "can't init *crudHTTP{} as crud.Operator")
	}

	if err = joinerOp.Join(crudOp, ahs.interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join *crudHTTP{} as crud.Operator with key '%s'", ahs.interfaceKey)
	}

	return nil
}
