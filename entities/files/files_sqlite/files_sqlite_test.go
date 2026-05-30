package files_sqlite

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

func TestFilesSQLite(t *testing.T) {
	require.NoError(t, os.Setenv("ENV", "TEST"))
	require.NoError(t, os.Setenv("SHOW_CONNECTS", "1"))

	cfg := logger.Config{
		Key:      strings.ReplaceAll(time.Now().Format(time.RFC3339)[:19], ":", "_"),
		LogLevel: logger.TraceLevel,
	}

	l, err := logger_zap.New(cfg)
	require.NoError(t, err)

	dsn := "test.sqlite"

	filesOp, filesCleaner, err := New(dsn, true, l)
	require.NotNil(t, filesOp)
	require.NotNil(t, filesCleaner)
	require.NoError(t, err)

	dir := "/test/data"

	files.FilesTestScenario(t, filesOp, dir, filesCleaner)
}

func TestFilesSQLiteSpecificFields(t *testing.T) {
	require.NoError(t, os.Setenv("ENV", "TEST"))
	require.NoError(t, os.Setenv("SHOW_CONNECTS", "1"))

	cfg := logger.Config{
		Key:      strings.ReplaceAll(time.Now().Format(time.RFC3339)[:19], ":", "_"),
		LogLevel: logger.TraceLevel,
	}

	l, err := logger_zap.New(cfg)
	require.NoError(t, err)

	filesOp, filesCleaner, err := New("test.sqlite", true, l)
	require.NotNil(t, filesOp)
	require.NotNil(t, filesCleaner)
	require.NoError(t, err)

	dir := "/test/specific"
	require.NoError(t, filesCleaner.Clean(dir))

	crc := int64(123456)
	data := files.Data{
		Path:     filepath.Join(dir, "file.txt"),
		Size:     10,
		CTime:    time.Now().Add(-2 * time.Hour).UTC(),
		MTime:    time.Now().Add(-time.Hour).UTC(),
		CRC:      &crc,
		MimeType: "text/plain",
	}

	require.NoError(t, filesOp.Save(data))

	item, err := filesOp.Read(data.Path)
	require.NoError(t, err)
	require.NotNil(t, item)
	require.Equal(t, data, item.Data)
	require.False(t, item.CreatedAt.IsZero())
	require.False(t, item.UpdatedAt.IsZero())
	require.Equal(t, item.CreatedAt, item.UpdatedAt)

	items, err := filesOp.List(dir, 1)
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, data, items[0].Data)
	require.False(t, items[0].CreatedAt.IsZero())
	require.False(t, items[0].UpdatedAt.IsZero())
}
