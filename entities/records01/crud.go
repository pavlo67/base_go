package records01

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/auth"

	"github.com/pavlo67/data/elements/crud"
	"github.com/pavlo67/data/elements/selectors"

	"github.com/pavlo67/data/types"
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

func (crudOp *records01CRUD) Save(key crud.Key, data interface{}, identity *auth.Identity) (*crud.Key, error) {
	if key.Type != CRUD01 {
		return nil, fmt.Errorf(onSave+": wrong key.Type (%#v) to save item (%#v)", key, data)
	}

	var item Item

	switch v := data.(type) {
	case Item:
		item = v
		if key.ID != nil {
			item.ID = key.ID
		}
	case *Item:
		if v == nil {
			return nil, errors.New(onSave + ": nil Item to save")
		}
		item = *v
		if key.ID != nil {
			item.ID = key.ID
		}
	case types.Record01:
		item = Item{
			ID:       key.ID,
			Record01: v,
		}
	case *types.Record01:
		if v == nil {
			return nil, errors.New(onSave + ": nil Record01 to save")
		}
		item = Item{
			ID:       key.ID,
			Record01: *v,
		}
	default:
		return nil, fmt.Errorf(onSave+": wrong data (%#v) to save with key (%#v)", data, key)
	}

	id, err := crudOp.recordsOp.Save(item, identity)
	if err != nil {
		return nil, errors.Wrap(err, onSave)
	}

	return &crud.Key{Type: CRUD01, ID: id}, nil
}

const onRead = "on records01/crud.Read()"

func (crudOp *records01CRUD) Read(key crud.Key, identity *auth.Identity) (interface{}, error) {
	if key.Type != CRUD01 {
		return nil, fmt.Errorf(onRead+": wrong key.Type (%#v)", key)
	}

	personItem, err := crudOp.recordsOp.Read(key.ID, identity)
	if err != nil || personItem == nil {
		return nil, fmt.Errorf(onRead+": got %#v / %s", personItem, err)
	}

	return personItem, nil
}

const onList = "on records01/crud.List()"

func (crudOp *records01CRUD) List(crudType crud.Type, _ selectors.Options, identity *auth.Identity) ([]interface{}, error) {
	if crudType != CRUD01 {
		return nil, fmt.Errorf(onList+": wrong crudType (%#v)", crudType)
	}

	// TODO!!! use selector
	personItems, err := crudOp.recordsOp.List(nil, identity)
	if err != nil {
		return nil, errors.Wrap(err, onList)
	}

	crudItems := make([]interface{}, len(personItems))
	for i, pi := range personItems {
		crudItems[i] = pi
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

const onCheckIfEqual = "on records01/crud.TestIfEqual()"

func (crudOp *records01CRUD) TestIfEqual(t *testing.T, expectedKey crud.Key, expected interface{}, toCheck interface{}) error {

	var itemExpected, itemToCheck Item

	switch v := expected.(type) {
	case Item:
		itemExpected = v
		itemExpected.ID = expectedKey.ID
	case *Item:
		if v == nil {
			return errors.New(onCheckIfEqual + ": nil Item expected")
		}
		itemExpected = *v
		itemExpected.ID = expectedKey.ID
	case types.Record01:
		itemExpected = Item{
			ID:       expectedKey.ID,
			Record01: v,
		}
	case *types.Record01:
		if v == nil {
			return errors.New(onCheckIfEqual + ": nil Record01 expected")
		}
		itemExpected = Item{
			ID:       expectedKey.ID,
			Record01: *v,
		}
	default:
		return fmt.Errorf(onSave+": wrong expected data (%#v)", expected)
	}

	switch v := toCheck.(type) {
	case Item:
		itemToCheck = v
	case *Item:
		if v == nil {
			return errors.New(onCheckIfEqual + ": nil Item to check")
		}
		itemToCheck = *v
	case types.Record01:
		itemToCheck = Item{
			Record01: v,
		}
	case *types.Record01:
		if v == nil {
			return errors.New(onCheckIfEqual + ": nil Record01 to check")
		}
		itemToCheck = Item{
			Record01: *v,
		}
	default:
		return fmt.Errorf(onSave+": wrong toCheck data (%#v)", toCheck)
	}

	descriptionExpected, descriptionToCheck := itemExpected.Description, itemToCheck.Description
	itemExpected.Description, itemToCheck.Description = types.Description01{}, types.Description01{}

	require.Equal(t, itemExpected, itemToCheck)
	descriptionExpected.TestIfEqual(t, descriptionToCheck)

	return nil
}
