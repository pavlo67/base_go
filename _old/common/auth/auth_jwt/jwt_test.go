package auth_jwt

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/auth"
)

func TestOperator(t *testing.T) {
	env := "test"
	err := os.Setenv("ENV", env)
	require.NoError(t, err)

	//l, err = logger_zap.Init(logger.Config{})
	//require.NoError(t, err)
	//require.NotNil(t, l)

	authOp, err := New("key.test")
	require.NoError(t, err)
	require.NotNil(t, authOp)

	auth.OperatorTestScenarioToken(t, authOp)
}
