package crud

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/rbac"

	"github.com/pavlo67/data/elements/selectors"
)

//const numRepeats = 3
//const toReadI = 0   // must be < numRepeats
//const toUpdateI = 1 // must be < numRepeats

func OperatorTestScenario(t *testing.T, crudOp Operator, crudCleanerOp db.Cleaner, crudType Type, itemToSave interface{}, changeItem ChangeItem) {
	require.NotNil(t, crudOp)
	require.NotNil(t, crudCleanerOp)

	var err error

	// prepare... ----------------------------------------------

	require.Equal(t, "test", os.Getenv("ENV"))

	adminIdentity := auth.IdentityWithRoles(rbac.RoleAdmin)
	require.NotNil(t, adminIdentity)

	require.NotNil(t, itemToSave)
	require.NotNil(t, changeItem)

	// clean old data ------------------------------------------

	err = crudCleanerOp.Clean()
	require.NoError(t, err)

	CountTestItems(t, crudOp, crudType, adminIdentity, 0)

	// add item ----------------------------------------------

	savedKey, err := crudOp.Save(Key{Type: crudType}, itemToSave, adminIdentity)
	require.NoError(t, err)
	require.NotNil(t, savedKey)
	require.Equal(t, crudType, savedKey.Type)
	require.NotEmpty(t, savedKey.ID)

	CountTestItems(t, crudOp, crudType, adminIdentity, 1)

	// read item ---------------------------------------------

	crudReaded, err := crudOp.Read(*savedKey, adminIdentity)
	require.NoError(t, err)
	require.NotNil(t, crudReaded)

	err = crudOp.TestIfEqual(t, *savedKey, itemToSave, crudReaded)
	require.NoError(t, err)

	// change item -------------------------------------------

	itemChanged, err := changeItem(crudReaded, *savedKey)

	savedChangedKey, err := crudOp.Save(*savedKey, itemChanged, adminIdentity)
	require.NoError(t, err)
	require.Equal(t, *savedKey, *savedChangedKey)

	CountTestItems(t, crudOp, crudType, adminIdentity, 1)

	// read crud ---------------------------------------------

	crudReaded, err = crudOp.Read(*savedKey, adminIdentity)
	require.NoError(t, err)
	require.NotNil(t, crudReaded)

	err = crudOp.TestIfEqual(t, *savedKey, itemChanged, crudReaded)
	require.NoError(t, err)

	// remove crud -------------------------------------------

	err = crudOp.Remove(*savedKey, adminIdentity)
	require.NoError(t, err)
	crudReaded, err = crudOp.Read(*savedKey, adminIdentity)
	require.Error(t, err)
	require.Nil(t, crudReaded)
	CountTestItems(t, crudOp, crudType, adminIdentity, 0)

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

}

func CountTestItems(t *testing.T, crudOp Operator, crudType Type, identity *auth.Identity, expectedCount int) {
	crudItems, err := crudOp.List(crudType, selectors.Options{}, identity)
	require.NoError(t, err)
	require.Equalf(t, expectedCount, len(crudItems), "crudItems = %#v", crudItems)
}

//func testData(t *testing.T, op Operator, expectedData, data interface{}, expectedID Key) {
//	if expectedData == nil {
//		require.Nil(t, data)
//		return
//	}
//	require.NotNil(t, data)
//
//	err := op.TestIfEqual(expectedData, expectedID, data)
//	require.NoError(t, err)
//}
