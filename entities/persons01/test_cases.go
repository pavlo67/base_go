package persons01

import (
	"fmt"

	"github.com/pavlo67/data/elements/crud"
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"

	"github.com/pavlo67/data/elements/contacts"
	"github.com/pavlo67/data/types"
)

var TestPersonToSave = types.Person01{
	Firstnames:  []string{"Erich", "Maria"},
	Middlename:  "???",
	Lastname:    "Remark",
	Nicknames:   []string{"erich1", "maria2"},
	Contacts:    []contacts.Item{{Type: "phone", Value: "777", Connected: []contacts.Item{{Type: "fax", Value: "888"}}}},
	Info:        common.Map{"info1": "data1", "info2": "data2"},
	Description: types.TestDescription01,
}

var _ crud.ChangeItem = ChangeTestCRUDItem

const onChangeItem = "on person01.ChangeTestCRUDItem()"

func ChangeTestCRUDItem(data interface{}, key crud.Key) (interface{}, error) {
	var item Item

	switch v := data.(type) {
	case Item:
		item = v
	case *Item:
		if v == nil {
			return nil, errors.New(onChangeItem + ": nil Item to change")
		}
		item = *v
	case types.Person01:
		item = Item{Person01: v}
	case *types.Person01:
		if v == nil {
			return nil, errors.New(onChangeItem + ": nil Person01 to change")
		}
		item = Item{Person01: *v}
	default:
		return nil, fmt.Errorf(onChangeItem+": wrong data (%#v) to change with key (%#v)", data, key)
	}

	return ChangeTestItem(item, key.ID), nil
}

func ChangeTestItem(personReaded Item, savedID ID) Item {
	personToSaveChanged := personReaded
	personToSaveChanged.ID = savedID
	personToSaveChanged.Firstnames = personToSaveChanged.Firstnames[:1]
	personToSaveChanged.Middlename += " (changed)"
	personToSaveChanged.Lastname += " (changed)"
	personToSaveChanged.Nicknames = append(personToSaveChanged.Nicknames, "changed")
	personToSaveChanged.Contacts = append(personToSaveChanged.Contacts, personToSaveChanged.Contacts...)

	if personToSaveChanged.Info == nil {
		personToSaveChanged.Info = common.Map{}
	}
	personToSaveChanged.Info["changed"] = "changed info"

	personToSaveChanged.Description = personToSaveChanged.Description.ChangeForTest()

	return personToSaveChanged
}
