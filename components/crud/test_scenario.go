package crud

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/db"
)

func OperatorTestScenario(t *testing.T, crudOp Operator, crudCleanerOp db.Cleaner,
	itemToSave Data, readValueRaw ReadValueRaw, changeItemForTest ChangeItemForTest, actor auth.Actor) {
	require.NotNil(t, crudOp)
	require.NotNil(t, crudCleanerOp)

	var err error

	// preparation ---------------------------------------------

	require.Equal(t, "test", os.Getenv("ENV"))

	require.NotNil(t, itemToSave)
	require.NotNil(t, changeItemForTest)

	crudType := itemToSave.Key.Type
	require.NotEmpty(t, crudType)

	// old data clean-up ---------------------------------------

	err = crudCleanerOp.Clean()
	require.NoError(t, err)

	CountTestItems(t, crudOp, crudType, actor, 0)

	t.Log("database is cleaned")

	// inserting -----------------------------------------------

	var savedKey *Key
	savedKey, itemToSave.History, err = crudOp.Save(itemToSave, actor)
	require.NoError(t, err)
	require.NotNil(t, savedKey)
	require.Equal(t, crudType, savedKey.Type)
	require.NotEmpty(t, savedKey.ID)

	CountTestItems(t, crudOp, crudType, actor, 1)

	t.Log("record is inserted")

	// reading -------------------------------------------------

	crudReaded, err := crudOp.Read(*savedKey, actor)
	require.NoError(t, err)
	require.NotNil(t, crudReaded)

	crudSaved := Data{
		Key:         *savedKey,
		Description: itemToSave.Description,
		Value:       itemToSave.Value,
	}

	urnOriginal := itemToSave.Description.URN

	TestIfEqual(t, crudSaved, *crudReaded, readValueRaw)
	require.NoError(t, err)

	t.Log("record is read")

	// updating ------------------------------------------------

	itemChanged, err := changeItemForTest(*crudReaded, *savedKey)
	require.NoError(t, err)
	require.NotNil(t, itemChanged)

	itemChanged.Description.URN = urnOriginal + "_changed"

	savedChangedKey, historyChanged, err := crudOp.Save(*itemChanged, actor)
	require.NoError(t, err)
	require.Equal(t, *savedKey, *savedChangedKey)

	CountTestItems(t, crudOp, crudType, actor, 1)

	t.Log("record is updated")

	//// updating (with unchanged .History) failure --------------
	//
	//savedChangedKey, _, err = crudOp.Save(*itemChanged, actor)
	//require.Error(t, err)
	//require.Nil(t, savedChangedKey)
	//
	//CountTestItems(t, crudOp, crudType, actor, 1)

	// updating ------------------------------------------------

	itemChanged.History = historyChanged
	savedChangedKey, itemChanged.History, err = crudOp.Save(*itemChanged, actor)
	require.NoError(t, err)
	require.Equal(t, *savedKey, *savedChangedKey)

	CountTestItems(t, crudOp, crudType, actor, 1)

	t.Log("record is updated again")

	// reading -------------------------------------------------

	crudReaded, err = crudOp.Read(*savedKey, actor)
	require.NoError(t, err)
	require.NotNil(t, crudReaded)

	crudSavedUpdated := Data{
		Key:         *savedKey,
		Description: itemChanged.Description,
		Value:       itemChanged.Value,
	}

	// // TODO: be careful, item.Description.URN left UNCHANGED
	// crudSavedUpdated.Description.URN = urnOriginal
	// ???

	TestIfEqual(t, crudSavedUpdated, *crudReaded, readValueRaw)
	require.NoError(t, err)

	t.Log("record is read")

	// removing ------------------------------------------------

	err = crudOp.Remove(*savedKey, actor)
	require.NoError(t, err)
	crudReaded, err = crudOp.Read(*savedKey, actor)
	require.Error(t, err)
	require.Nil(t, crudReaded)
	CountTestItems(t, crudOp, crudType, actor, 0)

	t.Log("record is removad")

}

const onCheckIfEqual = "on persons01/crud.TestIfEqual()"

func TestIfEqual(t *testing.T, expected, toCheck Data, readValueRaw ReadValueRaw) error { //
	require.Equal(t, expected.Key, toCheck.Key)
	expected.Description.TestIfEqual(t, toCheck.Description)

	value := toCheck.Value
	if jsonRawMessage, ok := value.(json.RawMessage); ok {
		var err error
		value, err = readValueRaw(jsonRawMessage)
		require.NoError(t, err)
	}

	require.Equal(t, expected.Value, value)

	return nil
}

func CountTestItems(t *testing.T, crudOp Operator, crudType Type, actor auth.Actor, expectedCount int) {
	crudItems, err := crudOp.List(crudType, Options{}, actor)
	require.NoError(t, err)
	require.Equalf(t, expectedCount, len(crudItems), "crudItems = %#v", crudItems)
}
