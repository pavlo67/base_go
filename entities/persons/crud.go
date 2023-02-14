package persons

import (
	"encoding/json"
	"fmt"

	"github.com/pavlo67/data/entities"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/rbac"

	"github.com/pavlo67/data/components/vcs"
)

var _ entities.OperatorCRUD = &personsCRUD{}

func OperatorCRUD(personsOp Operator, roles rbac.Roles) (entities.OperatorCRUD, error) {
	if personsOp == nil {
		return nil, errors.New("personsOp == nil")
	}

	return &personsCRUD{personsOp: personsOp, roles: roles}, nil
}

const CRUD entities.Type = "persons"

type personsCRUD struct {
	personsOp Operator
	roles     rbac.Roles
}

func (crudOp *personsCRUD) Types() ([]entities.Type, error) {
	return []entities.Type{CRUD}, nil
}

func (crudOp *personsCRUD) Roles() (rbac.Roles, error) {
	return crudOp.roles, nil
}

const onSave = "on persons01/crud.Add()"

func (crudOp *personsCRUD) Save(data entities.Data, actor auth.Actor) (*entities.Key, vcs.History, error) {
	if data.Key.Type != CRUD {
		return nil, nil, fmt.Errorf(onSave+": wrong key.ImporterInterfaceKey (%#v) to save item (%#v)", data.Key, data.Value)
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

	return &entities.Key{Type: CRUD, ID: id}, historyChanged, nil
}

const onRead = "on persons01/crud.Read()"

func (crudOp *personsCRUD) Read(key entities.Key, actor auth.Actor) (*entities.Data, error) {
	if key.Type != CRUD {
		return nil, fmt.Errorf(onRead+": wrong key.ImporterInterfaceKey (%#v)", key)
	}

	item, err := crudOp.personsOp.Read(key.ID, actor)
	if err != nil || item == nil {
		return nil, fmt.Errorf(onRead+": got %#v / %s", item, err)
	}

	return &entities.Data{
		Key: entities.Key{
			Type: CRUD,
			ID:   key.ID,
		},
		Description: item.Description,
		Value:       item.Person,
	}, nil
}

const onList = "on persons01/crud.List()"

func (crudOp *personsCRUD) List(crudType entities.Type, _ entities.Options, actor auth.Actor) ([]entities.Data, error) {
	if crudType != CRUD {
		return nil, fmt.Errorf(onList+": wrong crudType (%#v)", crudType)
	}

	// TODO!!! use selector
	items, err := crudOp.personsOp.List(nil, actor)
	if err != nil {
		return nil, errors.Wrap(err, onList)
	}

	crudItems := make([]entities.Data, len(items))
	for i, pi := range items {
		crudItems[i] = entities.Data{
			Key: entities.Key{
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

func (crudOp *personsCRUD) Remove(key entities.Key, actor auth.Actor) error {
	if key.Type != CRUD {
		return fmt.Errorf(onRemove+": wrong key.ImporterInterfaceKey (%#v)", key)
	}

	return crudOp.personsOp.Remove(key.ID, actor)
}
