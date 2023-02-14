package persons_pg

import (
	"testing"

	"github.com/pavlo67/data/entities"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/db/db_pg"
	"github.com/pavlo67/common/common/rbac"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/data/entities/persons"
)

func TestPersonsPgCRUD(t *testing.T) {
	cfgService, l := config.PrepareTests(t, "../../../_environments/", "test", "persons_pg.log")
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

	personsOp, _ := joinerOp.Interface(persons.InterfaceTestKey).(persons.Operator)
	require.NotNil(t, personsOp)

	personsCleanerOp, _ := joinerOp.Interface(persons.InterfaceCleanerKey).(db.Cleaner)
	require.NotNil(t, personsCleanerOp)

	crudOp, err := persons.OperatorCRUD(personsOp, rbac.Roles{entities.RoleTester})
	require.NoError(t, err)
	require.NotNil(t, crudOp)

	crudData := entities.Data{
		Key: entities.Key{
			Type: persons.CRUD,
			ID:   persons.TestItem.ID,
		},
		Description: persons.TestItem.Description,
		Value:       persons.TestItem.Person,
	}

	testActor := auth.Actor{Identity: auth.IdentityWithRoles(entities.RoleTester)}

	entities.OperatorTestScenario(t, crudOp, personsCleanerOp, crudData, persons.ReadValueRaw, persons.ChangeCRUDItemForTest, testActor)
}
