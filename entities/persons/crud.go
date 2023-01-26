package persons

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/rbac"

	"github.com/pavlo67/data/components/crud"
	"github.com/pavlo67/data/components/vcs"
)

var _ crud.Operator = &personsCRUD{}

func OperatorCRUD(personsOp Operator, roles rbac.Roles) (crud.Operator, error) {
	if personsOp == nil {
		return nil, errors.New("personsOp == nil")
	}

	return &personsCRUD{personsOp: personsOp, roles: roles}, nil
}

const CRUD crud.Type = "persons"

type personsCRUD struct {
	personsOp Operator
	roles     rbac.Roles
}

func (crudOp *personsCRUD) Types() ([]crud.Type, error) {
	return []crud.Type{CRUD}, nil
}

func (crudOp *personsCRUD) Roles() (rbac.Roles, error) {
	return crudOp.roles, nil
}

const onSave = "on persons01/crud.Save()"

func (crudOp *personsCRUD) Save(data crud.Data, actor auth.Actor) (*crud.Key, vcs.History, error) {
	if data.Key.Type != CRUD {
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
		if err := json.Unmarshal(v, &item.Person); err != nil {
			return nil, nil, fmt.Errorf(onSave+": can't unmarshal (%s) into item.Person", v)
		}
	case Person:
		item = Item{Person: v}
	case *Person:
		if v == nil {
			return nil, nil, errors.New(onSave + ": nil Person to save")
		}
		item = Item{Person: *v}
	default:
		return nil, nil, fmt.Errorf(onSave+": wrong data (%#v) to save with key (%#v)", data.Value, data.Key)
	}

	item.ID = data.Key.ID
	item.Description = data.Description
	id, _, historyChanged, err := crudOp.personsOp.Save(item, actor)
	if err != nil {
		return nil, nil, errors.Wrap(err, onSave)
	}

	return &crud.Key{Type: CRUD, ID: id}, historyChanged, nil
}

const onRead = "on persons01/crud.Read()"

func (crudOp *personsCRUD) Read(key crud.Key, actor auth.Actor) (*crud.Data, error) {
	if key.Type != CRUD {
		return nil, fmt.Errorf(onRead+": wrong key.Type (%#v)", key)
	}

	item, err := crudOp.personsOp.Read(key.ID, actor)
	if err != nil || item == nil {
		return nil, fmt.Errorf(onRead+": got %#v / %s", item, err)
	}

	return &crud.Data{
		Key: crud.Key{
			Type: CRUD,
			ID:   key.ID,
		},
		Description: item.Description,
		Value:       item.Person,
	}, nil
}

const onList = "on persons01/crud.List()"

func (crudOp *personsCRUD) List(crudType crud.Type, _ crud.Options, actor auth.Actor) ([]crud.Data, error) {
	if crudType != CRUD {
		return nil, fmt.Errorf(onList+": wrong crudType (%#v)", crudType)
	}

	// TODO!!! use selector
	items, err := crudOp.personsOp.List(nil, actor)
	if err != nil {
		return nil, errors.Wrap(err, onList)
	}

	crudItems := make([]crud.Data, len(items))
	for i, pi := range items {
		crudItems[i] = crud.Data{
			Key: crud.Key{
				Type: CRUD,
				ID:   pi.ID,
			},
			Description: pi.Description,
			Value:       pi.Person,
		}
	}

	return crudItems, nil

}

const onRemove = "on persons01/crud.Remove()"

func (crudOp *personsCRUD) Remove(key crud.Key, actor auth.Actor) error {
	if key.Type != CRUD {
		return fmt.Errorf(onRemove+": wrong key.Type (%#v)", key)
	}

	return crudOp.personsOp.Remove(key.ID, actor)
}
