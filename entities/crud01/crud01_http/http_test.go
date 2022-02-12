package crud01_http

import (
	"testing"
	"time"

	"github.com/pavlo67/data/entities/crud01/crud01_app"

	"github.com/pavlo67/data/components/crud"
	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/auth/auth_http"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/starter"

	auth2 "github.com/pavlo67/data/common/auth"

	"github.com/pavlo67/data/entities/persons01"
	"github.com/pavlo67/data/entities/records01"
)

func TestHTTPRecordsCRUD(t *testing.T) {
	cfgService, l := config.PrepareTests(t, "../../../_environments/", "test", "http_records01_pg.log")
	require.NotNil(t, cfgService)

	starters, err := crud01_app.Components(cfgService, cfgService, true)
	require.NoError(t, err)

	httpOptions := common.Map{"prefix": crud01_app.PrefixREST, "server_config": crud01_app.ServerConfig}

	starters = append(
		starters,
		starter.Starter{auth_http.Starter(), httpOptions, nil},
		starter.Starter{Starter(), httpOptions, nil},
	)

	joinerOp, err := starter.Run(starters, &cfgService, "CLI BUILD FOR TEST", l)
	require.NoError(t, err)
	require.NotNil(t, joinerOp)
	defer joinerOp.CloseAll()

	time.Sleep(time.Second)

	authOp, _ := joinerOp.Interface(auth_http.InterfaceKey).(auth.Operator)
	require.NotNil(t, authOp)

	recordsCleanerOp, _ := joinerOp.Interface(records01.InterfaceCleanerKey).(db.Cleaner)
	require.NotNil(t, recordsCleanerOp)

	crudOp, _ := joinerOp.Interface(InterfaceKey).(crud.Operator)
	require.NotNil(t, crudOp)

	crudData := crud.Data{
		Key: crud.Key{
			Type: records01.CRUD01,
			ID:   records01.TestItem.ID,
		},
		Description: records01.TestItem.Description,
		Value:       records01.TestItem.Record01,
	}

	testActor, err := auth2.Auth(cfgService, authOp, crud.RoleTester)
	require.NoError(t, err)
	require.NotNil(t, testActor)

	crud.OperatorTestScenario(t, crudOp, recordsCleanerOp, crudData, records01.ReadValueRaw, records01.ChangeCRUDItemForTest, *testActor)
}

func TestHTTPPersonsCRUD(t *testing.T) {
	cfgService, l := config.PrepareTests(t, "../../../_environments/", "test", "http_persons01_pg.log")
	require.NotNil(t, cfgService)

	starters, err := crud01_app.Components(cfgService, cfgService, true)
	require.NoError(t, err)

	httpOptions := common.Map{"prefix": crud01_app.PrefixREST, "server_config": crud01_app.ServerConfig}

	starters = append(
		starters,
		starter.Starter{auth_http.Starter(), httpOptions, nil},
		starter.Starter{Starter(), httpOptions, nil},
	)

	joinerOp, err := starter.Run(starters, &cfgService, "CLI BUILD FOR TEST", l)
	require.NoError(t, err)
	require.NotNil(t, joinerOp)
	defer joinerOp.CloseAll()

	time.Sleep(time.Second)

	authOp, _ := joinerOp.Interface(auth_http.InterfaceKey).(auth.Operator)
	require.NotNil(t, authOp)

	testActor, err := auth2.Auth(cfgService, authOp, crud.RoleTester)
	require.NoError(t, err)
	require.NotNil(t, testActor)

	personsCleanerOp, _ := joinerOp.Interface(persons01.InterfaceCleanerKey).(db.Cleaner)
	require.NotNil(t, personsCleanerOp)

	crudOp, _ := joinerOp.Interface(InterfaceKey).(crud.Operator)
	require.NotNil(t, crudOp)

	crudData := crud.Data{
		Key: crud.Key{
			Type: persons01.CRUD01,
			ID:   persons01.TestItem.ID,
		},
		Description: persons01.TestItem.Description,
		Value:       persons01.TestItem.Person01,
	}

	crud.OperatorTestScenario(t, crudOp, personsCleanerOp, crudData, persons01.ReadValueRaw, persons01.ChangeCRUDItemForTest, *testActor)
}
