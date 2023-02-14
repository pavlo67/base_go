package records_scenarios

import (
	"os"
	"testing"

	"github.com/pavlo67/data/components/item"
	"github.com/pavlo67/data/components/ns"
	"github.com/pavlo67/data/components/tags"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/data/components/selectors"
)

// TODO: test .History
// TODO: test .List() with selectors

const authID = auth.ID("1")

const testNewTag1 = "testNewTag1"

var embedded1 = []Content{
	{
		Title:   "56567",
		Summary: "3333333",
		Type:    "test...",
		Data:    "wqerwer",
	},
}

var item1Original = Item{
	Content: Content{
		Title:   "345456",
		Summary: "6578gj",
		Type:    "test",
		Data:    `{"AAA": "aaa", "BBB": 222}`,
	},
	Embedded: embedded1,
	Description: item.Description{
		Tags: []tags.Item{"1", "333"},
	},
}

var item2Original = Item{
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

var item3Original = Item{
	Content: Content{
		Title:   "34545ee6rt",
		Summary: "6578weqreegj",
		Type:    "test2",
		Data:    `wqerwer`,
	},
	Embedded: append(embedded1, embedded1...),
	Description: item.Description{
		Tags: []tags.Item{"qw1", "333"},
	},
}

func OperatorTestScenarioNoRBAC(t *testing.T, joinerOp joiner.Operator) {

	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	recordsOp, _ := joinerOp.Interface(InterfaceKey).(Operator)
	require.NotNil(t, recordsOp)

	cleanerOp, _ := joinerOp.Interface(InterfaceCleanerKey).(db.Cleaner)
	require.NotNil(t, cleanerOp)

	// clear ------------------------------------------------------------------

	err := cleanerOp.Clean()
	require.NoError(t, err)

	// prepare identity to process tests --------------------------------------

	identity := auth.Identity{ID: authID}

	// save record without identity: error ------------------------------------

	item1Saved, err := recordsOp.Save(item1Original, nil)
	require.Error(t, err)
	require.Empty(t, item1Saved)

	// save record without OwnerNSS: added automatically, ok ----------------

	item1WithoutNSS := item1Original
	item1WithoutNSS.OwnerNSS = ""
	item1WithoutNSS.ID, err = recordsOp.Save(item1WithoutNSS, &identity)
	require.NoError(t, err)
	require.NotEmpty(t, item1WithoutNSS.ID)

	item1Readed := readOkTest(t, recordsOp, item1WithoutNSS, identity)
	require.NotNil(t, item1Readed)
	require.Equal(t, item1Readed.OwnerNSS, ns.NSS(authID))

	// save record, ok ------------------------------------------------------

	item2Original.ID, err = recordsOp.Save(item2Original, &identity)
	require.NoError(t, err)
	require.NotEmpty(t, item2Original.ID)

	item2Readed := readOkTest(t, recordsOp, item2Original, identity)
	require.NotNil(t, item2Readed)

	item1ReadedAgain := readOkTest(t, recordsOp, item1WithoutNSS, identity)
	require.NotNil(t, item1ReadedAgain)

	// ----------------------------------------------------------------------

	dbTestNoRBAC(t, recordsOp, item1Original, item2Original, item3Original, identity)

	item2Readed = readOkTest(t, recordsOp, item2Original, identity)
	require.NotNil(t, item2Readed)

	item1ReadedAgain = readOkTest(t, recordsOp, item1WithoutNSS, identity)
	require.NotNil(t, item1ReadedAgain)

	//
	//// check .Remove(), .Read(), .List(), -------------------------------------
	//
	//err = recordsOp.Remove(itemOriginal3Saved.ID, &identity)
	//require.NoError(t, err)
	//
	//readFailTest(t, recordsOp, itemOriginal3Saved.ID, identity)
	//readOkTest(t, recordsOp, item1Saved, identity)

	// save record without OwnerNSS: added automatically, ok ----------------
	// save record without OwnerNSS: added automatically, ok ----------------

}

func dbTestNoRBAC(t *testing.T, recordsOp Operator, itemToSave, itemToUpdate, itemToUpdateAgain Item, identity auth.Identity) {

	// prepare selector tagged ----------------------------

	selectorTagged := selectors.Term{
		Key:    HasTag,
		Values: []string{testNewTag1},
	}

	// insert ---------------------------------------------

	itemToSave.OwnerNSS = ns.NSS(authID)

	itemSaved1ID, err := recordsOp.Save(itemToSave, &identity)
	require.NoError(t, err)
	require.NotEmpty(t, itemSaved1ID)

	// TODO!!!
	// require.Equal(t, itemToSave.Content, itemSaved1.Content)

	itemToSave.ID = itemSaved1ID

	// prepare selector parent ----------------------------

	selectorParent := selectors.Term{
		Key:    HasParent,
		Values: []string{string(itemSaved1ID)},
	}

	// check inserted -------------------------------------

	readOkTest(t, recordsOp, itemToSave, identity)

	items, err := recordsOp.List(&selectorTagged, &identity)
	require.NoError(t, err)
	require.Equal(t, 0, len(items))

	items, err = recordsOp.List(&selectorParent, &identity)
	require.NoError(t, err)
	require.Equal(t, 0, len(items))

	// update ---------------------------------------------

	itemToUpdate.ID = itemToSave.ID
	itemToUpdate.Tags, err = recordsOp.AddParent(append(itemToUpdate.Tags, testNewTag1), itemToSave.ID)
	require.NoError(t, err)

	itemSaved2ID, err := recordsOp.Save(itemToUpdate, &identity)
	require.NoError(t, err)
	require.Equal(t, itemToUpdate.ID, itemSaved2ID)

	// TODO!!!
	// require.Equal(t, itemToUpdate.Content, itemSaved2.Content)

	readOkTest(t, recordsOp, itemToUpdate, identity)

	items, err = recordsOp.List(&selectorTagged, &identity)
	require.NoError(t, err)
	require.Equal(t, 1, len(items))
	require.Equal(t, items[0].ID, itemToSave.ID)

	items, err = recordsOp.List(&selectorParent, &identity)
	require.NoError(t, err)
	require.Equal(t, 1, len(items))
	require.Equal(t, items[0].ID, itemToSave.ID)

	// prepare item to update again ------------------------------------------

	itemToUpdateAgain.ID = itemToSave.ID

	itemSaved3ID, err := recordsOp.Save(itemToUpdateAgain, &identity)
	require.NoError(t, err)
	require.Equal(t, itemToUpdateAgain.ID, itemSaved3ID)

	// TODO!!!
	//require.Equal(t, itemToUpdateAgain.Content, itemSaved3.Content)

	readOkTest(t, recordsOp, itemToUpdateAgain, identity)

	//return itemToSave
}
