package records_scenarios

import (
	"os"
	"testing"

	"github.com/pavlo67/data/entities/records"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/data/components/item"
	"github.com/pavlo67/data/components/tags"
	"github.com/stretchr/testify/require"
)

// TODO: test .History
// TODO: test .List() with selectors

const authID1 = auth.ID("1")
const authID2 = auth.ID("2")
const authID3 = auth.ID("3")

const testNewTag = "testNewTag"

var embedded = []records.Content{
	{
		Title:   "56567",
		Summary: "3333333",
		Type:    "test...",
		Data:    "wqerwer",
	},
}

var item11 = Item{
	Content: Content{
		Title:   "345456",
		Summary: "6578gj",
		Type:    "test",
		Data:    `{"AAA": "aaa", "BBB": 222}`,
	},
	Embedded: embedded,
	Description: item.Description{
		Tags: []tags.Item{"1", "333"},
	},
}

var item12 = Item{
	Content: Content{
		Title:   "345eeeee456rt",
		Summary: "6578eegj",
		Type:    "test1",
		Data:    `{"AAA": "awraa", "BBB": 22552}`,
	},
	Description: item.Description{
		Tags: []tags.Item{"1", "333"},
	},
}

var item22 = Item{
	Content: Content{
		Title:   "34545ee6rt",
		Summary: "6578weqreegj",
		Type:    "test2",
		Data:    `wqerwer`,
	},
	Embedded: append(embedded, embedded...),
	Description: item.Description{
		Tags: []tags.Item{"qw1", "333"},
	},
}

func OperatorTestScenario(t *testing.T, recordsOp Operator, cleanerOp db.Cleaner, l logger.Operator) {

	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	// clear ------------------------------------------------------------------------------

	err := cleanerOp.Clean()
	require.NoError(t, err)

	// prepare records to test  -----------------------------------------------------------

	// prepare records & auth.Identity -----------------------------------------

	item01 := item11
	item01.OwnerNSS = ""

	identity0 := auth.Identity{}
	identity1 := auth.Identity{ID: authID1}

	// save record without identity1: error ------------------------------------

	item01SavedID, err := recordsOp.Save(item01, &identity0)
	require.Error(t, err)
	require.Empty(t, item01SavedID)

	// save record without OwnerNSS: added automatically, ok -------------------

	require.Empty(t, item01.OwnerNSS)
	item01SavedID, err = recordsOp.Save(item01, &identity1)
	require.NoError(t, err)

	// TODO!!!
	// require.Equal(t, item01SavedID.OwnerNSS, identity1.Identity.ID)

	saveTest(t, recordsOp, item11, item12, item22) // item22Saved :=
	saveTest(t, recordsOp, item11, item12, item22) // item22SavedAgain :=

	// check .Remove(), .Read(), .List(), -------------------------------------

	//owner22Identity := auth.Identity{ID: item22Saved.OwnerNSS}
	//owner22ViewerIdentity := auth.Identity{ID: item22Saved.ViewerNSS}
	//
	//err = recordsOp.Remove(item22Saved.ID, &owner22Identity)
	//require.NoError(t, err)
	//
	//readFailTest(t, recordsOp, item22Saved.ID, owner22Identity)
	//readFailTest(t, recordsOp, item22Saved.ID, owner22ViewerIdentity)
	//readOkTest(t, recordsOp, item22SavedAgain, owner22Identity)
	//readOkTest(t, recordsOp, item22SavedAgain, owner22ViewerIdentity)

}

