package persons01

//// DEPRECATED
//func OperatorTestScenario(t *testing.T, joinerOp joiner.Operator, interfaceKey, interfaceCleanerKey joiner.InterfaceKey, personToSave entities.Person01) {
//
//	personsOp, _ := joinerOp.Interface(interfaceKey).(Operator)
//	require.NotNil(t, personsOp)
//
//	personsCleanerOp, _ := joinerOp.Interface(interfaceCleanerKey).(db.Cleaner)
//	require.NotNil(t, personsCleanerOp)
//
//	var err error
//
//	// prepare... ----------------------------------------------
//
//	require.Equal(t, "test", os.Getenv("ENV"))
//
//	adminIdentity := auth.IdentityWithRoles(rbac.RoleAdmin)
//	require.NotNil(t, adminIdentity)
//
//	// clean old data ------------------------------------------
//
//	err = personsCleanerOp.Clean()
//	require.NoError(t, err)
//
//	CountTestPersons(t, personsOp, adminIdentity, 0)
//
//	// add person ----------------------------------------------
//
//	itemToSave := Item{Person01: personToSave}
//
//	savedID, err := personsOp.Save(itemToSave, adminIdentity)
//	require.NoError(t, err)
//	require.NotEmpty(t, savedID)
//
//	CountTestPersons(t, personsOp, adminIdentity, 1)
//
//	// read person ---------------------------------------------
//
//	personReaded, err := personsOp.Read(savedID, adminIdentity)
//	require.NoError(t, err)
//	require.NotNil(t, personReaded)
//	CheckTestPerson(t, itemToSave, *personReaded)
//
//	// change person -------------------------------------------
//
//	personToSaveChanged := ChangeTestItem(*personReaded, savedID)
//
//	savedChangedID, err := personsOp.Save(personToSaveChanged, adminIdentity)
//	require.NoError(t, err)
//	require.Equal(t, personToSaveChanged.ID, savedChangedID)
//
//	CountTestPersons(t, personsOp, adminIdentity, 1)
//
//	// read person ---------------------------------------------
//
//	personReaded, err = personsOp.Read(savedID, adminIdentity)
//	require.NoError(t, err)
//	require.NotNil(t, personReaded)
//	require.Equal(t, personReaded.ID, savedID)
//	CheckTestPerson(t, personToSaveChanged, *personReaded)
//
//	// remove person -------------------------------------------
//
//	err = personsOp.Remove(savedID, adminIdentity)
//	require.NoError(t, err)
//	personReaded, err = personsOp.Read(savedID, adminIdentity)
//	require.Error(t, err)
//	require.Nil(t, personReaded)
//	CountTestPersons(t, personsOp, adminIdentity, 0)
//
//}
//
//// DEPRECATED
//func CountTestPersons(t *testing.T, personsOp Operator, identity *auth.Identity, expectedCount int) {
//	personItems, err := personsOp.List(nil, identity)
//	require.NoError(t, err)
//	require.Equalf(t, expectedCount, len(personItems), "personItems = %#v", personItems)
//}
//
//// DEPRECATED
//func CheckTestPerson(t *testing.T, personExpected, personToCheck Item) {
//	descriptionExpected, descriptionToCheck := personExpected.Description, personToCheck.Description
//	personExpected.Description, personToCheck.Description = crud.Description{}, crud.Description{}
//
//	require.Equal(t, personExpected, personToCheck)
//	require.Equal(t, descriptionExpected.URN, descriptionToCheck.URN)
//
//	if len(descriptionExpected.Tags) > 0 {
//		require.Equal(t, descriptionExpected.Tags, descriptionToCheck.Tags)
//	} else {
//		require.Equal(t, 0, len(descriptionToCheck.Tags))
//	}
//	if len(descriptionExpected.RelationsMap) > 0 {
//		require.Equal(t, descriptionExpected.RelationsMap, descriptionToCheck.RelationsMap)
//	} else {
//		require.Equal(t, 0, len(descriptionToCheck.RelationsMap))
//	}
//
//	require.Equal(t, descriptionExpected.ViewerNSS, descriptionToCheck.ViewerNSS)
//	require.Equal(t, descriptionExpected.OwnerNSS, descriptionToCheck.OwnerNSS)
//
//	require.True(t, len(descriptionToCheck.History) >= len(descriptionExpected.History))
//	require.Equal(t, descriptionExpected.History, descriptionToCheck.History[:len(descriptionExpected.History)])
//}
