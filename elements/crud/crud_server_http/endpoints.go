package crud_server_http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/server"
	"github.com/pavlo67/common/common/server/server_http"
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
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, _ server_http.PathParams, _ *auth.Identity) (server.Response, error) {

		credsJSON, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return server_http.ResponseRESTError(http.StatusBadRequest, errors.CommonError(common.WrongBodyKey, common.Map{"error": errors.Wrap(err, "can't read body")}), req)
		}

		var toAuth auth.Creds
		if err = json.Unmarshal(credsJSON, &toAuth); err != nil {
			return server_http.ResponseRESTError(http.StatusBadRequest, errors.CommonError(common.WrongJSONKey, common.Map{"error": errors.Wrapf(err, "can't unmarshal body: %s", credsJSON)}), req)
		}
		toAuth[auth.CredsIP] = req.RemoteAddr

		identity, err := authOp.Authenticate(toAuth)
		if err != nil {
			return server_http.ResponseRESTError(0, err, req)
		}

		return server_http.ResponseRESTOk(http.StatusOK, identity, req)
	},
}

var setCredsEndpoint = server_http.Endpoint{
	EndpointDescription: server_http.EndpointDescription{
		InternalKey: auth.IntefaceKeySetCreds,
		Method:      "POST",
	},

	//BodyParams: bodyParams,
	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, _ server_http.PathParams, identity *auth.Identity) (server.Response, error) {

		credsJSON, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return server_http.ResponseRESTError(http.StatusBadRequest, errors.CommonError(common.WrongBodyKey, common.Map{"error": errors.Wrap(err, "can't read body")}), req)
		}

		var toSet auth.Creds
		if err = json.Unmarshal(credsJSON, &toSet); err != nil {
			return server_http.ResponseRESTError(http.StatusBadRequest, errors.CommonError(common.WrongJSONKey, common.Map{"error": errors.Wrapf(err, "can't unmarshal body: %s", credsJSON)}), req)
		}
		toSet[auth.CredsIP] = req.RemoteAddr

		var authID auth.ID
		if identity != nil {
			authID = identity.ID
		}

		creds, err := authOp.SetCreds(authID, toSet)
		if err != nil {
			return server_http.ResponseRESTError(0, err, req)
		}

		return server_http.ResponseRESTOk(http.StatusOK, creds, req)
	},
}
