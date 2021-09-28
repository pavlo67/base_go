package persons01

import (
	"fmt"
	"testing"

	"github.com/pavlo67/data/elements/crud"
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/stretchr/testify/require"

	"github.com/pavlo67/data/elements/contacts"
	"github.com/pavlo67/data/types"
)

var TestPersonToSave = types.Person01{
	Firstnames: []string{"Erich", "Maria"},
	Middlename: "???",
	Lastname:   "Remark",
	Nicknames:  []string{"erich1", "maria2"},
	Contacts:   []contacts.Item{{Type: "phone", Value: "777", Connected: []contacts.Item{{Type: "fax", Value: "888"}}}},
	Info:       common.Map{"info1": "data1", "info2": "data2"},
	Description: types.Description01{
		URN:  "urn1",
		Tags: []string{"famous", "writer"},
		RelationsMap: types.Relations01Map{"r": types.Relation01{
			Key:  "r1key",
			NSS:  "nss_r1",
			Note: "wetr wert eryry",
		}},
		OwnerNSS:  "owner_nss",
		ViewerNSS: "viever_nss",
		// History:      nil,
	},
}

var _ crud.ChangeItem = ChangeTestItem

const onChangeItem = "on person01.ChangeTestItem()"

func ChangeTestItem(data interface{}, key crud.Key) (interface{}, error) {
	var item Item

	switch v := data.(type) {
	case Item:
		item = v
		item.ID = key.ID
	case *Item:
		if v == nil {
			return nil, errors.New(onChangeItem + ": nil Item to change")
		}
		item = *v
		item.ID = key.ID
	case types.Person01:
		item = Item{
			ID:       key.ID,
			Person01: v,
		}
	case *types.Person01:
		if v == nil {
			return nil, errors.New(onChangeItem + ": nil Person01 to change")
		}
		item = Item{
			ID:       key.ID,
			Person01: *v,
		}
	default:
		return nil, fmt.Errorf(onChangeItem+": wrong data (%#v) to change with key (%#v)", data, key)
	}

	return item, nil
}

func TestPersonItemToSaveChanged(personReaded Item, savedID ID) Item {
	personToSaveChanged := personReaded
	personToSaveChanged.ID = savedID
	personToSaveChanged.Firstnames = personToSaveChanged.Firstnames[:1]
	personToSaveChanged.Middlename += " (changed)"
	personToSaveChanged.Lastname += " (changed)"
	personToSaveChanged.Nicknames = personToSaveChanged.Nicknames[:1]
	personToSaveChanged.Contacts = append(personToSaveChanged.Contacts, personToSaveChanged.Contacts...)

	if personToSaveChanged.Info == nil {
		personToSaveChanged.Info = common.Map{}
	}
	personToSaveChanged.Info["changed"] = "changed info"

	personToSaveChanged.Description.URN += "_changed"
	personToSaveChanged.Description.Tags = append(personToSaveChanged.Description.Tags, "changed_tag")
	if personToSaveChanged.Description.RelationsMap == nil {
		personToSaveChanged.Description.RelationsMap = types.Relations01Map{}
	}
	personToSaveChanged.Description.RelationsMap["changed"] = types.Relation01{
		Key:  "chg",
		NSS:  "qwer",
		Note: "wqer qwer",
	}
	personToSaveChanged.Description.OwnerNSS += "_changed"
	personToSaveChanged.Description.ViewerNSS += "_changed"

	return personToSaveChanged
}

func CheckTestPerson(t *testing.T, personExpected, personToCheck types.Person01) {
	descriptionExpected, descriptionToCheck := personExpected.Description, personToCheck.Description
	personExpected.Description, personToCheck.Description = types.Description01{}, types.Description01{}

	require.Equal(t, personExpected, personToCheck)
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
}
