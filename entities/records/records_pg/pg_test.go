package records_pg

import (
	"testing"

	"github.com/pavlo67/data/entities"

	records2 "github.com/pavlo67/data/entities/records/records_scenarios"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/db/db_pg"
	"github.com/pavlo67/common/common/rbac"
	"github.com/pavlo67/common/common/starter"

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
		{Starter(), common.Map{"roles": rbac.Roles{entities.RoleTester}, "domain": "test"}, nil},
	}

	joinerOp, err := starter.Run(components, &cfgService, "CLI BUILD FOR TEST", l)
	require.NoError(t, err)
	require.NotNil(t, joinerOp)
	defer joinerOp.CloseAll()

	recordsOp, _ := joinerOp.Interface(records.InterfaceTestKey).(records.Operator)
	require.NotNil(t, recordsOp)

	recordsCleanerOp, _ := joinerOp.Interface(records.InterfaceCleanerKey).(db.Cleaner)
	require.NotNil(t, recordsCleanerOp)

	crudOp, err := records.OperatorCRUD(recordsOp, rbac.Roles{entities.RoleTester})
	require.NoError(t, err)
	require.NotNil(t, crudOp)

	crudData := entities.Data{
		Key: entities.Key{
			Type: records.CRUD,
			ID:   records2.TestItem.ID,
		},
		Description: records2.TestItem.Description,
		Value:       records2.TestItem.Record,
	}

	testActor := auth.Actor{Identity: auth.IdentityWithRoles(entities.RoleTester)}

	entities.OperatorTestScenario(t, crudOp, recordsCleanerOp, crudData, records2.ReadValueRaw, records2.ChangeCRUDItemForTest, testActor)
}
