package records01_pg

import (
	"testing"

	"github.com/pavlo67/common/common"

	"github.com/pavlo67/common/common/auth"

	"github.com/pavlo67/common/common/rbac"

	"github.com/pavlo67/data/components/crud"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/db/db_pg"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/data/entities/records01"
)

func TestRecordsPgCRUD(t *testing.T) {
	cfgService, l := config.PrepareTests(t, "../../../_environments/", "test", "records01_pg.log")
	require.NotNil(t, cfgService)

	//var cfg config.Access
	//err := cfgService.Value("files_fs", &cfg)
	//require.NoErrorf(t, err, "%#v", cfgService)

	components := []starter.Starter{
		{db_pg.Starter(), nil, nil},
		{Starter(), common.Map{"roles": rbac.Roles{crud.RoleTester}}, nil},
	}

	joinerOp, err := starter.Run(components, &cfgService, "CLI BUILD FOR TEST", l)
	require.NoError(t, err)
	require.NotNil(t, joinerOp)
	defer joinerOp.CloseAll()

	recordsOp, _ := joinerOp.Interface(records01.InterfaceTestKey).(records01.Operator)
	require.NotNil(t, recordsOp)

	recordsCleanerOp, _ := joinerOp.Interface(records01.InterfaceCleanerKey).(db.Cleaner)
	require.NotNil(t, recordsCleanerOp)

	crudOp, err := records01.OperatorCRUD(recordsOp, rbac.Roles{crud.RoleTester})
	require.NoError(t, err)
	require.NotNil(t, crudOp)

	crudData := crud.Data{
		Key: crud.Key{
			Type: records01.CRUD01,
			ID:   records01.TestItem.ID,
		},
		Description: records01.TestItem.Description,
		Value:       records01.TestItem.Record01,
	}

	testActor := auth.Actor{Identity: auth.IdentityWithRoles(crud.RoleTester)}

	crud.OperatorTestScenario(t, crudOp, recordsCleanerOp, crudData, records01.ReadValueRaw, records01.ChangeItemForTest, testActor)
}
