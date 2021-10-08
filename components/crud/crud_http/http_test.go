package crud_http

//func TestHTTPRecordsCRUD(t *testing.T) {
//	cfgService, l := config.PrepareTests(t, "../../../_environments/", "test", "dispatcher_records01_pg.log")
//	require.NotNil(t, cfgService)
//
//	components := append(
//		node_crud_settings.Components(true),
//		starter.Starter{Starter(), common.Map{
//			"prefix":        crud_node_http.PrefixREST,
//			"server_config": crud_node_http.ServerConfig,
//		}},
//	)
//
//	joinerOp, err := starter.Run(components, &cfgService, "CLI BUILD FOR TEST", l)
//	require.NoError(t, err)
//	require.NotNil(t, joinerOp)
//	defer joinerOp.CloseAll()
//
//	time.Sleep(time.Second)
//
//	recordsCleanerOp, _ := joinerOp.Interface(records01.InterfaceCleanerKey).(db.Cleaner)
//	require.NotNil(t, recordsCleanerOp)
//
//	crudOp, _ := joinerOp.Interface(InterfaceKey).(crud.Operator)
//	require.NotNil(t, crudOp)
//
//	crudData := crud.Data{
//		Key: crud.Key{
//			Type: records01.CRUD01,
//			ID:   records01.TestItem.ID,
//		},
//		Description: records01.TestItem.Description
