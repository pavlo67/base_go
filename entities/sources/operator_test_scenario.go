package sources

// import (
// 	"os"
// 	"testing"
// 	"time"

// 	"github.com/stretchr/testify/require"

// 	"encoding/json"

// 	"github.com/pavlo67/common/common"
// 	"github.com/pavlo67/common/common/crud"
// 	"github.com/pavlo67/common/common/identity"
// 	"github.com/pavlo67/common/common/logger"
// 	"github.com/pavlo67/common/common/selectors"
// )

// type OperatorTestCase struct {
// 	ActorKey
// 	db.Cleaner

// 	ToSave       Pack
// 	ToAddHistory []crud.ActionKey
// }

// var createdAt = time.Now().UTC()

// func TestCases(PacksOp ActorKey, cleanerOp db.Cleaner) []OperatorTestCase {
// 	return []OperatorTestCase{
// 		{
// 			ActorKey: PacksOp,
// 			Cleaner:  cleanerOp,
// 			ToSave: Pack{
// 				Key:     "test_key1",
// 				From:    "test_key2",
// 				To:      identity.Key("qwerwqer"),
// 				Options: common.Map{"1": float64(2)},
// 				TypeKey: "no_type",
// 				Content: map[string]string{"6": ";klj"},
// 				History: []crud.ActionKey{{
// 					Key:    crud.CreatedAction,
// 					DoneAt: createdAt,
// 				}},
// 			},

// 			ToAddHistory: []crud.ActionKey{{
// 				Key:    "action1",
// 				DoneAt: time.Now().UTC(),
// 			}},
// 		},
// 	}
// }

// const numRepeats = 3
// const toReadI = 0       // must be < numRepeats
// const toAddHistoryI = 1 // must be < numRepeats
// const toDeleteI = 2     // must be < numRepeats

// func ChechReaded(t *testing.T, readedPtr *Item, expectedID common.Key, expected Pack, l logger.ActorKey) {
// 	require.NotNil(t, readedPtr)

// 	readed := *readedPtr

// 	l.Infof("was saved: %#v", expected)
// 	l.Infof("is readed: %#v", readed.Pack)

// 	expectedBytes, _ := json.Marshal(expected.Content)
// 	require.Equal(t, expectedBytes, readed.Pack.ContentRaw)

// 	expected.Content = nil
// 	expected.ContentRaw = nil

// 	readed.Content = nil
// 	readed.ContentRaw = nil

// 	require.True(t, len(readed.History) > 0)
// 	require.True(t, readed.History[0].DoneAt.After(time.Time{}))
// 	require.True(t, readed.History[0].DoneAt.Before(time.Now()))

// 	readed.History = nil
// 	expected.History = nil

// 	require.Equal(t, expected, readed.Pack)
// 	require.Equal(t, expectedID, readed.Key)

// }

// func OperatorTestScenario(t *testing.T, testCases []OperatorTestCase, l logger.ActorKey) {
// 	if env, ok := os.LookupEnv("ENV"); !ok || env != "test" {
// 		t.Fatal("No test environment!!!")
// 	}

// 	for i, tc := range testCases {
// 		l.Debug(i)

// 		var id [numRepeats]common.Key
// 		var toSave [numRepeats]Pack

// 		// ClearDatabase ---------------------------------------------------------------------------------

// 		err := tc.Cleaner.Clean(nil, nil)
// 		require.NoError(t, err, "what is the error on .Cleaner()?")

// 		// test .Save ------------------------------------------------------------------------------------

// 		for i := 0; i < numRepeats; i++ {
// 			toSave[i] = tc.ToSave
// 			idI, err := tc.Save(&toSave[i], nil)
// 			require.NoError(t, err)
// 			require.NotEqual(t, common.Key(""), idI)
// 			id[i] = idI
// 		}

// 		// test .Read ------------------------------------------------------------------------------------

// 		readedSaved, err := tc.Read(id[toReadI], nil)
// 		require.NoError(t, err)

// 		ChechReaded(t, readedSaved, id[toReadI], tc.ToSave, l)

// 		// test .SetResults & .Read ----------------------------------------------------------------------

// 		_, err = tc.AddHistory(id[toAddHistoryI], tc.ToAddHistory, nil)
// 		require.NoError(t, err)

