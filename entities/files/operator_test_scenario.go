package files

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/base_go/lib/db"
)

func FilesTestScenario(t *testing.T, filesOp Operator, dir string, filesCleaner db.Operator) {
	t.Helper()

	// data preparation -------------------------------------------------

	now := time.Now()

	data1 := Data{
		Path:     filepath.Join(dir, "file1.txt"),
		IsDir:    false,
		Size:     100,
		CTime:    now.UTC(),
		MTime:    now.Add(time.Minute).UTC(),
		MimeType: "text/plain",
	}

	data2 := Data{
		Path:     filepath.Join(dir, "file2.json"),
		IsDir:    false,
		Size:     200,
		CTime:    now.Add(time.Hour).UTC(),
		MTime:    now.Add(time.Hour + time.Minute).UTC(),
		MimeType: "application/json",
	}

	data3 := Data{
		Path:     filepath.Join(dir, "file3.bin"),
		IsDir:    false,
		Size:     300,
		CTime:    now.Add(2 * time.Hour).UTC(),
		MTime:    now.Add(2*time.Hour + time.Minute).UTC(),
		MimeType: "application/octet-stream",
	}

	dataSub := Data{
		Path:     filepath.Join(dir, "sub", "file4.txt"),
		IsDir:    false,
		Size:     400,
		CTime:    now.Add(3 * time.Hour).UTC(),
		MTime:    now.Add(3*time.Hour + time.Minute).UTC(),
		MimeType: "text/plain",
	}

	dataSibling := Data{
		Path:     filepath.Join(filepath.Dir(dir), filepath.Base(dir)+"-sibling", "file5.txt"),
		IsDir:    false,
		Size:     500,
		CTime:    now.Add(4 * time.Hour).UTC(),
		MTime:    now.Add(4*time.Hour + time.Minute).UTC(),
		MimeType: "text/plain",
	}

	// ------------------------------------------------------------------

	err := filesCleaner.Clean("")
	require.Error(t, err)

	err = filesCleaner.Clean(dir)
	require.NoError(t, err)

	// + data1 ----------------------------------------------------------

	items, err := filesOp.List(dir, 0)
	require.NoError(t, err)
	require.Empty(t, items)

	err = filesOp.Save(data1)
	require.NoError(t, err)

	item1 := requireReadData(t, filesOp, data1)
	require.Equal(t, item1.CreatedAt, item1.UpdatedAt)

	requireListData(t, filesOp, dir, data1)

	err = filesOp.Save(data1)
	require.NoError(t, err)

	item1AfterResave := requireReadData(t, filesOp, data1)
	require.Equal(t, item1.CreatedAt, item1AfterResave.CreatedAt)
	require.True(t, !item1AfterResave.UpdatedAt.Before(item1.UpdatedAt))

	requireListData(t, filesOp, dir, data1)

	// + data2 ----------------------------------------------------------

	err = filesOp.Save(data2)
	require.NoError(t, err)

	item1AfterData2 := requireReadData(t, filesOp, data1)
	require.Equal(t, *item1AfterResave, *item1AfterData2)

	requireReadData(t, filesOp, data2)
	requireListData(t, filesOp, dir, data1, data2)

	// + data3 ----------------------------------------------------------

	err = filesOp.Save(data3)
	require.NoError(t, err)

	requireReadData(t, filesOp, data1)
	requireReadData(t, filesOp, data2)
	requireReadData(t, filesOp, data3)
	requireListData(t, filesOp, dir, data1, data2, data3)

	item1BeforeRemove := requireReadData(t, filesOp, data1)
	item3BeforeRemove := requireReadData(t, filesOp, data3)

	// - data2 ----------------------------------------------------------

	err = filesOp.Remove(data2.Path)
	require.NoError(t, err)

	requireListData(t, filesOp, dir, data1, data3)

	item1AfterRemove := requireReadData(t, filesOp, data1)
	item3AfterRemove := requireReadData(t, filesOp, data3)

	require.Equal(t, *item1BeforeRemove, *item1AfterRemove)
	require.Equal(t, *item3BeforeRemove, *item3AfterRemove)

	removed, err := filesOp.Read(data2.Path)
	require.NoError(t, err)
	require.Nil(t, removed)

	// clean directory tree ---------------------------------------------

	err = filesOp.Save(dataSub)
	require.NoError(t, err)

	err = filesOp.Save(dataSibling)
	require.NoError(t, err)

	requireReadData(t, filesOp, dataSub)
	requireReadData(t, filesOp, dataSibling)

	err = filesCleaner.Clean(dir)
	require.NoError(t, err)

	for _, path := range []string{data1.Path, data3.Path, dataSub.Path} {
		item, err := filesOp.Read(path)
		require.NoError(t, err)
		require.Nil(t, item)
	}

	requireReadData(t, filesOp, dataSibling)
}

func requireReadData(t *testing.T, op Operator, expected Data) *Item {
	t.Helper()

	actual, err := op.Read(expected.Path)
	require.NoError(t, err)
	require.NotNil(t, actual)

	require.Equal(t, expected, actual.Data, "file data differs: %s", expected.Path)
	require.False(t, actual.CreatedAt.IsZero(), "CreatedAt must be set: %s", expected.Path)
	require.False(t, actual.UpdatedAt.IsZero(), "UpdatedAt must be set: %s", expected.Path)

	return actual
}

func requireListData(t *testing.T, op Operator, dir string, expected ...Data) {
	t.Helper()

	actual, err := op.List(dir, 0)
	require.NoError(t, err)
	require.Len(t, actual, len(expected))

	actualData := make([]Data, 0, len(actual))
	for _, item := range actual {
		require.False(t, item.CreatedAt.IsZero(), "CreatedAt must be set: %s", item.Path)
		require.False(t, item.UpdatedAt.IsZero(), "UpdatedAt must be set: %s", item.Path)

		actualData = append(actualData, item.Data)
	}

	require.ElementsMatch(t, expected, actualData)
}
