package node_crud_settings

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/auth/auth_jwt"
	"github.com/pavlo67/common/common/auth/auth_server_http"
	"github.com/pavlo67/common/common/auth/auth_stub"
	"github.com/pavlo67/common/common/control"
	"github.com/pavlo67/common/common/server/server_http/server_http_jschmhr"
	"github.com/pavlo67/common/common/starter"
	"github.com/pavlo67/data/components/node_crud/node_crud_server_http"
)

func Components(logRequests bool) []starter.Starter {

	starters := []starter.Starter{
		// general purposes components
		{control.Starter(), nil},

		// auth/persons components
		{auth_stub.Starter(), common.Map{"interface_key": auth.InterfaceKey}},
		{auth_jwt.Starter(), nil},
		{auth_server_http.Starter(), common.Map{"auth_jwt_key": auth_jwt.InterfaceKey}},
	}

	starters = append(
		starters,

		// action managers
		starter.Starter{server_http_jschmhr.Starter(), nil},

		// actions starter (connecting specific actions to the corresponding action managers)
		starter.Starter{node_crud_server_http.Starter(), nil},
	)

	return starters
}
