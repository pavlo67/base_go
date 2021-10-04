package records01

import (
	"fmt"

	"github.com/pavlo67/data/components/crud"

	"github.com/pavlo67/data/entities"

	"github.com/pkg/errors"
)

var TestRecord = entities.Record01{
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
	Description: entities.TestDescription01,
}

var _ crud.ChangeItem = ChangeForTest

const onChangeItem = "on records01.ChangeForTest()"

func ChangeForTest(data interface{}, key crud.Key) (interface{}, error) {
	var item Item

	switch v := data.(type) {
	case Item:
		item = v
	case *Item:
		if v == nil {
			return nil, errors.New(onChangeItem + ": nil Item to change")
		}
		item = *v
	case entities.Record01:
		item = Item{Record01: v}
	case *entities.Record01:
		if v == nil {
			return nil, errors.New(onChangeItem + ": nil Record01 to change")
		}
		item = Item{Record01: *v}
	default:
		return nil, fmt.Errorf(onChangeItem+": wrong data (%#v) to change with key (%#v)", data, key)
	}

	return ChangeTestItem(item, key.ID), nil
}

func ChangeTestItem(recordReaded Item, savedID ID) Item {
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
