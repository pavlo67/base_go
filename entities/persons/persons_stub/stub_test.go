package persons_stub

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/rbac"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/data/entities/persons"

	"github.com/pavlo67/data/components/crud"
)

func TestPersonsStubCRUD(t *testing.T) {
	cfgService, l := config.PrepareTests(t, "../../../_environments/", "test", "persons_stub.log")
	require.NotNil(t, cfgService)

	//var cfg config.Access
	//err := cfgService.Value("files_fs", &cfg)
	//require.NoErrorf(t, err, "%#v", cfgService)

	components := []starter.Starter{
		{Starter(), nil, nil},
	}

	joinerOp, err := starter.Run(components, &cfgService, "CLI BUILD FOR TEST", l)
	require.NoError(t, err)
	require.NotNil(t, joinerOp)
	defer joinerOp.CloseAll()

	personsOp, _ := joinerOp.Interface(persons.InterfaceKey).(persons.Operator)
	require.NotNil(t, personsOp)

	personsCleanerOp, _ := joinerOp.Interface(persons.InterfaceCleanerKey).(db.Cleaner)
	require.NotNil(t, personsCleanerOp)

	crudOp, err := persons.OperatorCRUD(personsOp, rbac.Roles{crud.RoleTester})
	require.NoError(t, err)
	require.NotNil(t, crudOp)

	crudData := crud.Data{
		Key: crud.Key{
			Type: persons.CRUD,
			ID:   persons.TestItem.ID,
		},
		Description: persons.TestItem.Description,
		Value:       persons.TestItem.Person,
	}

	testActor := auth.Actor{Identity: auth.IdentityWithRoles(crud.RoleTester)}

	crud.OperatorTestScenario(t, crudOp, personsCleanerOp, crudData, persons.ReadValueRaw, persons.ChangeCRUDItemForTest, testActor)
}