func saveTest(t *testing.T, recordsOp Operator, itemToSave, itemToUpdate, itemToUpdateAgain Item) Item {
	//
	//identity1 := auth.Identity{ID: authID1}
	//identity2 := auth.Identity{ID: authID2}
	//identity3 := auth.Identity{ID: authID3}
	//
	//// prepare item to save --------------------------------------------------
	//
	//itemToSave.OwnerNSS = authID1
	//itemToSave.ViewerNSS = authID1
	//
	//// check .Add() with other identity: error -------------------------------
	//
	//itemSaved, err := recordsOp.Add(itemToSave, &identity2)
	//require.Error(t, err)
	//require.Empty(t, itemSaved)
	//
	//// check .Add() with owner identity: ok ----------------------------------
	//
	//itemSaved, err = recordsOp.Add(itemToSave, &identity1)
	//require.NoError(t, err)
	//require.NotEmpty(t, itemSaved)
	//require.Equal(t, itemToSave.Content, itemSaved.Content)
	//require.Equal(t, itemToSave.OwnerNSS, itemSaved.OwnerNSS)
	//require.Equal(t, itemToSave.ViewerNSS, itemSaved.ViewerNSS)
	//
	//itemToSave.ID = itemSaved.ID
	//
	//// check .Read(), .List() with owner/viewer identity ----------------------
	//
	//readOkTest(t, recordsOp, itemToSave, identity1)
	//
	//// check .Read(), .List() with other identity -----------------------------
	//
	//readFailTest(t, recordsOp, itemToSave.ID, identity2)
	//
	//// prepare item to update ------------------------------------------------
	//
	//itemToUpdate.ID = itemToSave.ID
	//itemToUpdate.OwnerNSS = authID1
	//itemToUpdate.ViewerNSS = authID2
	//
	//// check updating .Add() with other identity: error ----------------------
	//
	//itemSaved, err = recordsOp.Add(itemToUpdate, &identity2)
	//require.Error(t, err)
	//require.Empty(t, itemSaved)
	//
	//// check updating .Add() with owner identity: ok -------------------------
	//
	//itemSaved, err = recordsOp.Add(itemToUpdate, &identity1)
	//require.NoError(t, err)
	//require.NotEmpty(t, itemSaved)
	//require.Equal(t, itemToUpdate.ID, itemSaved.ID)
	//require.Equal(t, itemToUpdate.Content, itemSaved.Content)
	//require.Equal(t, itemToUpdate.OwnerNSS, itemSaved.OwnerNSS)
	//require.Equal(t, itemToUpdate.ViewerNSS, itemSaved.ViewerNSS)
	//
	//// check .Read(), .List() with owner identity -----------------------------
	//
	//readOkTest(t, recordsOp, itemToUpdate, identity1)
	//
	//// check .Read(), .List() with viewer identity ----------------------------
	//
	//readOkTest(t, recordsOp, itemToUpdate, identity2)
	//
	//// check .Read(), .List() with other identity -----------------------------
	//
	//readFailTest(t, recordsOp, itemToUpdate.ID, identity3)
	//
	//// prepare item to update again ------------------------------------------
	//
	//itemToUpdateAgain.ID = itemToSave.ID
	//itemToUpdateAgain.OwnerNSS = authID2
	//itemToUpdateAgain.ViewerNSS = authID2
	//
	//// check updating .Add() with other identity: error ----------------------
	//
	//itemSaved, err = recordsOp.Add(itemToUpdateAgain, &identity2)
	//require.Error(t, err)
	//require.Empty(t, itemSaved)
	//
	//// check updating .Add() with owner identity: ok -------------------------
	//
	//itemSaved, err = recordsOp.Add(itemToUpdateAgain, &identity1)
	//require.NoError(t, err)
	//require.NotEmpty(t, itemSaved)
	//require.Equal(t, itemToUpdateAgain.ID, itemSaved.ID)
	//require.Equal(t, itemToUpdateAgain.Content, itemSaved.Content)
	//require.Equal(t, itemToUpdateAgain.OwnerNSS, itemSaved.OwnerNSS)
	//require.Equal(t, itemToUpdateAgain.ViewerNSS, itemSaved.ViewerNSS)
	//
	//// check .Read(), .List() with owner/viewer identity ----------------------
	//
	//readOkTest(t, recordsOp, itemToUpdateAgain, identity2)
	//
	//// check .Read(), .List() with other identity -----------------------------
	//
	//readFailTest(t, recordsOp, itemToUpdateAgain.ID, identity1)
	//

	return itemToSave
}
