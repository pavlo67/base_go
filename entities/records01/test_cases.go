package records01

import (
	"encoding/json"
	"fmt"

	"github.com/pavlo67/data/components/crud"
	crud012 "github.com/pavlo67/data/components/crud01"

	"github.com/pavlo67/data/entities"

	"github.com/pkg/errors"
)

var TestItem = Item{
	Record01: entities.Record01{
		Content01: entities.Content01{
			Title:   "title1",
			Summary: "summary1",
			Type:    "something",
			Data:    "wqer ewr er/yhlk'; '",
		},
		Embedded: []entities.Content01{{
			Title:   "et1",
			Summary: "es1",
			Type:    "anything",
			Data:    "wertesrytr eu yuik",
		}},
	},
	Description: crud012.TestDescription01,
}

var _ crud.ReadValueRaw = ReadValueRaw

func ReadValueRaw(message json.RawMessage) (interface{}, error) {
	var record01 entities.Record01

	if err := json.Unmarshal(message, &record01); err != nil {
		return nil, fmt.Errorf(onSave+": can't unmarshal (%s) into item.Record01", message)
	}

	return record01, nil
}

var _ crud.ChangeItemForTest = ChangeCRUDItemForTest

const onChangeItem = "on records01.ChangeCRUDItemForTest()"

func ChangeCRUDItemForTest(data crud.Data, key crud.Key) (*crud.Data, error) {
	var item Item

	switch v := data.Value.(type) {
	case entities.Record01:
		item = Item{Record01: v}
	case *entities.Record01:
		if v == nil {
			return nil, errors.New(onChangeItem + ": nil Record01 to change")
		}
		item = Item{Record01: *v}
	case json.RawMessage:
		if err := json.Unmarshal(v, &item.Record01); err != nil {
			return nil, fmt.Errorf(onSave+": can't unmarshal (%s) into item.Record01", v)
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
		Value:       changedItem.Record01,
	}, nil
}

func ChangeItemForTest(recordReaded Item, savedID ID) Item {
	recordToSaveChanged := recordReaded
	recordToSaveChanged.ID = savedID

	recordToSaveChanged.Title += " (changed)"
	recordToSaveChanged.Summary += " (changed)"
	recordToSaveChanged.Type += " (changed)"
	recordToSaveChanged.Data += " (changed)"

	recordToSaveChanged.Embedded = append(recordToSaveChanged.Embedded, recordToSaveChanged.Embedded...)

	recordToSaveChanged.Description = recordToSaveChanged.Description.ChangeForTest()

	return recordToSaveChanged
}
