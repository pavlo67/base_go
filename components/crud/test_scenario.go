package crud

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/pavlo67/data/components/selectors"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/db"
)

func OperatorTestScenario(t *testing.T, crudOp Operator, crudCleanerOp db.Cleaner,
	itemToSave Data, readValueRaw ReadValueRaw, changeItemForTest ChangeItemForTest, actor auth.Actor) {
	require.NotNil(t, crudOp)
	require.NotNil(t, crudCleanerOp)

	var err error

	// prepare... ----------------------------------------------

	require.Equal(t, "test", os.Getenv("ENV"))

	require.NotNil(t, itemToSave)
	require.NotNil(t, changeItemForTest)

	crudType := itemToSave.Key.Type
	require.NotEmpty(t, crudType)

	// clean old data ------------------------------------------

	err = crudCleanerOp.Clean()
	require.NoError(t, err)

	CountTestItems(t, crudOp, crudType, actor, 0)

	// insert item -------------------------------------------

	t.Log("add item")

	var savedKey *Key
	savedKey, itemToSave.History, err = crudOp.Save(itemToSave, actor)
	require.NoError(t, err)
	require.NotNil(t, savedKey)
	require.Equal(t, crudType, savedKey.Type)
	require.NotEmpty(t, savedKey.ID)

	CountTestItems(t, crudOp, crudType, actor, 1)

	// read item ---------------------------------------------

	t.Log("read item")

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

	// update item -------------------------------------------

	t.Log("change item")

	itemChanged, err := changeItemForTest(*crudReaded, *savedKey)
	require.NoError(t, err)
	require.NotNil(t, itemChanged)

	itemChanged.Description.URN = urnOriginal + "_changed"

	var savedChangedKey *Key
	savedChangedKey, itemChanged.History, err = crudOp.Save(*itemChanged, actor)
	require.NoError(t, err)
	require.Equal(t, *savedKey, *savedChangedKey)

	CountTestItems(t, crudOp, crudType, actor, 1)

	// read item ---------------------------------------------

	t.Log("read item")

	crudReaded, err = crudOp.Read(*savedKey, actor)
	require.NoError(t, err)
	require.NotNil(t, crudReaded)

	crudSavedUpdated := Data{
		Key:         *savedKey,
		Description: itemChanged.Description,
		Value:       itemChanged.Value,
	}

	crudSavedUpdated.Description.URN = urnOriginal // UNCHANGED!!!

	TestIfEqual(t, crudSavedUpdated, *crudReaded, readValueRaw)
	require.NoError(t, err)

	// remove item -------------------------------------------

	t.Log("remove item")

	err = crudOp.Remove(*savedKey, actor)
	require.NoError(t, err)
	crudReaded, err = crudOp.Read(*savedKey, actor)
	require.Error(t, err)
	require.Nil(t, crudReaded)
	CountTestItems(t, crudOp, crudType, actor, 0)

}

//// add another crud --------------------------------------
//
//passwordToSaveAnother := "passwordToSaveAnother"
//crudToSaveAnother := types.Person01{
//	Identity: auth.Identity{
//		Nickname: "test_nickname2",
//		Roles:    rbac.Roles{rbac.RoleUser},
//	},
//}
//err = crudToSaveAnother.SetCreds(auth.Creds{auth.CredsPassword: passwordToSaveAnother})
//require.NoError(t, err)
//
//crudToSaveAnother.ID, err = crudOp.Save(crudToSaveAnother, adminIdentity)
//require.NoErrorf(t, err, "%#v", err)
//require.NotEmpty(t, crudToSaveAnother.ID)
//
//CountTestItems(t, crudOp, adminIdentity, 2)

//// add crud ----------------------------------------------
//
//crudItems, err = crudOp.List(adminIdentity)
//require.NoErrorf(t, err, "%#v", err)
//require.Equal(t, 3, len(crudItems))
//
//// list crud by itself: error ---------------------------
//
//crudItems, err = crudOp.List(&crud1Options)
//require.Errorf(t, err, "%#v", err)
//require.Empty(t, crudItems)

