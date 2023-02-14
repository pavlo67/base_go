package test

//func TestCRUD(t *testing.T) {
//
//	//t.Skip()
//
//	conf, err := config.Get(basis.CurrentPath() + "../../../cfg.json5")
//	if err != nil {
//		log.Fatal(err)
//	}
//	if conf == nil {
//		log.Fatal(errors.New("no config data after setup.Init()"))
//	}
//	componentKeys := config.ComponentKeys{
//		"mysql": "notebook",
//	}
//	mysqlConfig, errs := conf.MySQL(componentKeys["mysql"], nil)
//	err = basis.JoinErrors(errs)
//	if err != nil {
//		log.Fatal(err)
//	}
//	//err := config.LoadContext("../../../cfg.json5")
//	//if err != nil {
//	//	log.Fatal(err)
//	//}
//	//// creating Operator
//	//mysqlConfig, ok := config.Mysql["notebook"]
//	//if !ok {
//	//	log.Fatal(errors.Errorf("no mysql[notebook] section in config: %v", config.Mysql))
//	//}
//
//	domain := "aaa"
//	identity := confidenter.Identity{domain, "user", "111", ""}
//	identityBAD := confidenter.Identity{"aaa", "user", "999", ""}
//	//isBAD := identityBAD.String()
//	IS := identity.String()
//	//identityGroup := confidenter.Identity{"aaa", "group", "222", ""}
//	//isGroup := identityGroup.String()
//	//testController, err := controller.NewCRUDController(
//	//	map[confidenter.Identity]confidenter.Identity{
//	//		IS:    isGroup,
//	//		isBAD: isGroup,
//	//	},
//	//)
//
//	mysqlGroup, _ := groupsmysql.NewGroupsMySQL(
//		identity,
//		"aaa",
//		mysqlConfig,
//		"group",
//		"group_member",
//		controller.Managers{},
//	)
//
//	fountmysql, err := fountsmysql.NewMySQLFount(
//		//testController,
//		mysqlGroup,
//		mysqlConfig,
//		"fount",
//		"fount_tags",
//		"fount_stat",
//		"scanner_stat",
//		controller.Managers{rights.Create: IS, rights.Change: IS, rights.View: IS, rights.Delete: IS},
//	)
//	if err != nil {
//		t.Fatalf("can't init MySQLFount for tests: %v", err)
//	}
//	//fountsmysqlCRUD := fount.FountCRUD{fountsmysql}
//
//	s1 := rand.NewSource(time.Now().UnixNano())
//	r1 := rand.New(s1)
//
//	descriptionFount := crud.Records{
//		ID: "",
//		Details: map[string]string{
//			"Label": "Label CRUD TEST",
//			"Url":   "https://rss.unian.net/site/news_ukr.rss#" + strconv.Itoa(r1.Intn(1000)),
//		},
//		Managers: controller.Managers{rights.Owner: IS, rights.View: IS},
//	}
//
//	testCases := []crud.OperatorNewTestCase{
//		{
//			//&fountsmysqlCRUD,
//			fountmysql,
//			identity,
//			confidenter.Identity{},
//			confidenter.Identity{},
//			identityBAD,
//			identityBAD,
//			descriptionFount,
//			"Label",
//			false,
//			nil,
//			nil,
//			false,
//			false,
//		},
//	}
//	for _, testCase := range testCases {
//		fmt.Println("\n =============== Operator TEST: fount =================")
//		crud.OperatorNewTest(t, testCase)
//	}
//}
