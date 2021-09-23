package files01_fs

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/starter"
	"github.com/pavlo67/data/entities/files"
)

func TestFilesFS(t *testing.T) {
	cfgService, l := config.PrepareTests(t, "../../../_environments/", "test", "files_fs.log")
	require.NotNil(t, cfgService)

	var cfg config.Access
	err := cfgService.Value("files_fs", &cfg)
	require.NoErrorf(t, err, "%#v", cfgService)

	components := []starter.Starter{
		{Starter(), common.Map{"base_path": cfg.Path}},
	}

	joinerOp, err := starter.Run(components, &cfgService, "CLI BUILD FOR TEST", l)
	require.NoError(t, err)
	require.NotNil(t, joinerOp)
	defer joinerOp.CloseAll()

	files.FilesTestScenario(t, joinerOp, files.InterfaceKey, files.InterfaceKeyCleaner)
}
