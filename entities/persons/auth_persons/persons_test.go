package auth_persons

import (
	"testing"
)

func TestAuthPersons(t *testing.T) {
	//_, cfgService, l := config.Prepare(
	//	t,
	//	"../../../_environments/",
	//	"test",
	//	"", // "connect_test."+strconv.FormatInt(time.Now().Unix(), 10)+".log",
	//)
	//
	//starters := []starter.Starter{
	//	{persons_files.Starter(), nil},
	//	{Starter(), nil},
	//}
	//
	//label := "CLI/TEST BUILD"
	//joinerOp, err := starter.Run(starters, cfgService, label, l)
	//if err != nil {
	//	l.Fatal(err)
	//}
	//defer joinerOp.CloseAll()
	//
	//authOp, _ := joinerOp.Interface(auth.InterfaceKey).(auth.Operator)
	//if authOp == nil {
	//	l.Fatalf("no auth.Operator with key %s", auth.InterfaceKey)
	//}
	//
	//auth.OperatorTestScenarioPassword(t, authOp)
}
