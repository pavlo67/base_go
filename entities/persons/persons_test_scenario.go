package persons

import (
	"os"
	"testing"

	"github.com/pavlo67/common/common/joiner"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/rbac"

	"github.com/pavlo67/data/elements/contacts"
	"github.com/pavlo67/data/types"
)

func OperatorTestScenario(t *testing.T, joinerOp joiner.Operator) {

	personsOp, _ := joinerOp.Interface(InterfaceKey).(Operator)
	require.NotNil(t, personsOp)

	personsCleanerOp, _ := joinerOp.Interface(InterfaceCleanerKey).(db.Cleaner)
	require.NotNil(t, personsCleanerOp)

	var err error

	// prepare... ----------------------------------------------

	require.Equal(t, "test", os.Getenv("ENV"))

	require.NotNil(t, personsOp)
	require.NotNil(t, personsCleanerOp)

	adminIdentity := auth.IdentityWithRoles(rbac.RoleAdmin)
	require.NotNil(t, adminIdentity)

	// clean old data ------------------------------------------

	err = personsCleanerOp.Clean()
	require.NoError(t, err)

	countPersons(t, personsOp, adminIdentity, 0)

	// add person ----------------------------------------------

	personToSave := types.Person01{
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

	savedID, err := personsOp.Save(Item01{Person01: personToSave}, adminIdentity)
	require.NoError(t, err)
	require.NotEmpty(t, savedID)

	countPersons(t, personsOp, adminIdentity, 1)

	// read person ---------------------------------------------

	personReaded, err := personsOp.Read(savedID, adminIdentity)
	require.NoError(t, err)
	require.NotNil(t, personReaded)
	checkPerson(t, personToSave, personReaded.Person01)

	// change person -------------------------------------------

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

	savedChangedID, err := personsOp.Save(*personToSaveChanged, adminIdentity)
	require.NoError(t, err)
	require.Equal(t, personToSaveChanged.ID, savedChangedID)

	countPersons(t, personsOp, adminIdentity, 1)

	// read person ---------------------------------------------

	personReaded, err = personsOp.Read(savedID, adminIdentity)
	require.NoError(t, err)
	require.NotNil(t, personReaded)
	require.Equal(t, personReaded.ID, savedID)
	checkPerson(t, personToSaveChanged.Person01, personReaded.Person01)

	// remove person -------------------------------------------

	err = personsOp.Remove(savedID, adminIdentity)
	require.NoError(t, err)
	personReaded, err = personsOp.Read(savedID, adminIdentity)
	require.Error(t, err)
	require.Nil(t, personReaded)
	countPersons(t, personsOp, adminIdentity, 0)

	//// add another person --------------------------------------
	//
	//passwordToSaveAnother := "passwordToSaveAnother"
	//personToSaveAnother := types.Person01{
	//	Identity: auth.Identity{
	//		Nickname: "test_nickname2",
	//		Roles:    rbac.Roles{rbac.RoleUser},
	//	},
	//}
	//err = personToSaveAnother.SetCreds(auth.Creds{auth.CredsPassword: passwordToSaveAnother})
	//require.NoError(t, err)
	//
	//personToSaveAnother.ID, err = personsOp.Save(personToSaveAnother, adminIdentity)
	//require.NoErrorf(t, err, "%#v", err)
	//require.NotEmpty(t, personToSaveAnother.ID)
	//
	//countPersons(t, personsOp, adminIdentity, 2)

	//// add person ----------------------------------------------
	//
	//personItems, err = personsOp.List(adminIdentity)
	//require.NoErrorf(t, err, "%#v", err)
	//require.Equal(t, 3, len(personItems))
	//
	//// list persons by itself: error ---------------------------
	//
	//personItems, err = personsOp.List(&person1Options)
	//require.Errorf(t, err, "%#v", err)
	//require.Empty(t, personItems)

	//// change person by admin: ok ------------------------------
	//
	//person1ToChange := *personReaded
	//person1ToChange.Nickname += "_changed"
	//
	//person1Changed, err := personsOp.Change(person1ToChange, adminIdentity)
	//require.NoErrorf(t, err, "%#v", err)
	//require.Equal(t, person1ToChange.Identity, person1Changed.Identity)
	//
	//person1ChangedReaded, err := personsOp.Read(person1Changed.ID, adminIdentity)
	//require.NoErrorf(t, err, "%#v", err)
	//require.Equal(t, person1ToChange.Identity, person1ChangedReaded.Identity)
	//
	//// change person by itself: ok -----------------------------
	//
	//person1ToChange.Nickname += "_again"
	//
	//person1Changed, err = personsOp.Change(person1ToChange, &person1Options)
	//require.NoErrorf(t, err, "%#v", err)
	//require.Equal(t, person1ToChange.Identity, person1Changed.Identity)
	//
	//person1ChangedReaded, err = personsOp.Read(person1Changed.ID, &person1Options)
	//require.NoErrorf(t, err, "%#v", err)
	//require.Equal(t, person1ToChange.Identity, person1ChangedReaded.Identity)
	//
	//// change/read person by another person: error -------------
	//
	//person1ToChangeAgain := *person1ChangedReaded
	//person1ToChangeAgain.Nickname += "_again2"
	//
	//person1ChangedWrong, err := personsOp.Change(person1ToChangeAgain, &person2Options)
	//require.Errorf(t, err, "%#v", err)
	//require.Nil(t, person1ChangedWrong)
	//
	//person1ReadedWrong, err := personsOp.Read(personID1, &person2Options)
	//require.Errorf(t, err, "%#v", err)
	//require.Nil(t, person1ReadedWrong)
	//
	//person1Readed, err := personsOp.Read(personID1, &person1Options)
	//require.NoErrorf(t, err, "%#v", err)
	//require.NotNil(t, person1Readed)
	//require.Equal(t, person1Changed.Identity, person1Readed.Identity)
	//
	// remove person by admin: ok ------------------------------

	//err = personsOp.Remove(personToSaveAnother.ID, adminIdentity)
	//require.NoErrorf(t, err, "%#v", err)
	//
	//personAnotherReaded, err := personsOp.Read(personToSaveAnother.ID, adminIdentity)
	//require.Errorf(t, err, "%#v", err)
	//require.Nil(t, personAnotherReaded)
	//
	//countPersons(t, personsOp, adminIdentity, 1)

	//// remove person by itself: ok -----------------------------
	//
	//require.NotNil(t, person2Options.Identity)
	//err = personsOp.Remove(personID2, &person2Options)
	//require.NoErrorf(t, err, "%#v / %#v", person2Options.Identity, err)
	//
	//person2Readed, err := personsOp.Read(personID2, &person2Options)
	//require.Errorf(t, err, "%#v", err)
	//require.Nil(t, person2Readed)
	//
	//// remove person by another person: error ------------------
	//
	//err = personsOp.Remove(personID1, &person2Options)
	//require.Errorf(t, err, "%#v", err)
	//
	//person1Readed, err = personsOp.Read(personID1, adminIdentity)
	//require.NoErrorf(t, err, "%#v", err)
	//require.NotNil(t, person1Readed)
	//require.Equal(t, person1ChangedReaded.Identity, person1Readed.Identity)

}

func checkPerson(t *testing.T, personExpected, personToCheck types.Person01) {
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

func countPersons(t *testing.T, personsOp Operator, identity *auth.Identity, expectedCount int) {
	personItems, err := personsOp.List(nil, identity)
	require.NoError(t, err)
	require.Equalf(t, expectedCount, len(personItems), "personItems = %#v", personItems)

}
