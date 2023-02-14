package persons_files

import (
	"testing"
)

func TestPersonsFiles(t *testing.T) {
	//_, cfgService, l := config.Prepare(
	//	t,
	//	"../../../_environments/",
	//	"test",
	//	"", // "persons_test."+strconv.FormatInt(time.Now().Unix(), 10)+".log",
	//)
	//
	//components := []starter.Starter{
	//	{files_fs.Starter(), nil},
	//	{Starter(), nil},
	//}
	//
	//label := "PERSONS_FS/TEST BUILD"
	//joinerOp, err := starter.Run(components, cfgService, label, l)
	//if err != nil {
	//	l.Fatal(err)
	//}
	//defer joinerOp.CloseAll()
	//
	//personsOp, _ := joinerOp.Interface(persons.InterfaceKey).(persons.Operator)
	//require.NotNil(t, personsOp)
	//
	//personsCleanerOp, _ := joinerOp.Interface(persons.InterfaceCleanerKey).(db.Cleaner)
	//require.NotNil(t, personsCleanerOp)
	//
	//persons.OperatorTestScenario(t, personsOp, personsCleanerOp)
}
