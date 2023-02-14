package fileslocal

import (
	"log"
	"testing"

	"github.com/pavlo67/punctum/basis/filelib"
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"

	"github.com/pavlo67/partes/crud"
	"github.com/pavlo67/punctum/confidenter/controller"
	"github.com/pavlo67/punctum/confidenter/groups/groupsmysql"
	"github.com/pavlo67/punctum/things_old/fileinfo/fileinfomysql"
	"github.com/pavlo67/punctum/things_old/files"
)

var PartKeys = config.PartKeys{
	"identity":   "notebook",
	"fileslocal": "test",
	"mysql":      "notebook",
}

func TestCRUD(t *testing.T) {
	conf, err := joiner.Init(filelib.CurrentPath()+"../../../cfg.json5", PartKeys)
	if err != nil {
		t.Fatal(err)
	}
	if conf == nil {
		t.Fatal("no config data after setup.Run()")
	}

	identity, identityAnother, _, err := controller.IdentitiesForTests()
	if err != nil {
		t.Fatalf("on crud.IdentitiesForTestsOld(): %s", err)
	}

	starters := []starter.Starter{
		{groupsmysql.Starter(), ""},
		{fileinfomysql.Starter(), ""},
		{Starter(), ""},
	}
	closers, err := starter.Run(conf, PartKeys, starters, "TEST BUILD", false, false)
	for _, closer := range closers {
		defer closer.Close()
	}
	if err != nil {
		log.Println(err)
	}

	filesOp, ok := joiner.Component(files.InterfaceKey).(files.Operator)
	if !ok {
		t.Fatalf("no fileinfo.Operator found for crud test")
	}

	operatorCRUD := files.OperatorCRUD{filesOp}
	testCases, err := operatorCRUD.TestCases(identity, identityAnother)
	if err != nil {
		t.Fatalf("can't operatorCRUD.TestCases(%#v, %#v): %s", identity, identityAnother, err)
	}

	crud.OperatorTest(t, testCases, true)
}
