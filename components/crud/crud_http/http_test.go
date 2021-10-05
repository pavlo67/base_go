package crud_http

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/crud/crud_server_http"
)

func TestAuthHTTP(t *testing.T) {

	var cfgService config.Config
	cfgService, l = config.PrepareTests(
		t,
		"../../../_environments/",
		"test",
		"", // "connect_test."+strconv.FormatInt(time.Now().Unix(), 10)+".log",
	)

	var cfgServerHTTP config.Access
	err := cfgService.Value("server_http", &cfgServerHTTP)
	require.NoError(t, err)

	serverConfig := demo_settings.ServerConfig

	err = serverConfig.CompleteDirectly(crud_server_http.Endpoints, cfgServerHTTP.Host, cfgServerHTTP.Port, demo_settings.PrefixREST)
	require.NoError(t, err)

	crudOp, err := New(serverConfig)
	require.NoError(t, err)
	require.NotNil(t, crudOp)

	crud.OperatorTestScenarioPassword(t, crudOp)

}
