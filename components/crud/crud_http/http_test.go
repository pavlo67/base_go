package crud_http

import (
	"testing"
	"time"

	"github.com/pavlo67/data/entities/records01"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/data/components/crud/crud_node_http"

	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/db"
	"github.com/pavlo67/common/common/starter"
	"github.com/pavlo67/data/apps/node_crud/node_crud_settings"
	"github.com/pavlo67/data/components/crud"
	"github.com/pavlo67/data/entities/persons01"
	"github.com/stretchr/testify/require"
)

func TestHTTPRecordsCRUD(t *testing.T) {
	cfgService, l := config.PrepareTests(t, "../../../_environments/", "test", "dispatcher_records01_pg.log")
	require.NotNil(t, cfgService)

	components := append(
		node_crud_settings.Components(true),
		starter.Starter{Starter(), common.Map{
			"prefix":        crud_node_http.PrefixREST,
			"server_config": crud_node_http.ServerConfig,
		}},
	)

	joinerOp, err := starter.Run(components, &cfgService, "CLI BUILD FOR TEST", l)
	require.NoError(t, err)
	require.NotNil(t, joinerOp)
	defer joinerOp.CloseAll()

	time.Sleep(time.Second)

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

	crud.OperatorTestScenario(t, crudOp, recordsCleanerOp, crudData, records01.ChangeForTest)
}

func TestHTTPPersonsCRUD(t *testing.T) {
	cfgService, l := config.PrepareTests(t, "../../../_environments/", "test", "persons01_pg.log")
	require.NotNil(t, cfgService)

	components := append(
		node_crud_settings.Components(true),
		starter.Starter{Starter(), common.Map{
			"prefix":        crud_node_http.PrefixREST,
			"server_config": crud_node_http.ServerConfig,
		}},
	)

	joinerOp, err := starter.Run(components, &cfgService, "CLI BUILD FOR TEST", l)
	require.NoError(t, err)
	require.NotNil(t, joinerOp)
	defer joinerOp.CloseAll()

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

	crud.OperatorTestScenario(t, crudOp, personsCleanerOp, crudData, persons01.ChangeTestCRUDItem)
}

func TestCRUDHTTP(t *testing.T) {

	//var cfgService config.Config
	//cfgService, l = config.PrepareTests(
	//	t,
	//	"../../../_environments/",
	//	"test",
	//	"", // "connect_test."+strconv.FormatInt(time.Now().Unix(), 10)+".log",
	//)
	//
	//var cfgServerHTTP config.Access
	//err := cfgService.Value("server_http", &cfgServerHTTP)
	//require.NoError(t, err)
	//
	//serverConfig := demo_settings.ServerConfig
	//
	//err = serverConfig.CompleteDirectly(crud_server_http.Endpoints, cfgServerHTTP.Host, cfgServerHTTP.Port, demo_settings.PrefixREST)
	//require.NoError(t, err)
	//
	//crud.OperatorTestScenarioPassword(t, crudOp)

}
