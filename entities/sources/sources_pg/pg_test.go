package sources_pg

// import (
// 	"os"
// 	"testing"

// 	"github.com/stretchr/testify/require"

// 	"github.com/pavlo67/common/common/config"
// 	"github.com/pavlo67/common/common/libraries/filelib"
// 	"github.com/pavlo67/common/common/logger"
// 	"github.com/pavlo67/common/common/serializer"

// 	"github.com/pavlo67/workshop/components/packs"
// )

// const serviceName = "gatherer"

// func TestCRUD(t *testing.T) {
// 	env := "test"
// 	err := os.Setenv("ENV", env)
// 	require.NoError(t, err)

// 	l, err = logger.Init(logger.Config{})
// 	require.NoError(t, err)
// 	require.NotNil(t, l)

// 	configPath := filelib.CurrentPath() + "../../../environments/" + serviceName + "." + env + ".yaml"
// 	cfg, err := config.Get(configPath, serializer.MarshalerYAML)
// 	require.NoError(t, err)
// 	require.NotNil(t, cfg)

// 	cfgPostgres := config.Access{}
// 	err = cfg.Value("pg", &cfgPostgres)
// 	require.NoError(t, err)

// 	l.Infof("%#v", cfgPostgres)

// 	packsOp, cleanerOp, err := New(cfgPostgres, "", "")
// 	require.NoError(t, err)

// 	testCases := packs.TestCases(packsOp, cleanerOp)

// 	packs.OperatorTestScenario(t, testCases, l)
// }
