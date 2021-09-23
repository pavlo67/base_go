package persons_pg

//import (
//	"os"
//	"testing"
//
//	"github.com/stretchr/testify/require"
//
//	"github.com/pavlo67/common/common/config"
//	"github.com/pavlo67/common/common/libraries/filelib"
//	"github.com/pavlo67/common/common/logger"
//	"github.com/pavlo67/common/common/serializer"
//
//	"github.com/pavlo67/workshop/components/users"
//)
//
//const serviceName = "gatherer"
//
//func TestCRUD(t *testing.T) {
//	env := "test"
//	err := os.Setenv("ENV", env)
//	require.NoError(t, err)
//
//	l, err = logger.Init(logger.Config{})
//	require.NoError(t, err)
//	require.NotNil(t, l)
//
//	configPath := filelib.CurrentPath() + "../../../environments/" + serviceName + "." + env + ".yaml"
//	cfg, err := config.Get(configPath, serializer.MarshalerYAML)
//	require.NoError(t, err)
//	require.NotNil(t, cfg)
//
//	cfgPg := config.Access{}
//	err = cfg.Value("pg", &cfgPg)
//	require.NoError(t, err)
//
//	l.Infof("%#v", cfgPg)
//
//	//taggerOp, taggerCleanerOp, err := tagger_sqlite.New(cfgPg, "")
//	//require.NoError(t, err)
//
//	usersOp, cleanerOp, err := New(cfgPg, "storage", "", nil, nil) // taggerOp, taggerCleanerOp
//	require.NoError(t, err)
//	require.NotNil(t, usersOp)
//	require.NotNil(t, cleanerOp)
//
//	testCases := users.TestCases(usersOp, cleanerOp)
//
//	users.OperatorTestScenario(t, testCases, l)
//}
