package logger

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJ(t *testing.T, l OperatorJ) {

	j, outputPaths, err := l.J()
	require.NoError(t, err)
	require.NotNil(t, j)
	require.True(t, len(outputPaths) > 0)
	path := outputPaths[0]

	//l.Infof("outputPaths: %v", outputPaths)

	var opLen int64
	var f os.FileInfo

	// TODO!!! be careful: if log files are removed here the next checks fail

	f, _ = os.Stat(path)
	if f != nil && !f.IsDir() {
		opLen = f.Size()
	}

	l.Comment(`l.Comment()`)

	f, err = os.Stat(path)
	require.NoError(t, err)
	require.NotNil(t, f)
	require.False(t, f.IsDir())
	require.True(t, f.Size() > opLen)
	opLen = f.Size()

	j.Comment(`j.Comment()`)

	f, err = os.Stat(path)
	require.NoError(t, err)
	require.NotNil(t, f)
	require.False(t, f.IsDir())
	require.True(t, f.Size() > opLen)
	opLen = f.Size()

	l.Info("l.Info()")

	f, err = os.Stat(path)
	require.NoError(t, err)
	require.NotNil(t, f)
	require.False(t, f.IsDir())
	require.True(t, f.Size() > opLen)
	opLen = f.Size()

	j.Info("j.Info()")

	f, err = os.Stat(path)
	require.NoError(t, err)
	require.NotNil(t, f)
	require.False(t, f.IsDir())
	require.True(t, f.Size() > opLen)
	opLen = f.Size()

	l.Error("l.Error()")

	f, err = os.Stat(path)
	require.NoError(t, err)
	require.NotNil(t, f)
	require.False(t, f.IsDir())
	require.True(t, f.Size() > opLen)
	opLen = f.Size()

	j.Error("j.Error()")

	f, err = os.Stat(path)
	require.NoError(t, err)
	require.NotNil(t, f)
	require.False(t, f.IsDir())
	require.True(t, f.Size() > opLen)
	opLen = f.Size()

}
