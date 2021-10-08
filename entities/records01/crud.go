package records01

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/auth"

	"github.com/pavlo67/data/elements/selectors"

	"github.com/pavlo67/data/entities"

	"github.com/pavlo67/data/components/crud"
)

var _ crud.Operator = &records01CRUD{}

func OperatorCRUD(recordsOp Operator) (crud.Operator, error) {
	if recordsOp == nil {
		return nil, errors.New("recordsOp == nil")
	}

	return &records01CRUD{recordsOp: recordsOp}, nil
}

const CRUD01 crud.Type = "records01"

type records01CRUD struct {
	recordsOp Operator
}

func (crudOp *records01CRUD) Types() ([]crud.Type, error) {
	return []crud.Type{CRUD01}, nil
}

const onSave = "on records01/crud.Save()"

func (crudOp *records01CRUD) Save(data crud.Data, identity *auth.Identity) (*crud.Key, error) {
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
		if err := json.Unmarshal(v, &item.Record01); err != nil {
			return nil, fmt.Errorf(onSave+": can't unmarshal (%s) into item.Record01", v)
		}
	case entities.Record01:
		item = Item{Record01: v}
	case *entities.Record01:
		if v == nil {
			return nil, errors.New(onSave + ": nil Record01 to save")
		}
		item = Item{Record01: *v}
	default:
		return nil, fmt.Errorf(onSave+": wrong data (%#v) to save with key (%#v)", data.Value, data.Key)
	}

	item.ID = data.Key.ID
	item.Description = data.Description
	id, err := crudOp.recordsOp.Save(item, identity)
	if err != nil {
		return nil, errors.Wrap(err, onSave)
	}

	return &crud.Key{Type: CRUD01, ID: id}, nil
}

const onRead = "on records01/crud.Read()"

func (crudOp *records01CRUD) Read(key crud.Key, identity *auth.Identity) (*crud.Data, error) {
	if key.Type != CRUD01 {
		return nil, fmt.Errorf(onRead+": wrong key.Type (%#v)", key)
	}

	item, err := crudOp.recordsOp.Read(key.ID, identity)
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

func (crudOp *records01CRUD) List(crudType crud.Type, _ selectors.Options, identity *auth.Identity) ([]crud.Data, error) {
	if crudType != CRUD01 {
		return nil, fmt.Errorf(onList+": wrong crudType (%#v)", crudType)
	}

	// TODO!!! use selector
	items, err := crudOp.recordsOp.List(nil, identity)
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

func (crudOp *records01CRUD) Remove(key crud.Key, identity *auth.Identity) error {
	if key.Type != CRUD01 {
		return fmt.Errorf(onRemove+": wrong key.Type (%#v)", key)
	}

	return crudOp.recordsOp.Remove(key.ID, identity)
}
