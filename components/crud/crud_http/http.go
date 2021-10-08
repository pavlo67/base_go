package crud_http

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/httplib"
	"github.com/pavlo67/common/common/server/server_http"

	"github.com/pavlo67/data/elements/selectors"

	"github.com/pavlo67/data/components/crud"
)

var _ crud.Operator = &crudHTTP{}

type crudHTTP struct {
	serverConfig server_http.Config
}

const onNew = "on crudHTTP.New()"

func New(serverConfig server_http.Config) (crud.Operator, error) {
	crudOp := crudHTTP{
		serverConfig: serverConfig,
	}

	return &crudOp, nil
}

const onTypes = "on crudHTTP.Types()"

func (crudOp *crudHTTP) Types() ([]crud.Type, error) {
	ep := crudOp.serverConfig.EndpointsSettled[crud.IntefaceKeyTypes]
	serverURL := crudOp.serverConfig.Host + crudOp.serverConfig.Port + crudOp.serverConfig.Prefix + ep.Path

	// TODO!!!
	var creds *auth.Creds

	var types []crud.Type
	if err := httplib.Request(nil, serverURL, ep.Method, server_http.SetCreds(creds), nil, &types, l); err != nil {
		return nil, errors.Wrap(err, onTypes)
	}

	return types, nil
}

const onSave = "on crudHTTP.Save()"

func (crudOp *crudHTTP) Save(data crud.Data, _ *auth.Identity) (*crud.Key, error) {
	ep := crudOp.serverConfig.EndpointsSettled[crud.IntefaceKeySave]
	serverURL := crudOp.serverConfig.Host + crudOp.serverConfig.Port + crudOp.serverConfig.Prefix + ep.Path

	requestBody, err := json.Marshal(data)
	if err != nil {
		return nil, errors.Wrapf(err, onSave+": can't marshal data (%#v)", data)
	}

	// TODO!!!
	var creds *auth.Creds

	var key *crud.Key
	if err := httplib.Request(nil, serverURL, ep.Method, server_http.SetCreds(creds), requestBody, &key, l); err != nil {
		return nil, errors.Wrap(err, onSave)
	}

	return key, nil
}

const onRead = "on crudHTTP.Read()"

func (crudOp *crudHTTP) Read(key crud.Key, _ *auth.Identity) (*crud.Data, error) {
	ep := crudOp.serverConfig.EndpointsSettled[crud.IntefaceKeyRead]
	serverURL := crudOp.serverConfig.Host + crudOp.serverConfig.Port + crudOp.serverConfig.Prefix + ep.Path +
		"/" + string(key.Type) + "/" + key.ID.String()

	// TODO!!!
	var creds *auth.Creds

	var item crud.Data
	if err := httplib.Request(nil, serverURL, ep.Method, server_http.SetCreds(creds), nil, &item, l); err != nil {
		return nil, errors.Wrap(err, onRead)
	}

	return &item, nil
}

const onList = "on crudHTTP.List()"

func (crudOp *crudHTTP) List(crudType crud.Type, options selectors.Options, identity *auth.Identity) ([]crud.Data, error) {
	ep := crudOp.serverConfig.EndpointsSettled[crud.IntefaceKeyList]
	serverURL := crudOp.serverConfig.Host + crudOp.serverConfig.Port + crudOp.serverConfig.Prefix + ep.Path +
		"/" + string(crudType)

	// TODO!!! add selector too
	var creds *auth.Creds

	var items []crud.Data
	if err := httplib.Request(nil, serverURL, ep.Method, server_http.SetCreds(creds), nil, &items, l); err != nil {
		return nil, errors.Wrap(err, onList)
	}

	return items, nil
}

const onRemove = "on crudHTTP.Remove()"

func (crudOp *crudHTTP) Remove(key crud.Key, identity *auth.Identity) error {
	ep := crudOp.serverConfig.EndpointsSettled[crud.IntefaceKeyRemove]
	serverURL := crudOp.serverConfig.Host + crudOp.serverConfig.Port + crudOp.serverConfig.Prefix + ep.Path +
		"/" + string(key.Type) + "/" + key.ID.String()

	// TODO!!!
	var creds *auth.Creds

	if err := httplib.Request(nil, serverURL, ep.Method, server_http.SetCreds(creds), nil, nil, l); err != nil {
		return errors.Wrap(err, onRemove)
	}

	return nil
}
