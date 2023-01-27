package records

import (
	"encoding/json"
	"fmt"

	"github.com/pkg/errors"

	"github.com/pavlo67/data/components/crud"
)

var TestItem = Item{
	Record: Record{
		Content: Content{
			Title:   "title1",
			Summary: "summary1",
			Type:    "something",
			Data:    "wqer ewr er/yhlk'; '",
		},
		Embedded: []Content{{
			Title:   "et1",
			Summary: "es1",
			Type:    "anything",
			Data:    "wertesrytr eu yuik",
		}},
	},
	Description: crud.TestDescription,
}

var _ crud.ReadValueRaw = ReadValueRaw

func ReadValueRaw(message json.RawMessage) (interface{}, error) {
	var record Record

	if err := json.Unmarshal(message, &record); err != nil {
		return nil, fmt.Errorf(onSave+": can't unmarshal (%s) into item.Record", message)
	}

	return record, nil
}

var _ crud.ChangeItemForTest = ChangeCRUDItemForTest

const onChangeItem = "on records.ChangeCRUDItemForTest()"

func ChangeCRUDItemForTest(data crud.Data, key crud.Key) (*crud.Data, error) {
	var item Item

	switch v := data.Value.(type) {
	case Record:
		item = Item{Record: v}
	case *Record:
		if v == nil {
			return nil, errors.New(onChangeItem + ": nil Record01 to change")
		}
		item = Item{Record: *v}
	case json.RawMessage:
		if err := json.Unmarshal(v, &item.Record); err != nil {
			return nil, fmt.Errorf(onSave+": can't unmarshal (%s) into item.Record", v)
		}
	default:
		return nil, fmt.Errorf(onChangeItem+": wrong data (%#v) to change with key (%#v)", data, key)
	}

	changedItem := ChangeItemForTest(item, key.ID)

	return &crud.Data{
		Key: crud.Key{
			Type: CRUD,
			ID:   changedItem.ID,
		},
		Description: changedItem.Description,
		Value:       changedItem.Record,
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
