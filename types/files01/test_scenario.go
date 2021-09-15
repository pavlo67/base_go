package files01

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/joiner"
)

const path1 = "bbb/ccc"

var fileData1 = []byte("fileData1")

const path2 = "aaa"

var fileData2 = []byte("fileData2")

func FilesTestScenario(t *testing.T, joinerOp joiner.Operator, interfaceKey, interfaceCleanerKey joiner.InterfaceKey) {
	filesOp, _ := joinerOp.Interface(interfaceKey).(Operator)
	require.NotNil(t, filesOp)

	filesCleanerOp, _ := joinerOp.Interface(interfaceCleanerKey).(db.Cleaner)
	require.NotNil(t, filesCleanerOp)
	err := filesCleanerOp.Clean()
	require.NoError(t, err)

	path1Saved := saveTest(t, filesOp, path1, fileData1)
	require.NotEmpty(t, path1Saved)

	path2Saved := saveTest(t, filesOp, path2, fileData2)
	require.NotEmpty(t, path2Saved)
}

const noSuchFileStr = "no such file or directory"

func saveTest(t *testing.T, filesOp Operator, path string, data []byte) (pathCorrected string) {

	// check original path info ---------------------------------------------

	fi, err := filesOp.Stat(filepath.Dir(path), -1)

	var size0 int64
	//if err == nil {
	//	require.NotNil(t, fi)
	//	require.True(t, fi.IsDir)
	//	size0 = fi.Size
	//} else if errStr := strings.TrimSpace(err.Error()); len(errStr) >= len(noSuchFileStr) && errStr[len(errStr)-len(noSuchFileStr):] == noSuchFileStr {
	//	// no such file or directory: ok
	//} else {
	//	require.FailNow(t, "unexpected error", err)
	//}

	if err != nil {
		require.Nil(t, fi)
		require.True(t, os.IsNotExist(errors.Cause(err)))
	} else {
		require.NotNil(t, fi)
		size0 = fi.Size
	}

	// save file ------------------------------------------------------------

	pathSaved, err := filesOp.Save(path, "", data)
	require.NoError(t, err)
	require.NotEmpty(t, pathSaved)

	// check .Read(), .Items(), .Stat() --------------------------------------

	dataReaded, err := filesOp.Read(pathSaved)
	require.NoError(t, err)
	require.Equal(t, data, dataReaded)

	fis, err := filesOp.List(filepath.Dir(pathSaved), 0)
	require.NoError(t, err)

	// require.FailNowf(t, "%s --> %#v", filepath.Dir(pathSaved), fis)

	found := false
	for _, fi := range fis {
		if filepath.Base(fi.Path) == filepath.Base(pathSaved) {
			found = true
			require.Equalf(t, len(data), int(fi.Size), "%#v", fi)
		}
	}
	require.Truef(t, found, "%s / %#v", pathSaved, fis)

	fi, err = filesOp.Stat(filepath.Dir(pathSaved), -1)
	require.NoError(t, err)
	require.NotNil(t, fi)
	require.True(t, fi.IsDir)
	require.Equalf(t, size0+int64(len(data)), fi.Size, "%#v", fi)

	// remove file ----------------------------------------------------------

	err = filesOp.Remove(pathSaved)
	require.NoError(t, err)

	// check .Read(), .Items(), .Stat() --------------------------------------

	dataReaded, err = filesOp.Read(pathSaved)
	require.Error(t, err)
	require.Nil(t, dataReaded)

	fis, err = filesOp.List(filepath.Dir(pathSaved), 0)
	require.NoError(t, err)

	found = false
	for _, fi := range fis {
		if filepath.Base(fi.Path) == filepath.Base(pathSaved) {
			found = true
			require.FailNowf(t, "this file should be removed", "%#v", fi)
		}
	}

	fi, err = filesOp.Stat(filepath.Dir(pathSaved), -1)
	require.NoError(t, err)
	require.NotNil(t, fi)
	require.True(t, fi.IsDir)
	require.Equalf(t, size0, fi.Size, "%#v", fi)

	return pathSaved
}
