package crud

import (
	"fmt"
	"os"
	"testing"

	"github.com/pavlo67/data/elements/selectors"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/db"
)

type OperatorTestCase struct {
	Operator
	db.Cleaner

	ToCreate          interface{}
	ExpectedCreateErr bool
	ExpectedReadErr   bool
	ExpectedListErr   bool

	ToUpdate          interface{}
	ExpectedUpdateErr bool
	ExpectedDeleteErr bool
}

const numRepeats = 3
const toReadI = 0   // must be < numRepeats
const toUpdateI = 1 // must be < numRepeats

func OperatorTest(t *testing.T, testCases []OperatorTestCase) {

	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	for i, tc := range testCases {
		fmt.Println(i)

		// ClearDatabase ------------------------------------------------------------------------------------

		err := tc.Cleaner.Clean()
		require.NoError(t, err)

		// test Create --------------------------------------------------------------------------------------

		var ids [numRepeats]ID

		if tc.ExpectedCreateErr {
			_, err = tc.Save("", tc.ToCreate)
			require.Error(t, err)
			continue
		}

		for i := 0; i < numRepeats; i++ {
			ids[i], err = tc.Save("", tc.ToCreate)
			require.NoError(t, err)
			require.NotEmpty(t, ids[i])
		}

		// test Read ----------------------------------------------------------------------------------------

		if tc.ExpectedReadErr {
			_, err = tc.Read(ids[toReadI])
			require.Error(t, err)
			continue
		}

		readed, err := tc.Read(ids[toReadI])
		require.NoError(t, err)
		require.NotNil(t, readed)
		testData(t, tc, tc.ToCreate, readed, ids[toReadI])

		// test List ----------------------------------------------------------------------------------------

		if tc.ExpectedListErr {
			listed, err := tc.List(selectors.Options{})
			require.Equal(t, 0, len(listed))
			require.Error(t, err)
			continue
		}

		listed, err := tc.List(selectors.Options{})
		require.NoError(t, err)
		require.Equal(t, len(ids), len(listed))

		for i, l := range listed {
			testData(t, tc, tc.ToCreate, l, ids[i])
		}

		// test Update --------------------------------------------------------------------------------------

		if tc.ExpectedUpdateErr {
			_, err = tc.Save(ids[toUpdateI], tc.ToUpdate)
			require.Error(t, err)
			continue
		}

		for i := 0; i < 2; i++ {
			_, err = tc.Save(ids[toUpdateI], tc.ToUpdate)
			require.NoError(t, err)

			readed, err = tc.Read(ids[toUpdateI])
			require.NoError(t, err)
			testData(t, tc, tc.ToUpdate, readed, ids[toUpdateI])
		}

		//// can't update absent record
		//toUpdate[keyFields[0]] += "123"
		//nativeToUpdate, err = tc.StringMapToNative(toUpdate)
		//require.NoError(t, err)
		//err = tc.Update(tc.ISToUpdate, id[toUpdateI], nativeToUpdate)
		//require.Error(t, err)

		// test DeleteList --------------------------------------------------------------------------------------

		if tc.ExpectedDeleteErr {
			err = tc.Delete(ids[toUpdateI])
			require.Error(t, err)

			readed, err = tc.Read(ids[toUpdateI])
			require.NoError(t, err)
			testData(t, tc, tc.ToUpdate, readed, ids[toUpdateI])
			continue
		}

		err = tc.Delete(ids[toUpdateI])
		require.NoError(t, err)

		readed, err = tc.Read(ids[toUpdateI])
		require.Error(t, err)
		require.Nil(t, readed)
	}
}

func testData(t *testing.T, op Operator, expectedData, data interface{}, expectedID ID) {
	if expectedData == nil {
		require.Nil(t, data)
		return
	}
	require.NotNil(t, data)

	err := op.CheckIfEqual(expectedData, expectedID, data)
	require.NoError(t, err)
}
