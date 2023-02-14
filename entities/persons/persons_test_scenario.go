package persons

import (
	"os"
	"testing"

	"github.com/pavlo67/common/common"

	"github.com/pavlo67/common/common/db"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/rbac"
)

func OperatorTestScenario(t *testing.T, personsOp Operator, personsCleanerOp db.Cleaner) {

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

	passwordToSave := "passwordToSave"
	personToSave := Item{
		Identity: auth.Identity{
			Nickname: "test_nickname1",
			Roles:    rbac.Roles{rbac.RoleUser},
		},
	}
	err = personToSave.SetCreds(auth.Creds{auth.CredsPassword: passwordToSave})
	require.NoError(t, err)

	personToSave.ID, err = personsOp.Save(personToSave, adminIdentity)
	require.NoErrorf(t, err, "%#v", err)
	require.NotEmpty(t, personToSave.ID)

	countPersons(t, personsOp, adminIdentity, 1)

	// read person ---------------------------------------------

	personReaded, err := personsOp.Read(personToSave.ID, adminIdentity)
	require.NoErrorf(t, err, "%#v", err)
	require.True(t, personReaded.CheckCreds(auth.CredsPassword, passwordToSave))

	require.Equal(t, personToSave.Identity, personReaded.Identity)
	require.Equal(t, personToSave.Info, personReaded.Info)

	// change person -------------------------------------------

	personToChange := personReaded
	personToChange.Nickname += " (changed)"
	personToChange.Nickname += " (changed)"
	if personToChange.Info == nil {
		personToChange.Info = common.Map{}
	}

	personToChange.Info["change"] = "change"
	passwordToChange := "passwordToChange"
	err = personToChange.SetCreds(auth.Creds{auth.CredsPassword: passwordToChange})
	require.NoError(t, err)

	id, err := personsOp.Save(*personToChange, adminIdentity)
	require.NoErrorf(t, err, "%#v", err)
	require.Equal(t, personToChange.ID, id)

	countPersons(t, personsOp, adminIdentity, 1)

	// read person ---------------------------------------------

	personChangedReaded, err := personsOp.Read(personToSave.ID, adminIdentity)
	require.NoErrorf(t, err, "%#v", err)
	require.Truef(t, personChangedReaded.CheckCreds(auth.CredsPassword, passwordToChange), ".Creds(): %#v", personChangedReaded.Creds())

	require.Equal(t, personToChange.Identity, personChangedReaded.Identity)
	require.Equal(t, personToChange.Info, personChangedReaded.Info)

	// add another person --------------------------------------

	passwordToSaveAnother := "passwordToSaveAnother"
	personToSaveAnother := Item{
		Identity: auth.Identity{
			Nickname: "test_nickname2",
			Roles:    rbac.Roles{rbac.RoleUser},
		},
	}
	err = personToSaveAnother.SetCreds(auth.Creds{auth.CredsPassword: passwordToSaveAnother})
	require.NoError(t, err)

	personToSaveAnother.ID, err = personsOp.Save(personToSaveAnother, adminIdentity)
	require.NoErrorf(t, err, "%#v", err)
	require.NotEmpty(t, personToSaveAnother.ID)

	countPersons(t, personsOp, adminIdentity, 2)

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

	err = personsOp.Remove(personToSaveAnother.ID, adminIdentity)
	require.NoErrorf(t, err, "%#v", err)

	personAnotherReaded, err := personsOp.Read(personToSaveAnother.ID, adminIdentity)
	require.Errorf(t, err, "%#v", err)
	require.Nil(t, personAnotherReaded)

	countPersons(t, personsOp, adminIdentity, 1)

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

func countPersons(t *testing.T, personsOp Operator, identity *auth.Identity, expectedCount int) {
	personItems, err := personsOp.List(nil, identity)
	require.NoErrorf(t, err, "%#v", err)
	require.Equalf(t, expectedCount, len(personItems), "%#v", personItems)

}