//// change crud by admin: ok ------------------------------
//
//crud1ToChange := *crudReaded
//crud1ToChange.Nickname += "_changed"
//
//crud1Changed, err := crudOp.Change(crud1ToChange, adminIdentity)
//require.NoErrorf(t, err, "%#v", err)
//require.Equal(t, crud1ToChange.Identity, crud1Changed.Identity)
//
//crud1ChangedReaded, err := crudOp.Read(crud1Changed.ID, adminIdentity)
//require.NoErrorf(t, err, "%#v", err)
//require.Equal(t, crud1ToChange.Identity, crud1ChangedReaded.Identity)
//
//// change crud by itself: ok -----------------------------
//
//crud1ToChange.Nickname += "_again"
//
//crud1Changed, err = crudOp.Change(crud1ToChange, &crud1Options)
//require.NoErrorf(t, err, "%#v", err)
//require.Equal(t, crud1ToChange.Identity, crud1Changed.Identity)
//
//crud1ChangedReaded, err = crudOp.Read(crud1Changed.ID, &crud1Options)
//require.NoErrorf(t, err, "%#v", err)
//require.Equal(t, crud1ToChange.Identity, crud1ChangedReaded.Identity)
//
//// change/read crud by another crud: error -------------
//
//crud1ToChangeAgain := *crud1ChangedReaded
//crud1ToChangeAgain.Nickname += "_again2"
//
//crud1ChangedWrong, err := crudOp.Change(crud1ToChangeAgain, &crud2Options)
//require.Errorf(t, err, "%#v", err)
//require.Nil(t, crud1ChangedWrong)
//
//crud1ReadedWrong, err := crudOp.Read(crudID1, &crud2Options)
//require.Errorf(t, err, "%#v", err)
//require.Nil(t, crud1ReadedWrong)
//
//crud1Readed, err := crudOp.Read(crudID1, &crud1Options)
//require.NoErrorf(t, err, "%#v", err)
//require.NotNil(t, crud1Readed)
//require.Equal(t, crud1Changed.Identity, crud1Readed.Identity)
//
// remove crud by admin: ok ------------------------------

//err = crudOp.Remove(crudToSaveAnother.ID, adminIdentity)
//require.NoErrorf(t, err, "%#v", err)
//
//crudAnotherReaded, err := crudOp.Read(crudToSaveAnother.ID, adminIdentity)
//require.Errorf(t, err, "%#v", err)
//require.Nil(t, crudAnotherReaded)
//
//CountTestItems(t, crudOp, adminIdentity, 1)

//// remove crud by itself: ok -----------------------------
//
//require.NotNil(t, crud2Options.Identity)
//err = crudOp.Remove(crudID2, &crud2Options)
//require.NoErrorf(t, err, "%#v / %#v", crud2Options.Identity, err)
//
//crud2Readed, err := crudOp.Read(crudID2, &crud2Options)
//require.Errorf(t, err, "%#v", err)
//require.Nil(t, crud2Readed)
//
//// remove crud by another crud: error ------------------
//
//err = crudOp.Remove(crudID1, &crud2Options)
//require.Errorf(t, err, "%#v", err)
//
//crud1Readed, err = crudOp.Read(crudID1, adminIdentity)
//require.NoErrorf(t, err, "%#v", err)
//require.NotNil(t, crud1Readed)
//require.Equal(t, crud1ChangedReaded.Identity, crud1Readed.Identity)

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
	// t.Logf("111111111111111 %#v", actor)

	crudItems, err := crudOp.List(crudType, selectors.Options{}, actor)
	require.NoError(t, err)
	require.Equalf(t, expectedCount, len(crudItems), "crudItems = %#v", crudItems)
}
