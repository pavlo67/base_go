package loader_http

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/libraries/filelib"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/serializer"

	"github.com/pavlo67/workshop/constructions/filesloader"
)

func TestHTTP(t *testing.T) {
	env := "test"
	err := os.Setenv("ENV", env)
	require.NoError(t, err)

	l, err = logger.Init(logger.Config{})
	require.NoError(t, err)
	require.NotNil(t, l)

	configPath := filelib.CurrentPath() + "../../../environments/" + env + ".yaml"
	cfg, err := config.Get(configPath, serializer.MarshalerYAML)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	pathToStore := "./test_results/"

	flOp, cleanerOp, err := New(pathToStore)
	require.NoError(t, err)

	testCases := filesloader.TestCases(flOp, cleanerOp, pathToStore)

	filesloader.OperatorTestScenario(t, testCases, l)
}
