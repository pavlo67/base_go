package persons01_stub

import (
	"testing"

	"github.com/pavlo67/data/components/crud"

	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/starter"
	"github.com/pavlo67/data/entities/persons01"
	"github.com/stretchr/testify/require"
)

//func TestPersonsStub(t *testing.T) {
//	cfgService, l := config.PrepareTests(t, "../../../_environments/", "test", "persons_stub.log")
//	require.NotNil(t, cfgService)
//
//	//var cfg config.Access
//	//err := cfgService.Value("files_fs", &cfg)
//	//require.NoErrorf(t, err, "%#v", cfgService)
//
//	components := []starter.Starter{
//		{Starter(), nil},
//	}
//
//	joinerOp, err := starter.Run(components, &cfgService, "CLI BUILD FOR TEST", l)
//	require.NoError(t, err)
//	require.NotNil(t, joinerOp)
//	defer joinerOp.CloseAll()
//
//	persons01.OperatorTestScenario(t, joinerOp, persons01.InterfaceKey, persons01.InterfaceCleanerKey, persons01.TestItem)
//}

func TestPersonsStubCRUD(t *testing.T) {
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

	personsOp, _ := joinerOp.Interface(persons01.InterfaceKey).(persons01.Operator)
	require.NotNil(t, personsOp)

	personsCleanerOp, _ := joinerOp.Interface(persons01.InterfaceCleanerKey).(db.Cleaner)
	require.NotNil(t, personsCleanerOp)

	crudOp, err := persons01.OperatorCRUD(personsOp)
	require.NoError(t, err)
	require.NotNil(t, crudOp)

	crudData := crud.Data{
		Key: crud.Key{
			Type: persons01.CRUD01,
			ID:   persons01.TestItem.ID,
		},
		Description: persons01.TestItem.Description,
		Value:       persons01.TestItem.Person01,
	}

	crud.OperatorTestScenario(t, crudOp, personsCleanerOp, crudData, persons01.ReadValueRaw, persons01.ChangeItemForTest)
}