// 		readedUpdated, err := tc.Read(id[toAddHistoryI], nil)
// 		require.NoError(t, err)
// 		require.NotNil(t, readedUpdated)
// 		for i, r := range readedUpdated.History {
// 			readedUpdated.History[i].DoneAt = r.DoneAt.UTC()
// 		}

// 		ChechReaded(t, readedUpdated, id[toAddHistoryI], tc.ToSave, l)

// 		lenHistory0 := len(tc.ToSave.History)

// 		require.Equal(t, lenHistory0+len(tc.ToAddHistory), len(readedUpdated.History))
// 		require.Equal(t, tc.ToSave.History, readedUpdated.History[:lenHistory0])
// 		require.Equal(t, tc.ToAddHistory, readedUpdated.History[lenHistory0:])

// 		//require.True(t, reflect.DeepEqual(tc.ToAddHistory, readedUpdated.History[1:]), fmt.Sprintf("\nexpected = %#v\n  readed = %#v", tc.ToSetResults, readedUpdated.Results[0]))

// 		//// test .SetResults & .Read again  ---------------------------------------------------------------
// 		//
// 		//err = tc.Finish(id[toAddHistoryI], tc.ToSetResults, nil)
// 		//require.NoError(t, err)
// 		//
// 		//readedUpdated2, err := tc.Read(id[toAddHistoryI], nil)
// 		//require.NoError(t, err)
// 		//require.NotNil(t, readedUpdated2)
// 		//for i, r := range readedUpdated2.Results {
// 		//	readedUpdated2.Results[i].StartedAt = r.StartedAt.UTC()
// 		//	readedUpdated2.Results[i].FinishedAt = r.FinishedAt.UTC()
// 		//}
// 		//
// 		//ChechReaded(t, readedUpdated2, id[toAddHistoryI], tc.ToSave, l)
// 		//require.Equal(t, 2, len(readedUpdated2.Results))
// 		//require.True(t, reflect.DeepEqual(tc.ToSetResults, readedUpdated2.Results[0]))
// 		//require.True(t, reflect.DeepEqual(tc.ToSetResults, readedUpdated2.Results[1]))
// 		//
// 		//// check if another records are unchanged
// 		//
// 		//readedSaved, err = tc.Read(id[toReadI], nil)
// 		//require.NoError(t, err)
// 		//
// 		//ChechReaded(t, readedSaved, id[toReadI], tc.ToSave, l)
// 		//require.Equal(t, 0, len(readedSaved.Results))
// 		//
// 		//readedSaved, err = tc.Read(id[toDeleteI], nil)
// 		//require.NoError(t, err)
// 		//
// 		//ChechReaded(t, readedSaved, id[toDeleteI], tc.ToSave, l)
// 		//require.Equal(t, 0, len(readedSaved.Results))

// 		// test List -------------------------------------------------------------------------------------

// 		itemsAll, err := tc.List(nil, &crud.GetOptions{OrderBy: []string{"id"}})
// 		require.NoError(t, err)
// 		require.True(t, len(itemsAll) == numRepeats)

// 		ChechReaded(t, &itemsAll[toReadI], id[toReadI], tc.ToSave, l)
// 		ChechReaded(t, &itemsAll[toAddHistoryI], id[toAddHistoryI], tc.ToSave, l)
// 		ChechReaded(t, &itemsAll[toDeleteI], id[toDeleteI], tc.ToSave, l)

// 		itemsOne, err := tc.List(selectors.In("id", id[toDeleteI]), nil)
// 		require.NoError(t, err)
// 		require.True(t, len(itemsOne) == 1)

// 		ChechReaded(t, &itemsOne[0], id[toDeleteI], tc.ToSave, l)

// 		// test .Remove ----------------------------------------------------------------------------------

// 		//err = tc.Remove(id[toDeleteI], nil)
// 		//require.NoError(t, err)
// 		//
// 		//readDeleted, err := tc.Read(id[toDeleteI], nil)
// 		//require.ErrStr(t, err)
// 		//require.Nil(t, readDeleted)
// 		//
// 		//itemsAll, err = tc.ListTags(nil, nil)
// 		//require.NoError(t, err)
// 		//require.True(t, len(itemsAll) == numRepeats-1)
// 	}
// }
