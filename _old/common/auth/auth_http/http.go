package auth_http

import (
	"encoding/json"

	"github.com/pavlo67/common/common/httplib"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/server_http"
)

var _ auth.Operator = &authHTTP{}

type authHTTP struct {
	serverConfig server_http.Config
	logFilePath  string
}

const onNew = "on authHTTP.New()"

func New(serverConfig server_http.Config, logFilePath string) (auth.Operator, error) {
	authOp := authHTTP{
		serverConfig: serverConfig,
		logFilePath:  logFilePath,
	}

	return &authOp, nil
}

func (authOp *authHTTP) SetCreds(actor auth.Actor, toSet auth.Creds) (*auth.Creds, error) {
	ep := authOp.serverConfig.EndpointsSettled[auth.IntefaceKeySetCreds]
	serverURL := authOp.serverConfig.Host + authOp.serverConfig.Port + authOp.serverConfig.Prefix + ep.Path

	requestBody, err := json.Marshal(toSet)
	if err != nil {
		return nil, errors.Wrapf(err, onAuthenticate+": can't marshal toSet(%#v)", toSet)
	}

	var creds *auth.Creds
	if err := httplib.Request(nil, serverURL, ep.Method, server_http.SetCreds(actor.Creds), requestBody, &creds, l, authOp.logFilePath); err != nil {
		return nil, err
	}

	return creds, nil
}

// Authenticate can require to call .SetCredsByKey() first and to use some session-generated creds
const onAuthenticate = "on authHTTP.Authenticate()"

func (authOp *authHTTP) Authenticate(toAuth auth.Creds) (*auth.Actor, error) {
	ep := authOp.serverConfig.EndpointsSettled[auth.IntefaceKeyAuthenticate]
	serverURL := authOp.serverConfig.Host + authOp.serverConfig.Port + authOp.serverConfig.Prefix + ep.Path

	requestBody, err := json.Marshal(toAuth)
	if err != nil {
		return nil, errors.Wrapf(err, onAuthenticate+": can't marshal toAuth(%#v)", toAuth)
	}

	var actor *auth.Actor
	if err = httplib.Request(nil, serverURL, ep.Method, nil, requestBody, &actor, l, authOp.logFilePath); err != nil {
		return nil, err
	}

	return actor, nil
}
