package persons01

import (
	"fmt"

	"github.com/pavlo67/data/components/crud"

	"github.com/pavlo67/data/entities"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"

	"github.com/pavlo67/data/elements/contacts"
)

var TestItem = Item{
	Person01: entities.Person01{
		Firstnames: []string{"Erich", "Maria"},
		Middlename: "???",
		Lastname:   "Remark",
		Nicknames:  []string{"erich1", "maria2"},
		Contacts:   []contacts.Item{{Type: "phone", Value: "777", Connected: []contacts.Item{{Type: "fax", Value: "888"}}}},
		Info:       common.Map{"info1": "data1", "info2": "data2"},
	},
	Description: crud.TestDescription01,
}

var _ crud.ChangeItem = ChangeTestCRUDItem

const onChangeItem = "on person01.ChangeTestCRUDItem()"

func ChangeTestCRUDItem(data crud.Data, key crud.Key) (*crud.Data, error) {
	var item Item

	switch v := data.Value.(type) {
	case entities.Person01:
		item = Item{Person01: v}
	case *entities.Person01:
		if v == nil {
			return nil, errors.New(onChangeItem + ": nil Person01 to change")
		}
		item = Item{Person01: *v}
	default:
		return nil, fmt.Errorf(onChangeItem+": wrong data (%#v) to change with key (%#v)", data, key)
	}

	changedItem := ChangeTestItem(item, key.ID)

	return &crud.Data{
		Key: crud.Key{
			Type: CRUD01,
			ID:   changedItem.ID,
		},
		Description: changedItem.Description,
		Value:       changedItem.Person01,
	}, nil
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
