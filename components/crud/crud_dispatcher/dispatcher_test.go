package crud_dispatcher

import (
	"testing"

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
		{records01_pg.Starter(), nil},
		{Starter(), nil},
	}

	joinerOp, err := starter.Run(components, &cfgService, "CLI BUILD FOR TEST", l)
	require.NoError(t, err)
	require.NotNil(t, joinerOp)
	defer joinerOp.CloseAll()

	recordsOp, _ := joinerOp.Interface(records01.InterfaceKey).(records01.Operator)
	require.NotNil(t, recordsOp)

	recordsCleanerOp, _ := joinerOp.Interface(records01.InterfaceCleanerKey).(db.Cleaner)
	require.NotNil(t, recordsCleanerOp)

	crudOp, err := records01.OperatorCRUD(recordsOp)
	require.NoError(t, err)
	require.NotNil(t, crudOp)

	crud.OperatorTestScenario(t, crudOp, recordsCleanerOp, records01.CRUD01, records01.TestRecord, records01.ChangeForTest)
}

func TestDispatcherPersonsPgCRUD(t *testing.T) {
	cfgService, l := config.PrepareTests(t, "../../../_environments/", "test", "persons01_pg.log")
	require.NotNil(t, cfgService)

	components := []starter.Starter{
		{db_pg.Starter(), nil},
		{persons01_pg.Starter(), nil},
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

	crud.OperatorTestScenario(t, crudOp, personsCleanerOp, persons01.CRUD01, persons01.TestPersonToSave, persons01.ChangeTestCRUDItem)
}
