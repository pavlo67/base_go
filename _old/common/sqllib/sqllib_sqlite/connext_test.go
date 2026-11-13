package sqllib_sqlite

//func TestConnect(t *testing.T) {
//	return
//
//	envs, _ := config.PrepareTests(
//		t,
//		"../../../_envs/",
//		"", // "connect_test."+strconv.FormatInt(time.Now().Unix(), 10)+".log",
//	)
//
//	//var cfgSqlite config.Access
//	//err := envs.Value("sqllib_sqlite", &cfgSqlite)
//	//require.NoError(t, err)
//	//
//	//cfgSqlite.Path, err = filelib.Dir(cfgSqlite.Path)
//	//require.NoError(t, err)
//	//require.NotEmpty(t, cfgSqlite.Path)
//	//
//	//cfgSqlite.Path += "test_connect.sqlite"
//
//	var cfgSqlite config.Access
//	err := envs.Value("db_sqlite", &cfgSqlite)
//	require.NoError(t, err)
//
//	db, err := Connect(cfgSqlite)
//	require.NoError(t, err)
//	require.NotNil(t, db)
//
//	sqllib.TestDB(t, db)
//
//}
