package records_scenarios

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/auth"
)

func readOkTest(t *testing.T, recordsOp Operator, item Item, identity auth.Identity) *Item {
	itemReaded, err := recordsOp.Read(item.ID, &identity)
	require.NoError(t, err)
	require.NotNil(t, itemReaded)

	require.Equal(t, item.ID, itemReaded.ID)
	require.Equal(t, item.Content, itemReaded.Content)
	//require.Equal(t, item.OwnerNSS, itemReaded.OwnerNSS)
	//require.Equal(t, item.ViewerNSS, itemReaded.ViewerNSS)

	items, err := recordsOp.List(nil, &identity)
	require.NoError(t, err)

	var itemFound *Item
	for _, itemListed := range items {
		if itemListed.ID == item.ID {
			itemFound = &itemListed
			require.Equal(t, item.ID, itemListed.ID)
			require.Equal(t, item.Content, itemListed.Content)
			//require.Equal(t, item.OwnerNSS, itemListed.OwnerNSS)
			//require.Equal(t, item.ViewerNSS, itemListed.ViewerNSS)
		}
	}
	require.NotNilf(t, itemFound, "%#v", items)

	return itemReaded
}

func readFailTest(t *testing.T, recordsOp Operator, itemID ID, identity auth.Identity) {
	itemReaded, err := recordsOp.Read(itemID, &identity)
	require.Error(t, err)
	require.Nil(t, itemReaded)

	items, err := recordsOp.List(nil, &identity)
	require.NoError(t, err)

	for _, itemListed := range items {
		if itemListed.ID == itemID {
			require.FailNow(t, "the item shouldn't be in list ", "%s -> %#v", itemID, itemListed)
		}
	}
}
