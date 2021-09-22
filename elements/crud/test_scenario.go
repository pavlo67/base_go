package crud

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/db"

	"github.com/pavlo67/data/elements/selectors"
)

type OperatorTestCase struct {
	Type
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

func OperatorTestScenario(t *testing.T, testCases []OperatorTestCase) {

	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
		t.Fatal("No test environment!!!")
	}

	for i, tc := range testCases {
		fmt.Println(i)

		// ClearDatabase ------------------------------------------------------------------------------------

		err := tc.Cleaner.Clean()
		require.NoError(t, err)

		// test Create --------------------------------------------------------------------------------------

		var keys [numRepeats]Key

		if tc.ExpectedCreateErr {
			_, err = tc.Save(Key{Type: tc.Type}, tc.ToCreate)
			require.Error(t, err)
			continue
		}

		for i := 0; i < numRepeats; i++ {
			key, err := tc.Save(Key{Type: tc.Type}, tc.ToCreate)

			require.NoError(t, err)
			require.NotNil(t, key)
			require.Equal(t, key.Type, tc.Type)
			require.NotEmpty(t, key.ID)

			keys[i] = *key

		}

		// test Read ----------------------------------------------------------------------------------------

		if tc.ExpectedReadErr {
			_, err = tc.Read(keys[toReadI])
			require.Error(t, err)
			continue
		}

		readed, err := tc.Read(keys[toReadI])
		require.NoError(t, err)
		require.NotNil(t, readed)
		testData(t, tc, tc.ToCreate, readed, keys[toReadI])

		// test List ----------------------------------------------------------------------------------------

		if tc.ExpectedListErr {
			listed, err := tc.List(selectors.Options{})
			require.Equal(t, 0, len(listed))
			require.Error(t, err)
			continue
		}

		listed, err := tc.List(selectors.Options{})
		require.NoError(t, err)
		require.Equal(t, len(keys), len(listed))

		for i, l := range listed {
			testData(t, tc, tc.ToCreate, l, keys[i])
		}

		// test Update --------------------------------------------------------------------------------------

		if tc.ExpectedUpdateErr {
			_, err = tc.Save(keys[toUpdateI], tc.ToUpdate)
			require.Error(t, err)
			continue
		}

		for i := 0; i < 2; i++ {
			_, err = tc.Save(keys[toUpdateI], tc.ToUpdate)
			require.NoError(t, err)

			readed, err = tc.Read(keys[toUpdateI])
			require.NoError(t, err)
			testData(t, tc, tc.ToUpdate, readed, keys[toUpdateI])
		}

		//// can't update absent record
		//toUpdate[keyFields[0]] += "123"
		//nativeToUpdate, err = tc.StringMapToNative(toUpdate)
		//require.NoError(t, err)
		//err = tc.Update(tc.ISToUpdate, id[toUpdateI], nativeToUpdate)
		//require.Error(t, err)

		// test DeleteList --------------------------------------------------------------------------------------

		if tc.ExpectedDeleteErr {
			err = tc.Remove(keys[toUpdateI])
			require.Error(t, err)

			readed, err = tc.Read(keys[toUpdateI])
			require.NoError(t, err)
			testData(t, tc, tc.ToUpdate, readed, keys[toUpdateI])
			continue
		}

		err = tc.Remove(keys[toUpdateI])
		require.NoError(t, err)

		readed, err = tc.Read(keys[toUpdateI])
		require.Error(t, err)
		require.Nil(t, readed)
	}
}

func testData(t *testing.T, op Operator, expectedData, data interface{}, expectedID Key) {
	if expectedData == nil {
		require.Nil(t, data)
		return
	}
	require.NotNil(t, data)

	err := op.CheckIfEqual(expectedData, expectedID, data)
	require.NoError(t, err)
}
