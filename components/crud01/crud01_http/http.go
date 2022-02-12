package crud01_http

import (
	"encoding/json"

	"github.com/pavlo67/data/components/crud01/crud01_server_http"

	"github.com/pavlo67/data/components/vcs"

	"github.com/pavlo67/data/components/selectors"

	"github.com/pavlo67/data/components/crud"

	"github.com/pavlo67/common/common/rbac"

	"github.com/pavlo67/common/common"

	"github.com/pavlo67/common/common/auth"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/httplib"
	"github.com/pavlo67/common/common/server/server_http"
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

func (crudOp *crudHTTP) Roles() (rbac.Roles, error) {
	return nil, common.ErrNotSupported
}

const onTypes = "on crudHTTP.Types()"

func (crudOp *crudHTTP) Types() ([]crud.Type, error) {
	ep := crudOp.serverConfig.EndpointsSettled[crud.IntefaceKeyTypes]
	serverURL := crudOp.serverConfig.Host + crudOp.serverConfig.Port + crudOp.serverConfig.Prefix + ep.Path

	var types []crud.Type
	if err := httplib.Request(nil, serverURL, ep.Method, server_http.SetCreds(nil), nil, &types, l); err != nil {
		return nil, errors.Wrap(err, onTypes)
	}

	return types, nil
}

const onSave = "on crudHTTP.Save()"

func (crudOp *crudHTTP) Save(data crud.Data, actor auth.Actor) (*crud.Key, vcs.History, error) {
	ep := crudOp.serverConfig.EndpointsSettled[crud.IntefaceKeySave]
	serverURL := crudOp.serverConfig.Host + crudOp.serverConfig.Port + crudOp.serverConfig.Prefix + ep.Path

	requestBody, err := json.Marshal(data)
	if err != nil {
		return nil, nil, errors.Wrapf(err, onSave+": can't marshal data (%#v)", data)
	}

	var saveResult crud01_server_http.SaveResult
	if err := httplib.Request(nil, serverURL, ep.Method, server_http.SetCreds(actor.Creds), requestBody, &saveResult, l); err != nil {
		return nil, nil, errors.Wrap(err, onSave)
	}

	return &saveResult.Key, saveResult.History, nil
}

const onRead = "on crudHTTP.Read()"

func (crudOp *crudHTTP) Read(key crud.Key, actor auth.Actor) (*crud.Data, error) {
	ep := crudOp.serverConfig.EndpointsSettled[crud.IntefaceKeyRead]
	serverURL := crudOp.serverConfig.Host + crudOp.serverConfig.Port + crudOp.serverConfig.Prefix + ep.Path +
		"/" + string(key.Type) + "/" + key.ID.String()

	var dataRaw crud.DataRaw
	if err := httplib.Request(nil, serverURL, ep.Method, server_http.SetCreds(actor.Creds), nil, &dataRaw, l); err != nil {
		return nil, errors.Wrap(err, onRead)
	}

	return &crud.Data{dataRaw.Key, dataRaw.Description, dataRaw.Value}, nil
}

const onList = "on crudHTTP.List()"

func (crudOp *crudHTTP) List(crudType crud.Type, options selectors.Options, actor auth.Actor) ([]crud.Data, error) {
	ep := crudOp.serverConfig.EndpointsSettled[crud.IntefaceKeyList]
	serverURL := crudOp.serverConfig.Host + crudOp.serverConfig.Port + crudOp.serverConfig.Prefix + ep.Path +
		"/" + string(crudType)

	var dataRaws []crud.DataRaw
	if err := httplib.Request(nil, serverURL, ep.Method, server_http.SetCreds(actor.Creds), nil, &dataRaws, l); err != nil {
		return nil, errors.Wrap(err, onList)
	}

	items := make([]crud.Data, len(dataRaws))
	for i, dataRaw := range dataRaws {
		items[i] = crud.Data{dataRaw.Key, dataRaw.Description, dataRaw.Value}
	}

	return items, nil
}

const onRemove = "on crudHTTP.Remove()"

func (crudOp *crudHTTP) Remove(key crud.Key, actor auth.Actor) error {
	ep := crudOp.serverConfig.EndpointsSettled[crud.IntefaceKeyRemove]
	serverURL := crudOp.serverConfig.Host + crudOp.serverConfig.Port + crudOp.serverConfig.Prefix + ep.Path +
		"/" + string(key.Type) + "/" + key.ID.String()

	if err := httplib.Request(nil, serverURL, ep.Method, server_http.SetCreds(actor.Creds), nil, nil, l); err != nil {
		return errors.Wrap(err, onRemove)
	}

	return nil
}
