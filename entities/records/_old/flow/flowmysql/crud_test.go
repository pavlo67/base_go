package flowmysql

import (
	"os"

	"log"
	"testing"
)

func TestMain(m *testing.M) {
	if _, ok := os.LookupEnv("TEST_ENVIRONMENT"); !ok {
		log.Fatalln("No test environment!!!")
	}
	os.Exit(m.Run())
}

//func TestCRUD(t *testing.T) {
//
//	//t.Skip()
//	// fragment table structure must be corrected
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
//
//	domain := "aaa"
//	identity := confidenter.Identity{domain, "user", "111", ""}
//	identityBAD := confidenter.Identity{"aaa", "user", "999", ""}
//	//isBAD := identityBAD.String()
//	IS := identity.String()
//	//identityGroup := confidenter.Identity{"aaa", "group", "1", ""}
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
//	fragmentMySQL, err := flowmysql.NewFlowMySQL(
//		identity,
//		//testController,
//		mysqlGroup,
//		mysqlConfig,
//		"flow",
//		controller.Managers{rights.Create: IS},
//	)
//	if err != nil {
//		t.Fatalf("can't init NewFlowMySQL for tests: %v", err)
//	}
//	//fragmentMySQLCRUD := flow.OperatorCRUD{fragmentMySQL}
//	s1 := rand.NewSource(time.Now().UnixNano())
//	r1 := rand.New(s1)
//	is := "aaa/flow/" + strconv.Itoa(r1.Intn(1000))
//	descriptionFragment := crud.Records{
//		ID: "",
//		Details: map[string]string{
//			"Content":    "Subj CRUD TEST",
//			"OriginalID": is,
//			"FountURL":   "http://test.com/rss",
//		},
//		Managers: controller.Managers{rights.Owner: IS, rights.View: IS},
//	}
//
//	testCases := []crud.OperatorNewTestCase{
//		{
//			//&fragmentMySQLCRUD,
//			fragmentMySQL,
//			identity,
//			identityBAD,
//			confidenter.Identity{},
//			identityBAD,
//			identityBAD,
//			descriptionFragment,
//			"Content",
//			false,
//			nil,
//			nil,
//			false,
//			false,
//		},
//	}
//	for _, testCase := range testCases {
//		fmt.Println("\n =============== Operator TEST: fragment =================")
//		crud.OperatorNewTest(t, testCase)
//	}
//}
