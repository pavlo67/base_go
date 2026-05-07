package files_sqlite

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/config"

	"github.com/pavlo67/data/entities/files"
)

func TestFilesSQLite(t *testing.T) {
	os.Setenv("ENV", "TEST")
	os.Setenv("SHOW_CONNECTS", "1")

	_, l := config.PrepareTests(
		t,
		"../../../_environments/",
		"test",
		"", // "persons_test."+strconv.FormatInt(time.Now().Unix(), 10)+".log",
	)

	dsn := "test.sqlite"

	filesOp, filesCleaner, err := New(dsn, l)
	require.NotNil(t, filesOp)
	require.NotNil(t, filesCleaner)
	require.NoError(t, err)

	files.FilesTestScenario(t, filesOp, filesCleaner)
}
