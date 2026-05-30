package files_sqlite

import (
	"os"
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
