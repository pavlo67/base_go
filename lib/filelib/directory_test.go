package filelib

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFullPath(t *testing.T) {
	parent := filepath.Join(string(filepath.Separator), "tmp", "base_go_filelib_test")
	root := filepath.Join(parent, "root")

	fullPath, err := FullPath(root, "relative/file.txt")
	require.NoError(t, err)
	require.Equal(t, filepath.Join(root, "relative", "file.txt"), fullPath)

	fullPath, err = FullPath(root, root)
	require.NoError(t, err)
	require.Equal(t, root, fullPath)

	inside := filepath.Join(root, "inside", "file.txt")
	fullPath, err = FullPath(root, inside)
	require.NoError(t, err)
	require.Equal(t, inside, fullPath)

	_, err = FullPath(root, filepath.Join(parent, "outside.txt"))
	require.Error(t, err)
}
