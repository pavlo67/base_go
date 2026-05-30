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

	file1 := File{
		Path:     filepath.Join(dir, "file1.txt"),
		IsDir:    false,
		Size:     100,
		CTime:    now.UTC(),
		MTime:    now.Add(time.Minute).UTC(),
		MimeType: "text/plain",
	}

	file2 := File{
		Path:     filepath.Join(dir, "file2.json"),
		IsDir:    false,
		Size:     200,
		CTime:    now.Add(time.Hour).UTC(),
		MTime:    now.Add(time.Hour + time.Minute).UTC(),
		MimeType: "application/json",
	}

	file3 := File{
		Path:     filepath.Join(dir, "file3.bin"),
		IsDir:    false,
		Size:     300,
		CTime:    now.Add(2 * time.Hour).UTC(),
		MTime:    now.Add(2*time.Hour + time.Minute).UTC(),
		MimeType: "application/octet-stream",
	}

	// ------------------------------------------------------------------

	err := filesCleaner.Clean(dir)
	require.NoError(t, err)

	// + file1 ----------------------------------------------------------

	items, err := filesOp.List(dir, 0)
	require.NoError(t, err)
	require.Empty(t, items)

	err = filesOp.Save(file1)
	require.NoError(t, err)

	item1 := requireReadFile(t, filesOp, file1)
	require.Equal(t, item1.CreatedAt, item1.UpdatedAt)

	requireListFiles(t, filesOp, dir, file1)

	err = filesOp.Save(file1)
	require.NoError(t, err)

	item1AfterResave := requireReadFile(t, filesOp, file1)
	require.Equal(t, item1.CreatedAt, item1AfterResave.CreatedAt)
	require.True(t, !item1AfterResave.UpdatedAt.Before(item1.UpdatedAt))

	requireListFiles(t, filesOp, dir, file1)

	// + file2 ----------------------------------------------------------

	err = filesOp.Save(file2)
	require.NoError(t, err)

	item1AfterFile2 := requireReadFile(t, filesOp, file1)
	require.Equal(t, *item1AfterResave, *item1AfterFile2)

	requireReadFile(t, filesOp, file2)
	requireListFiles(t, filesOp, dir, file1, file2)

	// + file3 ----------------------------------------------------------

	err = filesOp.Save(file3)
	require.NoError(t, err)

	requireReadFile(t, filesOp, file1)
	requireReadFile(t, filesOp, file2)
	requireReadFile(t, filesOp, file3)
	requireListFiles(t, filesOp, dir, file1, file2, file3)

	item1BeforeRemove := requireReadFile(t, filesOp, file1)
	item3BeforeRemove := requireReadFile(t, filesOp, file3)

	// - file2 ----------------------------------------------------------

	err = filesOp.Remove(file2.Path)
	require.NoError(t, err)

	requireListFiles(t, filesOp, dir, file1, file3)

	item1AfterRemove := requireReadFile(t, filesOp, file1)
	item3AfterRemove := requireReadFile(t, filesOp, file3)

	require.Equal(t, *item1BeforeRemove, *item1AfterRemove)
	require.Equal(t, *item3BeforeRemove, *item3AfterRemove)

	removed, err := filesOp.Read(file2.Path)
	require.NoError(t, err)
	require.Nil(t, removed)
}

func requireReadFile(t *testing.T, op Operator, expected File) *Item {
	t.Helper()

	actual, err := op.Read(expected.Path)
	require.NoError(t, err)
	require.NotNil(t, actual)

	require.Equal(t, expected, actual.File, "file data differs: %s", expected.Path)
	require.False(t, actual.CreatedAt.IsZero(), "CreatedAt must be set: %s", expected.Path)
	require.False(t, actual.UpdatedAt.IsZero(), "UpdatedAt must be set: %s", expected.Path)

	return actual
}

func requireListFiles(t *testing.T, op Operator, dir string, expected ...File) {
	t.Helper()

	actual, err := op.List(dir, 0)
	require.NoError(t, err)
	require.Len(t, actual, len(expected))

	actualFiles := make([]File, 0, len(actual))
	for _, item := range actual {
		require.False(t, item.CreatedAt.IsZero(), "CreatedAt must be set: %s", item.Path)
		require.False(t, item.UpdatedAt.IsZero(), "UpdatedAt must be set: %s", item.Path)

		actualFiles = append(actualFiles, item.File)
	}

	require.ElementsMatch(t, expected, actualFiles)
}
