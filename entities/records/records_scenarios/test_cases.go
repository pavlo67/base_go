package records_scenarios

import (
	"encoding/json"
	"fmt"

	"github.com/pavlo67/data/entities"

	"github.com/pavlo67/data/entities/records"

	"github.com/pkg/errors"
)

var TestItem = records.Item{
	Record: records.Record{
		Content: records.Content{
			Title:   "title1",
			Summary: "summary1",
			Type:    "something",
			Data:    "wqer ewr er/yhlk'; '",
		},
		Additions: []records.Content{{
			Title:   "et1",
			Summary: "es1",
			Type:    "anything",
			Data:    "wertesrytr eu yuik",
		}},
	},
	Description: entities.TestDescription,
}

var _ entities.ReadValueRaw = ReadValueRaw

func ReadValueRaw(message json.RawMessage) (interface{}, error) {
	var record records.Record

	if err := json.Unmarshal(message, &record); err != nil {
		return nil, fmt.Errorf(records.onSave+": can't unmarshal (%s) into item.Record", message)
	}

	return record, nil
}

var _ entities.ChangeItemForTest = ChangeCRUDItemForTest

const onChangeItem = "on records.ChangeCRUDItemForTest()"

func ChangeCRUDItemForTest(data entities.Data, key entities.Key) (*entities.Data, error) {
	var item records.Item

	switch v := data.Value.(type) {
	case records.Record:
		item = records.Item{Record: v}
	case *records.Record:
		if v == nil {
			return nil, errors.New(onChangeItem + ": nil Record01 to change")
		}
		item = records.Item{Record: *v}
	case json.RawMessage:
		if err := json.Unmarshal(v, &item.Record); err != nil {
			return nil, fmt.Errorf(records.onSave+": can't unmarshal (%s) into item.Record", v)
		}
	default:
		return nil, fmt.Errorf(onChangeItem+": wrong data (%#v) to change with key (%#v)", data, key)
	}

	changedItem := ChangeItemForTest(item, key.ID)

	return &entities.Data{
		Key: entities.Key{
			Type: records.CRUD,
			ID:   changedItem.ID,
		},
		Description: changedItem.Description,
		Value:       changedItem.Record,
	}, nil
}

func ChangeItemForTest(recordReaded records.Item, savedID records.ID) records.Item {
	recordToSaveChanged := recordReaded
	recordToSaveChanged.ID = savedID

	recordToSaveChanged.Title += " (changed)"
	recordToSaveChanged.Summary += " (changed)"
	recordToSaveChanged.Type += " (changed)"
	recordToSaveChanged.Data += " (changed)"

	recordToSaveChanged.Additions = append(recordToSaveChanged.Additions, recordToSaveChanged.Additions...)

	recordToSaveChanged.Description = recordToSaveChanged.Description.ChangeForTest()

	return recordToSaveChanged
}
