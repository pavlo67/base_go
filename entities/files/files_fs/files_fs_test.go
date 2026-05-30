package files_fs

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/base_go/entities/files"
	"github.com/pavlo67/base_go/lib/logger"
	"github.com/pavlo67/base_go/lib/logger/logger_zap"
)

func TestFilesFS(t *testing.T) {
	require.NoError(t, os.Setenv("ENV", "TEST"))

	cfg := logger.Config{
		Key:      strings.ReplaceAll(time.Now().Format(time.RFC3339)[:19], ":", "_"),
		LogLevel: logger.TraceLevel,
	}

	l, err := logger_zap.New(cfg)
	require.NoError(t, err)

	root := t.TempDir()
	filesOp, filesCleaner, err := New(root, l)
	require.NoError(t, err)
	require.NotNil(t, filesOp)
	require.NotNil(t, filesCleaner)

	dir := filepath.Join(root, "data")
	files.FilesTestScenario(t, filesOp, dir, filesCleaner)
}

func TestFilesFSSpecificFields(t *testing.T) {
	require.NoError(t, os.Setenv("ENV", "TEST"))

	cfg := logger.Config{
		Key:      strings.ReplaceAll(time.Now().Format(time.RFC3339)[:19], ":", "_"),
		LogLevel: logger.TraceLevel,
	}

	l, err := logger_zap.New(cfg)
	require.NoError(t, err)

	root := t.TempDir()
	filesOp, filesCleaner, err := New(root, l)
	require.NoError(t, err)
	require.NotNil(t, filesOp)
	require.NotNil(t, filesCleaner)

	dir := root
	require.NoError(t, filesCleaner.Clean(dir))

	mtime := time.Now().Add(-time.Hour).UTC().Truncate(time.Second)
	data := files.Data{
		Path:     filepath.Join(dir, "file.txt"),
		Size:     10,
		MTime:    mtime,
		MimeType: "text/plain; charset=utf-8",
	}

	require.NoError(t, filesOp.Save(data))

	item := requireReadFSData(t, filesOp, data)
	require.True(t, item.CTime.IsZero())
	require.True(t, item.CreatedAt.IsZero())
	require.Equal(t, item.MTime, item.UpdatedAt)

	items, err := filesOp.List(dir, 1)
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, item.Data, items[0].Data)
	require.True(t, items[0].CreatedAt.IsZero())
	require.Equal(t, items[0].MTime, items[0].UpdatedAt)
}

func TestFilesFSSymlinkMimeType(t *testing.T) {
	require.NoError(t, os.Setenv("ENV", "TEST"))

	cfg := logger.Config{
		Key:      strings.ReplaceAll(time.Now().Format(time.RFC3339)[:19], ":", "_"),
		LogLevel: logger.TraceLevel,
	}

	l, err := logger_zap.New(cfg)
	require.NoError(t, err)

	root := t.TempDir()
	filesOp, filesCleaner, err := New(root, l)
	require.NoError(t, err)
	require.NotNil(t, filesOp)
	require.NotNil(t, filesCleaner)

	dir := root
	require.NoError(t, filesCleaner.Clean(dir))

	data := files.Data{
		Path: filepath.Join(dir, "file.txt"),
		Size: 10,
	}
	require.NoError(t, filesOp.Save(data))
	require.NoError(t, os.Symlink("file.txt", filepath.Join(dir, "link")))

	link, err := filesOp.Read(filepath.Join(dir, "link"))
	require.NoError(t, err)
	require.NotNil(t, link)
	require.Equal(t, MimeTypeSymlink, link.MimeType)
}

func TestFilesFSRootBoundary(t *testing.T) {
	require.NoError(t, os.Setenv("ENV", "TEST"))

	cfg := logger.Config{
		Key:      strings.ReplaceAll(time.Now().Format(time.RFC3339)[:19], ":", "_"),
		LogLevel: logger.TraceLevel,
	}

	l, err := logger_zap.New(cfg)
	require.NoError(t, err)

	parent := t.TempDir()
	root := filepath.Join(parent, "root")
	outside := filepath.Join(parent, "outside.txt")

	filesOp, filesCleaner, err := New(root, l)
	require.NoError(t, err)
	require.NotNil(t, filesOp)
	require.NotNil(t, filesCleaner)

	require.NoError(t, os.WriteFile(outside, []byte("outside"), 0o644))
	require.NoError(t, filesOp.Save(files.Data{Path: filepath.Join(root, "inside.txt"), Size: 10}))

	require.NoError(t, filesCleaner.Clean(root))

	rootInfo, err := os.Stat(root)
	require.NoError(t, err)
	require.True(t, rootInfo.IsDir())

	inside, err := filesOp.Read(filepath.Join(root, "inside.txt"))
	require.NoError(t, err)
	require.Nil(t, inside)

	outsideData, err := os.ReadFile(outside)
	require.NoError(t, err)
	require.Equal(t, []byte("outside"), outsideData)

	require.Error(t, filesCleaner.Clean(outside))
	require.Error(t, filesOp.Save(files.Data{Path: outside, Size: 1}))
	require.Error(t, filesOp.Remove(root, true))
}

func requireReadFSData(t *testing.T, op files.Operator, expected files.Data) *files.Item {
	t.Helper()

	actual, err := op.Read(expected.Path)
	require.NoError(t, err)
	require.NotNil(t, actual)

	require.Equal(t, filepath.ToSlash(expected.Path), actual.Path)
	require.Equal(t, expected.IsDir, actual.IsDir)
	require.Equal(t, expected.Size, actual.Size)
	require.Equal(t, expected.MimeType, actual.MimeType)
	require.True(t, actual.CreatedAt.IsZero())
	require.False(t, actual.UpdatedAt.IsZero())

	return actual
}
