package crud_serverhttp_jschmhr0

import (
	"encoding/json"
	"html"
	"net/http"
	"reflect"

	"github.com/julienschmidt/httprouter"
	"github.com/pavlo67/partes/confidenter/sessions"
	"github.com/pavlo67/partes/crud"
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/viewshtml"
	"github.com/pkg/errors"
)

const pageLengthDefault = 200

const onRead = "on crud_serverhttp_jschmhr.Read()"

func (rcOp *crud_serverhttp_jschmhr) Read(session *sessions.Item, r *http.Request, params httprouter.Params) (serverhttp.HTMLResponse, error) {
	defer r.Body.Close()
	crudType := params.ByName("type")

	crudOp, err := rcOp.getOp(crudType)
	if err != nil {
		return serverhttp.HTMLError(0, "невідомий тип даних"), errors.Wrap(err, onReadList)
	}

	id := params.ByName("id")

	item, err := crudOp.Read(session.UserIS(), id)
	if err != nil {
		return serverhttp.HTMLError(0, "не вдається прочитати дані"), errors.Wrapf(err, onRead)
	}

	description, err := crudOp.Describe()
	if err != nil {
		l.Errorf("can't crudOp.Describe(): %s", err)
	}

	data := crud.ReadData{
		Description: description,
		Item:        item,
	}

	dataJSON, _ := json.Marshal(data)

	return serverhttp.HTMLResponse{
		Data: map[string]string{
			"caput":  "??? 123",
			"corpus": string(dataJSON),
		},
	}, nil
}

const onReadList = "on crud_serverhttp_jschmhr.ReadList()"

func (rcOp *crud_serverhttp_jschmhr) ReadList(session *sessions.Item, r *http.Request, params httprouter.Params) (serverhttp.HTMLResponse, error) {
	defer r.Body.Close()
	crudType := params.ByName("type")

	crudOp, err := rcOp.getOp(crudType)
	if err != nil {
		return serverhttp.HTMLError(0, "невідомий тип даних"), err
	}

	description, err := crudOp.Describe()
	if err != nil {
		return serverhttp.HTMLError(0, "нема опису даних"), errors.Wrap(err, onReadList+": can't crudOp.Describe()")
	}

	options, optionsHTTP, err := content.ListOptionsFromRequest(r, pageLengthDefault, description.SortByDefault)
	if err != nil || options == nil {
		return serverhttp.HTMLError(0, "не вдається прочитати параметри запиту"), errors.Wrap(err, onReadList)
	}

	var items []interface{}
	items, optionsHTTP.AllCnt, err = crudOp.ReadList(session.UserIS(), options)
	if err != nil {
		return serverhttp.HTMLError(0, "не вдається прочитати дані"), errors.Wrap(err, onReadList)
	}

	table := crud.Table(description, items, crudOp)
	pagination := viewshtml.Pagination(options.Limits, options.SortBy, optionsHTTP)
	save := `<button id="saveEditedInTable" data-crud_type="` + html.EscapeString(crudType) + `">зберегти зміни</button>`
	return serverhttp.HTMLResponse{
		Data: map[string]string{
			"caput":  crudType,
			"corpus": table + "<p>" + save + "<p>" + pagination,
		},
	}, nil
}

type UpdateListREST struct {
	ID    map[string]string `json:"id"`
	Key   string            `json:"key"`
	Value string            `json:"value,omitempty"`
}

const onUpdateListREST = "on crud_serverhttp_jschmhr.UpdateListREST()"

func (rcOp *crud_serverhttp_jschmhr) UpdateListREST(session *sessions.Item, r *http.Request, params httprouter.Params) (serverhttp.RESTResponse, error) {
	defer r.Body.Close()
	crudType := params.ByName("type")

	var data []UpdateListREST
	err := json.NewDecoder(r.Body).Decode(&data)

	l.Info(crudType, data, err)

	//
	//crudOp, err := rcOp.getOp(crudType)
	//if err != nil {
	//	return serverhttp.RESTError(basis.ErrCantPerform), err
	//}
	//
	//description, err := crudOp.Describe()
	//if err != nil {
	//	return serverhttp.RESTError(basis.ErrCantPerform), err
	//}

	//data := reflect.New(reflect.ValueOf(description.Exemplar).Elem().Type()).Interface()
	//err = json.NewDecoder(r.Body).Decode(&data)
	//if err != nil {
	//	return serverhttp.RESTError(basis.ErrJSONFormat), errors.Wrap(err, "on decode r.Body")
	//}
	//
	//id, err := crudOp.IDFromNative(data)
	//if err != nil {
	//	return serverhttp.RESTError(basis.ErrJSONFormat), errors.Wrapf(err, "on crudOp.IDFromNative(%#v)", data)
	//}
	//
	//var info string
	//if id == "" {
	//	id, err = crudOp.Create(user.Identity().String(), data)
	//	info = `Створено запис (` + id + `)`
	//
	//} else {
	//	var res crud.Result
	//	res, err = crudOp.Update(user.Identity().String(), data)
	//	if res.NumOk < 1 {
	//		info = `Запис (` + id + `) не зазнав змін.`
	//	} else {
	//		info = `Запис (` + id + `) оновлено.`
	//	}
	//}
	//if err != nil {
	//	return serverhttp.RESTError(basis.ErrCantPerform), err
	//}
	//
	return serverhttp.RESTResponse{}, nil
}

