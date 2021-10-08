package crud_server_http

import (
	"fmt"

	"github.com/pavlo67/data/components/crud"
	"github.com/pavlo67/data/components/crud/crud_dispatcher"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"
)

func Starter() starter.Operator {
	return &crudServerHTTPStarter{}
}

var _ starter.Operator = &crudServerHTTPStarter{}

type crudServerHTTPStarter struct{}

// --------------------------------------------------------------------------

var l logger.Operator
var crudDispatcherOp crud.Operator

func (ds *crudServerHTTPStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ds *crudServerHTTPStarter) Prepare(_ *config.Config, _ common.Map) error {
	return nil
}

func (ds *crudServerHTTPStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	if crudDispatcherOp, _ = joinerOp.Interface(crud_dispatcher.InterfaceKey).(crud.Operator); crudDispatcherOp == nil {
		return fmt.Errorf("no crud.Operator with key %s", crud_dispatcher.InterfaceKey)
	}

	return Endpoints.Join(joinerOp)
}
