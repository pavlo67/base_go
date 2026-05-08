package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetCallInfo(t *testing.T) {
	callInfo := GetCallInfo()

	require.Equal(t, "github.com/pavlo67/common/common/logger", callInfo.PackageFullName)
	require.Equal(t, "logger", callInfo.PackageName)
	require.Equal(t, "caller_test.go", callInfo.FileName)
	require.Equal(t, "TestGetCallInfo", callInfo.FuncName)

}
