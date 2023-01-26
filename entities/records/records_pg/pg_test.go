package records_pg

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/db/db_pg"
	"github.com/pavlo67/common/common/rbac"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/data/components/crud"

	"github.com/pavlo67/data/entities/records"
)

func TestRecordsPgCRUD(t *testing.T) {
	cfgService, l := config.PrepareTests(t, "../../../_environments/", "test", "records_pg.log")
	require.NotNil(t, cfgService)

	//var cfg config.Access
	//err := cfgService.Value("files_fs", &cfg)
	//require.NoErrorf(t, err, "%#v", cfgService)

	components := []starter.Starter{
		{db_pg.Starter(), nil, nil},
		{Starter(), common.Map{"roles": rbac.Roles{crud.RoleTester}, "domain": "test"}, nil},
	}

	joinerOp, err := starter.Run(components, &cfgService, "CLI BUILD FOR TEST", l)
	require.NoError(t, err)
	require.NotNil(t, joinerOp)
	defer joinerOp.CloseAll()

	recordsOp, _ := joinerOp.Interface(records.InterfaceTestKey).(records.Operator)
	require.NotNil(t, recordsOp)

	recordsCleanerOp, _ := joinerOp.Interface(records.InterfaceCleanerKey).(db.Cleaner)
	require.NotNil(t, recordsCleanerOp)

	crudOp, err := records.OperatorCRUD(recordsOp, rbac.Roles{crud.RoleTester})
	require.NoError(t, err)
	require.NotNil(t, crudOp)

	crudData := crud.Data{
		Key: crud.Key{
			Type: records.CRUD,
			ID:   records.TestItem.ID,
		},
		Description: records.TestItem.Description,
		Value:       records.TestItem.Record,
	}

	testActor := auth.Actor{Identity: auth.IdentityWithRoles(crud.RoleTester)}

	crud.OperatorTestScenario(t, crudOp, recordsCleanerOp, crudData, records.ReadValueRaw, records.ChangeCRUDItemForTest, testActor)
}
