package records01

import (
	"encoding/json"
	"fmt"

	"github.com/pavlo67/data/components/vcs"

	"github.com/pavlo67/data/components/selectors"

	"github.com/pavlo67/common/common/rbac"
	"github.com/pavlo67/data/components/crud"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/auth"

	"github.com/pavlo67/data/entities"
)

const CRUD01 crud.Type = "records01"

var _ crud.Operator = &records01CRUD{}

func OperatorCRUD(recordsOp Operator, roles rbac.Roles) (crud.Operator, error) {
	if recordsOp == nil {
		return nil, errors.New("recordsOp == nil")
	}

	return &records01CRUD{recordsOp: recordsOp, roles: roles}, nil
}

type records01CRUD struct {
	recordsOp Operator
	roles     rbac.Roles
}

func (crudOp *records01CRUD) Types() ([]crud.Type, error) {
	return []crud.Type{CRUD01}, nil
}

func (crudOp *records01CRUD) Roles() (rbac.Roles, error) {
	return crudOp.roles, nil
}

const onSave = "on records01/crud.Save()"

func (crudOp *records01CRUD) Save(data crud.Data, actor auth.Actor) (*crud.Key, vcs.History, error) {
	if data.Key.Type != CRUD01 {
		return nil, nil, fmt.Errorf(onSave+": wrong key.Type (%#v) to save item (%#v)", data.Key, data.Value)
	}

	var item Item

	switch v := data.Value.(type) {
	case Item:
		item = v
	case *Item:
		if v == nil {
			return nil, nil, errors.New(onSave + ": nil Item to save")
		}
		item = *v
	case json.RawMessage:
		if err := json.Unmarshal(v, &item.Record01); err != nil {
			return nil, nil, fmt.Errorf(onSave+": can't unmarshal (%s) into item.Record01", v)
		}
	case entities.Record01:
		item = Item{Record01: v}
	case *entities.Record01:
		if v == nil {
			return nil, nil, errors.New(onSave + ": nil Record01 to save")
		}
		item = Item{Record01: *v}
	default:
		return nil, nil, fmt.Errorf(onSave+": wrong data (%#v) to save with key (%#v)", data.Value, data.Key)
	}

	item.ID = data.Key.ID
	item.Description = data.Description
	id, historyChangedStr, err := crudOp.recordsOp.Save(item, actor)
	if err != nil {
		return nil, nil, errors.Wrap(err, onSave)
	}

	return &crud.Key{Type: CRUD01, ID: id}, historyChangedStr, nil
}

const onRead = "on records01/crud.Read()"

func (crudOp *records01CRUD) Read(key crud.Key, actor auth.Actor) (*crud.Data, error) {
	if key.Type != CRUD01 {
		return nil, fmt.Errorf(onRead+": wrong key.Type (%#v)", key)
	}

	item, err := crudOp.recordsOp.Read(key.ID, actor)
	if err != nil || item == nil {
		return nil, fmt.Errorf(onRead+": got %#v / %s", item, err)
	}

	return &crud.Data{
		Key: crud.Key{
			Type: CRUD01,
			ID:   key.ID,
		},
		Description: item.Description,
		Value:       item.Record01,
	}, nil
}

const onList = "on records01/crud.List()"

func (crudOp *records01CRUD) List(crudType crud.Type, _ selectors.Options, actor auth.Actor) ([]crud.Data, error) {
	if crudType != CRUD01 {
		return nil, fmt.Errorf(onList+": wrong crudType (%#v)", crudType)
	}

	// TODO!!! use selector
	items, err := crudOp.recordsOp.List(nil, actor)
	if err != nil {
		return nil, errors.Wrap(err, onList)
	}

	crudItems := make([]crud.Data, len(items))
	for i, pi := range items {
		crudItems[i] = crud.Data{
			Key: crud.Key{
				Type: CRUD01,
				ID:   pi.ID,
			},
			Description: pi.Description,
			Value:       pi.Record01,
		}
	}

	return crudItems, nil

}

const onRemove = "on records01/crud.Remove()"

func (crudOp *records01CRUD) Remove(key crud.Key, actor auth.Actor) error {
	if key.Type != CRUD01 {
		return fmt.Errorf(onRemove+": wrong key.Type (%#v)", key)
	}

	return crudOp.recordsOp.Remove(key.ID, actor)
}
