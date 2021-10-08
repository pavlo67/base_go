package crud_dispatcher

import (
	"testing"

	"github.com/pavlo67/common/common"

	"github.com/pavlo67/data/entities/persons01/persons01_pg"

	"github.com/pavlo67/data/entities/persons01"

	"github.com/pavlo67/data/entities/records01/records01_pg"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/db/db_pg"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/data/entities/records01"

	"github.com/pavlo67/data/components/crud"
)

func TestDispatcherRecordsPgCRUD(t *testing.T) {
	cfgService, l := config.PrepareTests(t, "../../../_environments/", "test", "dispatcher_records01_pg.log")
	require.NotNil(t, cfgService)

	components := []starter.Starter{
		{db_pg.Starter(), nil},
		{records01_pg.Starter(), common.Map{"crud_key": records01.InterfaceCRUDKey}},
		{Starter(), nil},
	}

	joinerOp, err := starter.Run(components, &cfgService, "CLI BUILD FOR TEST", l)
	require.NoError(t, err)
	require.NotNil(t, joinerOp)
	defer joinerOp.CloseAll()

	crudOp, _ := joinerOp.Interface(InterfaceKey).(crud.Operator)
	require.NotNil(t, crudOp)

	recordsCleanerOp, _ := joinerOp.Interface(records01.InterfaceCleanerKey).(db.Cleaner)
	require.NotNil(t, recordsCleanerOp)

	crudData := crud.Data{
		Key: crud.Key{
			Type: records01.CRUD01,
			ID:   records01.TestItem.ID,
		},
		Description: records01.TestItem.Description,
		Value:       records01.TestItem.Record01,
	}

	crud.OperatorTestScenario(t, crudOp, recordsCleanerOp, crudData, records01.ReadValueRaw, records01.ChangeItemForTest)
}

func TestDispatcherPersonsPgCRUD(t *testing.T) {
	cfgService, l := config.PrepareTests(t, "../../../_environments/", "test", "persons01_pg.log")
	require.NotNil(t, cfgService)

	components := []starter.Starter{
		{db_pg.Starter(), nil},
		{persons01_pg.Starter(), common.Map{"crud_key": persons01.InterfaceCRUDKey}},
		{Starter(), nil},
	}

	joinerOp, err := starter.Run(components, &cfgService, "CLI BUILD FOR TEST", l)
	require.NoError(t, err)
	require.NotNil(t, joinerOp)
	defer joinerOp.CloseAll()

	crudOp, _ := joinerOp.Interface(InterfaceKey).(crud.Operator)
	require.NotNil(t, crudOp)

	personsCleanerOp, _ := joinerOp.Interface(persons01.InterfaceCleanerKey).(db.Cleaner)
	require.NotNil(t, personsCleanerOp)

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
