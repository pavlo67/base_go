package auth_server_http

import (
	"fmt"

	"github.com/pavlo67/common/common/auth/auth_jwt"

	"github.com/pavlo67/common/common/server_http"
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"
)

const InterfaceKey joiner.InterfaceKey = "auth_server_http"

func Starter() starter.Operator {
	return &authServerHTTPStarter{}
}

var _ starter.Operator = &authServerHTTPStarter{}

type authServerHTTPStarter struct {
	authKey    joiner.InterfaceKey
	authJWTKey joiner.InterfaceKey

	interfaceKey joiner.InterfaceKey
}

// ------------------------------------------------------------------------------------------------

var l logger.Operator
var authOp, authJWTOp auth.Operator

func (ashs *authServerHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

const onRun = "on authServerHTTPStarter.Run()"

func (ashs *authServerHTTPStarter) Run(_ *config.Envs, options common.Map, joinerOp joiner.Operator, l_ logger.OperatorJ) error {

	l = l_

	ashs.authKey = joiner.InterfaceKey(options.StringDefault("auth_key", string(auth.InterfaceKey)))
	ashs.authJWTKey = joiner.InterfaceKey(options.StringDefault("auth_jwt_key", string(auth_jwt.InterfaceKey)))
	ashs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	// middleware -------------------------------------------------------

	authJWTOp, _ = joinerOp.Interface(ashs.authJWTKey).(auth.Operator)
	if authJWTOp == nil {
		return fmt.Errorf(onRun+": no auth.Operator with key %s", ashs.authJWTKey)
	}

	middleware, err := OnRequestMiddleware(authJWTOp)
	if err != nil || middleware == nil {
		return fmt.Errorf(onRun+": can't create server_http.OnRequestMiddleware(authJWTOp), got %#v, %s", middleware, err)
	}

	srvOp, _ := joinerOp.Interface(server_http.InterfaceKey).(server_http.Operator)
	if srvOp == nil {
		return fmt.Errorf(onRun+": no server_http.Operator with key %s", server_http.InterfaceKey)
	}

	if err = srvOp.HandleMiddleware(middleware); err != nil {
		return errors.Wrap(err, onRun)
	}

	// endpoints --------------------------------------------------------

	if authOp, _ = joinerOp.Interface(ashs.authKey).(auth.Operator); authOp == nil {
		return fmt.Errorf(onRun+": no auth.Operator with key %s", ashs.authKey)
	}

	return Endpoints.Join(joinerOp)
}
