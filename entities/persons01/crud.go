package persons01

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/pavlo67/data/components/crud"

	"github.com/pavlo67/data/entities"

	"github.com/pavlo67/data/common/auth"

	"github.com/pavlo67/data/elements/selectors"
)

var _ crud.Operator = &persons01CRUD{}

func OperatorCRUD(personsOp Operator) (crud.Operator, error) {
	if personsOp == nil {
		return nil, errors.New("personsOp == nil")
	}

	return &persons01CRUD{personsOp: personsOp}, nil
}

const CRUD01 crud.Type = "persons01"

type persons01CRUD struct {
	personsOp Operator
}

func (crudOp *persons01CRUD) Types() ([]crud.Type, error) {
	return []crud.Type{CRUD01}, nil
}

const onSave = "on persons01/crud.Save()"

func (crudOp *persons01CRUD) Save(data crud.Data, actor auth.Actor) (*crud.Key, error) {
	if data.Key.Type != CRUD01 {
		return nil, fmt.Errorf(onSave+": wrong key.Type (%#v) to save item (%#v)", data.Key, data.Value)
	}

	var item Item

	switch v := data.Value.(type) {
	case Item:
		item = v
	case *Item:
		if v == nil {
			return nil, errors.New(onSave + ": nil Item to save")
		}
		item = *v
	case json.RawMessage:
		if err := json.Unmarshal(v, &item.Person01); err != nil {
			return nil, fmt.Errorf(onSave+": can't unmarshal (%s) into item.Person01", v)
		}
	case entities.Person01:
		item = Item{Person01: v}
	case *entities.Person01:
		if v == nil {
			return nil, errors.New(onSave + ": nil Person01 to save")
		}
		item = Item{Person01: *v}
	default:
		return nil, fmt.Errorf(onSave+": wrong data (%#v) to save with key (%#v)", data.Value, data.Key)
	}

	item.ID = data.Key.ID
	item.Description = data.Description
	id, err := crudOp.personsOp.Save(item, actor)
	if err != nil {
		return nil, errors.Wrap(err, onSave)
	}

	return &crud.Key{Type: CRUD01, ID: id}, nil
}

const onRead = "on persons01/crud.Read()"

func (crudOp *persons01CRUD) Read(key crud.Key, actor auth.Actor) (*crud.Data, error) {
	if key.Type != CRUD01 {
		return nil, fmt.Errorf(onRead+": wrong key.Type (%#v)", key)
	}

	item, err := crudOp.personsOp.Read(key.ID, actor)
	if err != nil || item == nil {
		return nil, fmt.Errorf(onRead+": got %#v / %s", item, err)
	}

	return &crud.Data{
		Key: crud.Key{
			Type: CRUD01,
			ID:   key.ID,
		},
		Description: item.Description,
		Value:       item.Person01,
	}, nil
}

const onList = "on persons01/crud.List()"

func (crudOp *persons01CRUD) List(crudType crud.Type, _ selectors.Options, actor auth.Actor) ([]crud.Data, error) {
	if crudType != CRUD01 {
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
				Type: CRUD01,
				ID:   pi.ID,
			},
			Description: pi.Description,
			Value:       pi.Person01,
		}
	}

	return crudItems, nil

}

const onRemove = "on persons01/crud.Remove()"

func (crudOp *persons01CRUD) Remove(key crud.Key, actor auth.Actor) error {
	if key.Type != CRUD01 {
		return fmt.Errorf(onRemove+": wrong key.Type (%#v)", key)
	}

	return crudOp.personsOp.Remove(key.ID, actor)
}
