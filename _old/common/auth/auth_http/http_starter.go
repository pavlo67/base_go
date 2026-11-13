package auth_http

import (
	"github.com/pavlo67/common/common/auth/auth_server_http"
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/server_http"
	"github.com/pavlo67/common/common/starter"
)

func Starter() starter.Operator {
	return &authHTTPStarter{}
}

const InterfaceKey joiner.InterfaceKey = "auth_http"

var l logger.Operator
var _ starter.Operator = &authHTTPStarter{}

type authHTTPStarter struct {
	serverConfig server_http.Config
	interfaceKey joiner.InterfaceKey
}

func (ahs *authHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

//} else if endpointsPtr, ok := options["endpoints"].(*server_http.Endpoints); ok {
//	ihs.endpoints = *endpointsPtr

func (ahs *authHTTPStarter) Run(cfg *config.Envs, options common.Map, joinerOp joiner.Operator, l_ logger.OperatorJ) error {
	l = l_

	var access config.Access
	if err := cfg.Value("auth_http", &access); err != nil {
		return err
	}

	var ok bool
	if ahs.serverConfig, ok = options["server_config"].(server_http.Config); !ok {
		return errors.New("no server config for authHTTPStarter")
	}
	ignoreAbsent, _ := options["ignore_absent"].(bool)

	logFilePath, _ := options["log_file"].(string)

	if err := ahs.serverConfig.CompleteDirectly(auth_server_http.Endpoints, access.Host, access.Port, ignoreAbsent); err != nil {
		return err
	}

	ahs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	authOp, err := New(ahs.serverConfig, logFilePath)
	if err != nil {
		return errors.Wrap(err, "can't init *authHTTP{} as auth.Operator")
	}

	if err = joinerOp.Join(authOp, ahs.interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join *authHTTP{} as auth.Operator with key '%s'", ahs.interfaceKey)
	}

	return nil
}
