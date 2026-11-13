package auth_server_http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/server_http"
)

var Endpoints = server_http.Endpoints{
	authenticateEndpoint,
	setCredsEndpoint,
}

//var bodyParams = json.RawMessage(`{
//   "in": "body",
//	"name": "credentials",
//	"description": "user's email/login & password'",
//	"schema": {
//		"type": "object",
//		"required":"password",
//		"properties": {
//			"email":    {"type": "string"},
//			"nickname": {"type": "string"},
//			"password": {"type": "string"}
//          ...
//		}
//	}
//
//}`)

var authenticateEndpoint = server_http.Endpoint{
	EndpointDescription: server_http.EndpointDescription{
		InternalKey: auth.IntefaceKeyAuthenticate,
		Method:      "POST",
	},

	//BodyParams: bodyParams,
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, _ server_http.PathParams, _ *auth.Identity) (server_http.Response, error) {

		credsJSON, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return server_http.ResponseRESTError(http.StatusBadRequest, errors.CommonError(common.WrongBodyKey, common.Map{"error": errors.Wrap(err, "can't read body")}), req)
		}

		var toAuth auth.Creds
		if err = json.Unmarshal(credsJSON, &toAuth); err != nil {
			return server_http.ResponseRESTError(http.StatusBadRequest, errors.CommonError(common.WrongJSONKey, common.Map{"error": errors.Wrapf(err, "can't unmarshal body: %s", credsJSON)}), req)
		}
		toAuth[auth.CredsIP] = req.RemoteAddr

		actor, err := authOp.Authenticate(toAuth)
		if err != nil {
			return server_http.ResponseRESTError(0, err, req)
		} else if actor == nil || actor.Identity == nil {
			return server_http.ResponseRESTError(0, auth.ErrNotAuthenticated, req)
		}

		toSet := auth.Creds{
			auth.CredsNickname: actor.Nickname,
			auth.CredsID:       string(actor.ID),
		}

		if len(actor.Roles) > 0 {
			rolesJSON, err := json.Marshal(actor.Roles)
			if err != nil {
				return server_http.ResponseRESTError(0, err, req)
			}
			toSet[auth.CredsRolesJSON] = string(rolesJSON)
		}

		jwtCreds, err := authJWTOp.SetCreds(auth.Actor{}, toSet)
		if err != nil || jwtCreds == nil {
			return server_http.ResponseRESTError(0, fmt.Errorf("got %#v / %s", jwtCreds, err), req)
		}
		actor.Creds = *jwtCreds

		return server_http.ResponseRESTOk(http.StatusOK, actor, req)
	},
}

var setCredsEndpoint = server_http.Endpoint{
	EndpointDescription: server_http.EndpointDescription{
		InternalKey: auth.IntefaceKeySetCreds,
		Method:      "POST",
	},

	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, _ server_http.PathParams, identity *auth.Identity) (server_http.Response, error) {

		credsJSON, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return server_http.ResponseRESTError(http.StatusBadRequest, errors.CommonError(common.WrongBodyKey, common.Map{"error": errors.Wrap(err, "can't read body")}), req)
		}

		var toSet auth.Creds
		if err = json.Unmarshal(credsJSON, &toSet); err != nil {
			return server_http.ResponseRESTError(http.StatusBadRequest, errors.CommonError(common.WrongJSONKey, common.Map{"error": errors.Wrapf(err, "can't unmarshal body: %s", credsJSON)}), req)
		}
		toSet[auth.CredsIP] = req.RemoteAddr

		creds, err := authOp.SetCreds(auth.Actor{Identity: identity}, toSet)
		if err != nil {
			return server_http.ResponseRESTError(0, err, req)
		}

		return server_http.ResponseRESTOk(http.StatusOK, creds, req)
	},
}
