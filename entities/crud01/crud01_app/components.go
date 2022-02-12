package crud01_app

import (
	"github.com/pavlo67/data/components/crud"
	"github.com/pavlo67/data/entities/crud01/crud01_dispatcher"
	"github.com/pavlo67/data/entities/crud01/crud01_server_http"
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/auth/auth_jwt"
	"github.com/pavlo67/common/common/auth/auth_server_http"
	"github.com/pavlo67/common/common/auth/auth_stub"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/control"
	"github.com/pavlo67/common/common/db/db_pg"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/rbac"
	"github.com/pavlo67/common/common/server/server_http/server_http_jschmhr"
	"github.com/pavlo67/common/common/starter"

	"github.com/pavlo67/data/entities/persons01/persons01_pg"
	"github.com/pavlo67/data/entities/records01/records01_pg"
)

const onComponents = "on node_crud.Components()"

const dbPgTestInterfaceKey joiner.InterfaceKey = "db_pg_test"

func Components(cfgService, cfgTests config.Config, logRequests bool) ([]starter.Starter, error) {

	var actors []auth.Actor
	if err := cfgService.Value("actors", &actors); err != nil {
		return nil, errors.Wrap(err, onComponents)
	}

	components := []starter.Starter{
		// general purposes components
		{control.Starter(), nil, nil},
		{db_pg.Starter(), nil, nil},
		{db_pg.Starter(), common.Map{"interface_key": dbPgTestInterfaceKey}, &cfgTests},
		{server_http_jschmhr.Starter(), nil, nil},

		// auth components
		{auth_stub.Starter(), common.Map{"interface_key": auth.InterfaceKey}, nil},
		{auth_jwt.Starter(), nil, nil},
		{auth_server_http.Starter(), common.Map{"auth_jwt_key": auth_jwt.InterfaceKey}, nil},

		// CRUD components
		{persons01_pg.Starter(), common.Map{"roles": rbac.Roles{rbac.RoleAdmin}}, nil},
		{persons01_pg.Starter(), common.Map{"roles": rbac.Roles{crud.RoleTester}, "db_key": dbPgTestInterfaceKey}, nil},
		{records01_pg.Starter(), common.Map{"roles": rbac.Roles{rbac.RoleAdmin}}, nil},
		{records01_pg.Starter(), common.Map{"roles": rbac.Roles{crud.RoleTester}, "db_key": dbPgTestInterfaceKey}, nil},
		{crud01_dispatcher.Starter(), nil, nil},
		{crud01_server_http.Starter(), nil, nil},

		// app starter
		{Starter(), nil, nil},
	}

	return components, nil
}
