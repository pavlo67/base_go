package records

import (
	"encoding/json"
	"fmt"

	"github.com/pavlo67/data/entities"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/rbac"

	"github.com/pavlo67/data/components/vcs"
)

const CRUD entities.Type = "records"

var _ entities.OperatorCRUD = &recordsCRUD{}

func OperatorCRUD(recordsOp Operator, roles rbac.Roles) (entities.OperatorCRUD, error) {
	if recordsOp == nil {
		return nil, errors.New("recordsOp == nil")
	}

	return &recordsCRUD{recordsOp: recordsOp, roles: roles}, nil
}

type recordsCRUD struct {
	recordsOp Operator
	roles     rbac.Roles
}

func (crudOp *recordsCRUD) Types() ([]entities.Type, error) {
	return []entities.Type{CRUD}, nil
}

func (crudOp *recordsCRUD) Roles() (rbac.Roles, error) {
	return crudOp.roles, nil
}

const onSave = "on records/crud.Save()"

func (crudOp *recordsCRUD) Save(data entities.Data, actor auth.Actor) (*entities.Key, vcs.History, error) {
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
		if err := json.Unmarshal(v, &item.Record); err != nil {
			return nil, nil, fmt.Errorf(onSave+": can't unmarshal (%s) into item.Record", v)
		}
	case Record:
		item = Item{Record: v}
	case *Record:
		if v == nil {
			return nil, nil, errors.New(onSave + ": nil Record01 to save")
		}
		item = Item{Record: *v}
	default:
		return nil, nil, fmt.Errorf(onSave+": wrong data (%#v) to save with key (%#v)", data.Value, data.Key)
	}

	item.ID = data.Key.ID
	item.Description = data.Description
	id, _, historyChanged, err := crudOp.recordsOp.Save(item, actor)
	if err != nil {
		return nil, nil, errors.Wrap(err, onSave)
	}

	return &entities.Key{Type: CRUD, ID: id}, historyChanged, nil
}

const onRead = "on records/crud.Read()"

func (crudOp *recordsCRUD) Read(key entities.Key, actor auth.Actor) (*entities.Data, error) {
	if key.Type != CRUD {
		return nil, fmt.Errorf(onRead+": wrong key.Type (%#v)", key)
	}

	item, err := crudOp.recordsOp.Read(key.ID, actor)
	if err != nil || item == nil {
		return nil, fmt.Errorf(onRead+": got %#v / %s", item, err)
	}

	return &entities.Data{
		Key: entities.Key{
			Type: CRUD,
			ID:   key.ID,
		},
		Description: item.Description,
		Value:       item.Record,
	}, nil
}

const onList = "on records/crud.List()"

func (crudOp *recordsCRUD) List(crudType entities.Type, _ entities.Options, actor auth.Actor) ([]entities.Data, error) {
	if crudType != CRUD {
		return nil, fmt.Errorf(onList+": wrong crudType (%#v)", crudType)
	}

	// TODO!!! use selector
	items, err := crudOp.recordsOp.List(nil, actor)
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
			Value:       pi.Record,
		}
	}

	return crudItems, nil

}

const onRemove = "on records/crud.Remove()"

func (crudOp *recordsCRUD) Remove(key entities.Key, actor auth.Actor) error {
	if key.Type != CRUD {
		return fmt.Errorf(onRemove+": wrong key.Type (%#v)", key)
	}

	return crudOp.recordsOp.Remove(key.ID, actor)
}
