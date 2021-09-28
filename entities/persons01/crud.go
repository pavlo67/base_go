package persons01

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

func (crudOp *persons01CRUD) Save(key crud.Key, data interface{}, identity *auth.Identity) (*crud.Key, error) {
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
	case types.Person01:
		item = Item{
			ID:       key.ID,
			Person01: v,
		}
	case *types.Person01:
		if v == nil {
			return nil, errors.New(onSave + ": nil Person01 to save")
		}
		item = Item{
			ID:       key.ID,
			Person01: *v,
		}
	default:
		return nil, fmt.Errorf(onSave+": wrong data (%#v) to save with key (%#v)", data, key)
	}

	id, err := crudOp.personsOp.Save(item, identity)
	if err != nil {
		return nil, errors.Wrap(err, onSave)
	}

	return &crud.Key{Type: CRUD01, ID: id}, nil
}

const onRead = "on persons01/crud.Read()"

func (crudOp *persons01CRUD) Read(key crud.Key, identity *auth.Identity) (interface{}, error) {
	if key.Type != CRUD01 {
		return nil, fmt.Errorf(onRead+": wrong key.Type (%#v)", key)
	}

	personItem, err := crudOp.personsOp.Read(key.ID, identity)
	if err != nil || personItem == nil {
		return nil, fmt.Errorf(onRead+": got %#v / %s", personItem, err)
	}

	return personItem, nil
}

const onList = "on persons01/crud.List()"

func (crudOp *persons01CRUD) List(crudType crud.Type, _ selectors.Options, identity *auth.Identity) ([]interface{}, error) {
	if crudType != CRUD01 {
		return nil, fmt.Errorf(onList+": wrong crudType (%#v)", crudType)
	}

	// TODO!!! use selector
	personItems, err := crudOp.personsOp.List(nil, identity)
	if err != nil {
		return nil, errors.Wrap(err, onList)
	}

	crudItems := make([]interface{}, len(personItems))
	for i, pi := range personItems {
		crudItems[i] = pi
	}

	return crudItems, nil

}

const onRemove = "on persons01/crud.Remove()"

func (crudOp *persons01CRUD) Remove(key crud.Key, identity *auth.Identity) error {
	if key.Type != CRUD01 {
		return fmt.Errorf(onRemove+": wrong key.Type (%#v)", key)
	}

	return crudOp.personsOp.Remove(key.ID, identity)
}

const onCheckIfEqual = "on persons01/crud.TestIfEqual()"

func (crudOp *persons01CRUD) TestIfEqual(t *testing.T, expectedKey crud.Key, expected interface{}, toCheck interface{}) error {

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
	case types.Person01:
		itemExpected = Item{
			ID:       expectedKey.ID,
			Person01: v,
		}
	case *types.Person01:
		if v == nil {
			return errors.New(onCheckIfEqual + ": nil Person01 expected")
		}
		itemExpected = Item{
			ID:       expectedKey.ID,
			Person01: *v,
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
	case types.Person01:
		itemToCheck = Item{
			Person01: v,
		}
	case *types.Person01:
		if v == nil {
			return errors.New(onCheckIfEqual + ": nil Person01 to check")
		}
		itemToCheck = Item{
			Person01: *v,
		}
	default:
		return fmt.Errorf(onSave+": wrong toCheck data (%#v)", toCheck)
	}

	descriptionExpected, descriptionToCheck := itemExpected.Description, itemToCheck.Description
	itemExpected.Description, itemToCheck.Description = types.Description01{}, types.Description01{}

	require.Equal(t, itemExpected, itemToCheck)
	require.Equal(t, descriptionExpected.URN, descriptionToCheck.URN)

	if len(descriptionExpected.Tags) > 0 {
		require.Equal(t, descriptionExpected.Tags, descriptionToCheck.Tags)
	} else {
		require.Equal(t, 0, len(descriptionToCheck.Tags))
	}
	if len(descriptionExpected.RelationsMap) > 0 {
		require.Equal(t, descriptionExpected.RelationsMap, descriptionToCheck.RelationsMap)
	} else {
		require.Equal(t, 0, len(descriptionToCheck.RelationsMap))
	}

	require.Equal(t, descriptionExpected.ViewerNSS, descriptionToCheck.ViewerNSS)
	require.Equal(t, descriptionExpected.OwnerNSS, descriptionToCheck.OwnerNSS)

	require.True(t, len(descriptionToCheck.History) >= len(descriptionExpected.History))
	require.Equal(t, descriptionExpected.History, descriptionToCheck.History[:len(descriptionExpected.History)])

	return nil
}
