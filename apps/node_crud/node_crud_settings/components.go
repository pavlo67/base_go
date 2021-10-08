package node_crud_settings

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/auth/auth_jwt"
	"github.com/pavlo67/common/common/auth/auth_server_http"
	"github.com/pavlo67/common/common/auth/auth_stub"
	"github.com/pavlo67/common/common/control"
	"github.com/pavlo67/common/common/db/db_pg"
	"github.com/pavlo67/common/common/server/server_http/server_http_jschmhr"
	"github.com/pavlo67/common/common/starter"
	"github.com/pavlo67/data/components/crud/crud_dispatcher"
	"github.com/pavlo67/data/components/crud/crud_node_http"
	"github.com/pavlo67/data/components/crud/crud_server_http"
	"github.com/pavlo67/data/entities/persons01"
	"github.com/pavlo67/data/entities/persons01/persons01_pg"
	"github.com/pavlo67/data/entities/records01"
	"github.com/pavlo67/data/entities/records01/records01_pg"
)

func Components(logRequests bool) []starter.Starter {

	starters := []starter.Starter{
		// general purposes components
		{control.Starter(), nil},
		{db_pg.Starter(), nil},

		{server_http_jschmhr.Starter(), nil},

		// auth components
		{auth_stub.Starter(), common.Map{"interface_key": auth.InterfaceKey}},
		{auth_jwt.Starter(), nil},
		{auth_server_http.Starter(), common.Map{"auth_jwt_key": auth_jwt.InterfaceKey}},

		{persons01_pg.Starter(), common.Map{"crud_key": persons01.InterfaceCRUDKey}},
		{records01_pg.Starter(), common.Map{"crud_key": records01.InterfaceCRUDKey}},
		{crud_dispatcher.Starter(), nil},

		// actions starter (connecting specific actions to the corresponding action managers)
		{crud_server_http.Starter(), nil},
		{crud_node_http.Starter(), nil},
	}

	return starters
}
