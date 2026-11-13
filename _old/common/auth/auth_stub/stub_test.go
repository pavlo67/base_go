package auth_stub

import (
	"os"
	"testing"

	"github.com/pavlo67/common/common/logger/logger_test"

	"github.com/pavlo67/common/common/auth"
	"github.com/stretchr/testify/require"
)

func TestAuthStub(t *testing.T) {

	os.Setenv("ENV", "test")

	l = logger_test.New(t, "", "", false, nil)

	authOp, err := New(nil)
	require.NoError(t, err)
	require.NotNil(t, authOp)

	auth.OperatorTestScenarioPassword(t, authOp)
}
