package crud01_server_http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/pavlo67/data/components/vcs"

	"github.com/pavlo67/data/components/selectors"

	"github.com/pavlo67/data/components/crud"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/server"
	"github.com/pavlo67/common/common/server/server_http"
)

var Endpoints = server_http.Endpoints{
	typesEndpoint,
	saveEndpoint,
	readEndpoint,
	listEndpoint,
	removeEndpoint,
}

var typesEndpoint = server_http.Endpoint{
	EndpointDescription: server_http.EndpointDescription{
		InternalKey: crud.IntefaceKeyTypes,
		Method:      "DELETE",
		PathParams:  []string{"type", "id"},
	},

	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server.Response, error) {
		types, err := crudDispatcherOp.Types()
		if err != nil {
			return server_http.ResponseRESTError(0, err, req)
		}

		return server_http.ResponseRESTOk(http.StatusOK, types, req)
	},
}

type SaveResult struct {
	crud.Key
	vcs.History
}

var saveEndpoint = server_http.Endpoint{
	EndpointDescription: server_http.EndpointDescription{
		InternalKey: crud.IntefaceKeySave,
		Method:      "POST",
	},

	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, _ server_http.PathParams, identity *auth.Identity) (server.Response, error) {
		dataJSON, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return server_http.ResponseRESTError(http.StatusBadRequest, errors.CommonError(common.WrongBodyKey, common.Map{"error": errors.Wrap(err, "can't read body")}), req)
		}
		var dataRaw crud.DataRaw
		if err = json.Unmarshal(dataJSON, &dataRaw); err != nil {
			return server_http.ResponseRESTError(http.StatusBadRequest, errors.CommonError(common.WrongJSONKey, common.Map{"error": errors.Wrapf(err, "can't unmarshal body: %s", dataJSON)}), req)
		}

		key, historyChanged, err := crudDispatcherOp.Save(crud.Data{dataRaw.Key, dataRaw.Description, dataRaw.Value}, auth.Actor{Identity: identity})
		if err != nil {
			return server_http.ResponseRESTError(0, err, req)
		} else if key == nil {
			return server_http.ResponseRESTError(0, fmt.Errorf("key == nil"), req)
		}

		return server_http.ResponseRESTOk(http.StatusOK, SaveResult{Key: *key, History: historyChanged}, req)
	},
}

var readEndpoint = server_http.Endpoint{
	EndpointDescription: server_http.EndpointDescription{
		InternalKey: crud.IntefaceKeyRead,
		Method:      "GET",
		PathParams:  []string{"type", "id"},
	},

	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server.Response, error) {
		key := crud.Key{
			Type: crud.Type(params["type"]),
			ID:   crud.NewID(params["id"]),
		}

		item, err := crudDispatcherOp.Read(key, auth.Actor{Identity: identity})
		if err != nil {
			return server_http.ResponseRESTError(0, err, req)
		}

		return server_http.ResponseRESTOk(http.StatusOK, item, req)
	},
}

var listEndpoint = server_http.Endpoint{
	EndpointDescription: server_http.EndpointDescription{
		InternalKey: crud.IntefaceKeyList,
		Method:      "GET",
		PathParams:  []string{"type"},
	},

	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server.Response, error) {
		crudType := crud.Type(params["type"])

		// TODO!!!
		var selector selectors.Options

		items, err := crudDispatcherOp.List(crudType, selector, auth.Actor{Identity: identity})
		if err != nil {
			return server_http.ResponseRESTError(0, err, req)
		}

		return server_http.ResponseRESTOk(http.StatusOK, items, req)
	},
}

var removeEndpoint = server_http.Endpoint{
	EndpointDescription: server_http.EndpointDescription{
		InternalKey: crud.IntefaceKeyRemove,
		Method:      "DELETE",
		PathParams:  []string{"type", "id"},
	},

	WorkerHTTP: func(serverOp server_http.Operator, req *http.Request, params server_http.PathParams, identity *auth.Identity) (server.Response, error) {
		key := crud.Key{
			Type: crud.Type(params["type"]),
			ID:   crud.NewID(params["id"]),
		}

		if err := crudDispatcherOp.Remove(key, auth.Actor{Identity: identity}); err != nil {
			return server_http.ResponseRESTError(0, err, req)
		}

		return server_http.ResponseRESTOk(http.StatusOK, nil, req)
	},
}
