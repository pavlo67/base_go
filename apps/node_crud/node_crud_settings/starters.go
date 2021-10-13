package node_crud_settings

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/auth/auth_jwt"
	"github.com/pavlo67/common/common/auth/auth_server_http"
	"github.com/pavlo67/common/common/auth/auth_stub"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/control"
	"github.com/pavlo67/common/common/db/db_pg"
	"github.com/pavlo67/common/common/rbac"
	"github.com/pavlo67/common/common/server/server_http/server_http_jschmhr"
	"github.com/pavlo67/common/common/starter"
	"github.com/pkg/errors"

	"github.com/pavlo67/data/entities/persons01"
	"github.com/pavlo67/data/entities/persons01/persons01_pg"
	"github.com/pavlo67/data/entities/records01"
	"github.com/pavlo67/data/entities/records01/records01_pg"

	"github.com/pavlo67/data/components/crud/crud_dispatcher"
	"github.com/pavlo67/data/components/crud/crud_node_http"
	"github.com/pavlo67/data/components/crud/crud_server_http"
)

const onComponents = "on node_crud.Starters()"

func Starters(cfgService config.Config, logRequests bool) ([]starter.Starter, error) {

	var actors []auth.Actor
	if err := cfgService.Value("actors", &actors); err != nil {
		return nil, errors.Wrap(err, onComponents)
	}

	starters := []starter.Starter{
		// general purposes components
		{control.Starter(), nil},
		{db_pg.Starter(), nil},
		{server_http_jschmhr.Starter(), nil},

		// auth components
		{auth_stub.Starter(), common.Map{"interface_key": auth.InterfaceKey}},
		{auth_jwt.Starter(), nil},
		{auth_server_http.Starter(), common.Map{"auth_jwt_key": auth_jwt.InterfaceKey}},

		// CRUD components
		{persons01_pg.Starter(), common.Map{"crud_key": persons01.InterfaceCRUDKey, "roles": rbac.Roles{rbac.RoleAdmin}}},
		{records01_pg.Starter(), common.Map{"crud_key": records01.InterfaceCRUDKey, "roles": rbac.Roles{rbac.RoleAdmin}}},
		{crud_dispatcher.Starter(), nil},
		{crud_server_http.Starter(), nil},

		// app starter
		{crud_node_http.Starter(), nil},
	}

	return starters, nil
}
