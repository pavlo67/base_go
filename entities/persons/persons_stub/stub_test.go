package persons_stub

import (
	"testing"

	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/starter"
	"github.com/pavlo67/data/entities/persons"
	"github.com/stretchr/testify/require"
)

func TestPersonsStub(t *testing.T) {
	cfgService, l := config.PrepareTests(t, "../../../_environments/", "test", "persons_stub.log")
	require.NotNil(t, cfgService)

	//var cfg config.Access
	//err := cfgService.Value("files_fs", &cfg)
	//require.NoErrorf(t, err, "%#v", cfgService)

	components := []starter.Starter{
		{Starter(), nil},
	}

	joinerOp, err := starter.Run(components, &cfgService, "CLI BUILD FOR TEST", l)
	require.NoError(t, err)
	require.NotNil(t, joinerOp)
	defer joinerOp.CloseAll()

	persons.OperatorTestScenario(t, joinerOp)
}