const onSaveREST = "on crud_serverhttp_jschmhr.SaveREST()"

func (rcOp *crud_serverhttp_jschmhr) SaveREST(session *sessions.Item, r *http.Request, params httprouter.Params) (serverhttp.RESTResponse, error) {
	defer r.Body.Close()
	crudType := params.ByName("type")

	crudOp, err := rcOp.getOp(crudType)
	if err != nil {
		return serverhttp.RESTError(basis.ErrCantPerform), err
	}

	description, err := crudOp.Describe()
	if err != nil {
		return serverhttp.RESTError(basis.ErrCantPerform), errors.Wrap(err, onSaveREST)
	}

	data := reflect.New(reflect.ValueOf(description.Exemplar).Elem().Type()).Interface()
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return serverhttp.RESTError(basis.ErrJSONFormat), errors.Wrap(err, "on decode r.Body")
	}

	id, err := crudOp.IDFromNative(data)
	if err != nil {
		return serverhttp.RESTError(basis.ErrJSONFormat), errors.Wrapf(err, "on crudOp.IDFromNative(%#v)", data)
	}

	var info string
	if id == "" {
		id, err = crudOp.Create(session.UserIS(), data)
		info = `Створено запис (` + id + `)`

	} else {
		var res crud.Result
		res, err = crudOp.Update(session.UserIS(), data)
		if res.NumOk < 1 {
			info = `Запис (` + id + `) не зазнав змін.`
		} else {
			info = `Запис (` + id + `) оновлено.`
		}
	}
	if err != nil {
		return serverhttp.RESTError(basis.ErrCantPerform), err
	}

	return serverhttp.RESTResponse{
		Data: crud.ResultData{
			IDs:  []string{id},
			Info: info,
		},
	}, nil

}

const onDeleteREST = "on crud_serverhttp_jschmhr.DeleteREST()"

func (rcOp *crud_serverhttp_jschmhr) DeleteREST(session *sessions.Item, r *http.Request, params httprouter.Params) (serverhttp.RESTResponse, error) {
	defer r.Body.Close()
	crudType := params.ByName("type")

	crudOp, err := rcOp.getOp(crudType)
	if err != nil {
		return serverhttp.RESTError(basis.ErrCantPerform), err
	}

	id := params.ByName("id")

	res, err := crudOp.Delete(session.UserIS(), id)
	if err != nil {
		return serverhttp.RESTError(basis.ErrCantPerform), err
	}

	var info string
	if res.NumOk < 1 {
		info = `Запис (` + id + `) не зазнав змін.`
	} else {
		info = `Запис (` + id + `) вилучено.`
	}

	return serverhttp.RESTResponse{
		Data: crud.ResultData{
			IDs:  []string{id},
			Info: info,
		},
	}, nil
}

const onReadREST = "on crud_serverhttp_jschmhr.ReadREST()"

func (rcOp *crud_serverhttp_jschmhr) ReadREST(session *sessions.Item, r *http.Request, params httprouter.Params) (serverhttp.RESTResponse, error) {
	defer r.Body.Close()
	crudType := params.ByName("type")

	crudOp, err := rcOp.getOp(crudType)
	if err != nil {
		return serverhttp.RESTError(basis.ErrCantPerform), err
	}

	id := params.ByName("id")

	item, err := crudOp.Read(session.UserIS(), id)
	if err != nil {
		return serverhttp.RESTError(basis.ErrCantPerform), err
	}

	description, err := crudOp.Describe()
	if err != nil {
		l.Errorf("can't crudOp.Describe(): %s", err)
	}

	return serverhttp.RESTResponse{
		Data: crud.ReadData{
			Description: description,
			Item:        item,
		},
	}, nil
}

const onReadListREST = "on crud_serverhttp_jschmhr.ReadListREST()"

func (rcOp *crud_serverhttp_jschmhr) ReadListREST(session *sessions.Item, r *http.Request, params httprouter.Params) (serverhttp.RESTResponse, error) {
	defer r.Body.Close()
	crudType := params.ByName("type")

	crudOp, err := rcOp.getOp(crudType)
	if err != nil {
		return serverhttp.RESTError(basis.ErrCantPerform), err
	}

	description, err := crudOp.Describe()
	if err != nil {
		l.Errorf("can't crudOp.Describe(): %s", err)
	}

	options, _, err := content.ListOptionsFromRequest(r, pageLengthDefault, description.SortByDefault)
	if err != nil || options == nil {
		return serverhttp.RESTError(errors.New("не вдається прочитати параметри запиту")), err
	}

	items, allCnt, err := crudOp.ReadList(session.UserIS(), options)
	if err != nil {
		return serverhttp.RESTError(basis.ErrCantPerform), err
	}

	return serverhttp.RESTResponse{
		Data: crud.ReadListData{
			Description: description,
			Items:       items,
			AllCnt:      allCnt,
		},
	}, nil
}
