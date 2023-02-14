package persons

import (
	"encoding/json"
	"fmt"

	"github.com/pavlo67/data/entities"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"

	"github.com/pavlo67/data/entities/contacts"
)

var TestItem = Item{
	Person: Person{
		Firstnames: []string{"Erich", "Maria"},
		Middlename: "???",
		Lastname:   "Remark",
		Nicknames:  []string{"erich1", "maria2"},
		Contacts:   []contacts.Item{{Type: "phone", Value: "777", Connected: []contacts.Item{{Type: "fax", Value: "888"}}}},
		Info:       common.Map{"info1": "data1", "info2": "data2"},
	},
	Description: entities.TestDescription,
}

var _ entities.ReadValueRaw = ReadValueRaw

func ReadValueRaw(message json.RawMessage) (interface{}, error) {
	var person Person

	if err := json.Unmarshal(message, &person); err != nil {
		return nil, fmt.Errorf(onSave+": can't unmarshal (%s) into item.Person", message)
	}

	return person, nil
}

var _ entities.ChangeItemForTest = ChangeCRUDItemForTest

const onChangeItem = "on item.ChangeCRUDItemForTest()"

func ChangeCRUDItemForTest(data entities.Data, key entities.Key) (*entities.Data, error) {
	var item Item

	switch v := data.Value.(type) {
	case Person:
		item = Item{Person: v}
	case *Person:
		if v == nil {
			return nil, errors.New(onChangeItem + ": nil Person to change")
		}
		item = Item{Person: *v}
	case json.RawMessage:
		if err := json.Unmarshal(v, &item.Person); err != nil {
			return nil, fmt.Errorf(onSave+": can't unmarshal (%s) into item.Person", v)
		}
	default:
		return nil, fmt.Errorf(onChangeItem+": wrong data (%#v) to change with key (%#v)", data, key)
	}

	changedItem := ChangeItemForTest(item, key.ID)

	return &entities.Data{
		Key: entities.Key{
			Type: CRUD,
			ID:   changedItem.ID,
		},
		Description: changedItem.Description,
		Value:       changedItem.Person,
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
