package cfg

//import (
//	"log"
//	"os"
//	"testing"
//
//	"github.com/pavlo67/common/common/config"
//	"github.com/pavlo67/common/common/filelib"
//
//	"github.com/stretchr/testify/require"
//)
//
//func TestMain(m *testing.M) {
//	if err := os.Setenv("ENV", "test"); err != nil {
//		log.Fatalln("No test environment!!!")
//	}
//	os.Exit(m.Run())
//}
//
//func TestReadFile(t *testing.T) {
//
//	ExpectedStrings := map[string]string{
//		"test1": "1.000",
//		"test2": "test2",
//	}
//
//	ExpectedFlags := map[string]bool{
//		"debug": true,
//	}
//
//	ExpectedMysqls := map[string]config.Access{
//		"notebook": {
//			Host: "localhost",
//			Port: 3306,
//			User: "root",
//			Pass: "",
//			Path: "notebook_go_test",
//		},
//	}
//
//	//ExpectedIndex := map[string]serverhttp_jschmhr.ComponentsIndex{
//	//	"main": {
//	//		EndpointsConfig: map[string]serverhttp_jschmhr.Endpoint{
//	//			"ep1": {Method: "GET", LocalPath: "/"},
//	//		},
//	//		Listeners: map[string]serverhttp_jschmhr.Listener{
//	//			"lst1": {Label: "id1"},
//	//		},
//	//	},
//	//}
//
//	cfg, err := config.Get("cfg_test.json5")
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	for k, es := range ExpectedStrings {
//		s, errs := cfg.String(k, nil)
//		if len(errs) > 0 {
//			log.Fatalf("unexpected errs (%v) for string key: %s", errs, k)
//		}
//		require.Equal(t, s, es, "bad string for key: "+k)
//
//		kUnexp := k + "unexpected"
//		_, errs = cfg.String(kUnexp, nil)
//		if len(errs) < 1 {
//			log.Fatalf("no expected errs for string key: %s", kUnexp)
//		}
//	}
//
//	for k, ef := range ExpectedFlags {
//		f, errs := cfg.Bool(k, nil)
//		if len(errs) > 0 {
//			log.Fatalf("unexpected errs (%v) for bool flag key: %s", errs, k)
//		}
//		require.Equal(t, f, ef, "bad bool flag for key: "+k)
//
//		kUnexp := k + "unexpected"
//		_, errs = cfg.Bool(kUnexp, nil)
//		if len(errs) < 1 {
//			log.Fatalf("no expected errs for bool flag key: %s", kUnexp)
//		}
//	}
//
//	for k, emc := range ExpectedMysqls {
//		mc, errs := cfg.MySQL(k, nil)
//		if len(errs) > 0 {
//			log.Fatalf("unexpected errs (%v) for mysql key: %s", errs, k)
//		}
//		require.Equal(t, mc, emc, "bad mysql credentials for key: "+k)
//
//		kUnexp := k + "unexpected"
//		_, errs = cfg.MySQL(kUnexp, nil)
//		if len(errs) < 1 {
//			log.Fatalf("no expected errs for mysql credentials key: %s", kUnexp)
//		}
//	}
//
//}
