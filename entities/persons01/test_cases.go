package persons01

import (
	"encoding/json"
	"fmt"

	"github.com/pavlo67/data/entities/crud01"

	"github.com/pavlo67/data/components/contacts"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"

	"github.com/pavlo67/data/entities"

	"github.com/pavlo67/data/components/crud"
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
	Description: crud01.TestDescription,
}

var _ crud.ReadValueRaw = ReadValueRaw

func ReadValueRaw(message json.RawMessage) (interface{}, error) {
	var person01 entities.Person01

	if err := json.Unmarshal(message, &person01); err != nil {
		return nil, fmt.Errorf(onSave+": can't unmarshal (%s) into item.Person01", message)
	}

	return person01, nil
}

var _ crud.ChangeItemForTest = ChangeCRUDItemForTest

const onChangeItem = "on person01.ChangeCRUDItemForTest()"

func ChangeCRUDItemForTest(data crud.Data, key crud.Key) (*crud.Data, error) {
	var item Item

	switch v := data.Value.(type) {
	case entities.Person01:
		item = Item{Person01: v}
	case *entities.Person01:
		if v == nil {
			return nil, errors.New(onChangeItem + ": nil Person01 to change")
		}
		item = Item{Person01: *v}
	case json.RawMessage:
		if err := json.Unmarshal(v, &item.Person01); err != nil {
			return nil, fmt.Errorf(onSave+": can't unmarshal (%s) into item.Person01", v)
		}
	default:
		return nil, fmt.Errorf(onChangeItem+": wrong data (%#v) to change with key (%#v)", data, key)
	}

	changedItem := ChangeItemForTest(item, key.ID)

	return &crud.Data{
		Key: crud.Key{
			Type: CRUD01,
			ID:   changedItem.ID,
		},
		Description: changedItem.Description,
		Value:       changedItem.Person01,
	}, nil
}

func ChangeItemForTest(personReaded Item, savedID ID) Item {
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
