package crud01_dispatcher

import (
	"fmt"

	"github.com/pavlo67/common/common/rbac"
	"github.com/pavlo67/data/components/crud"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"
)

const InterfaceKey joiner.InterfaceKey = "crud_dispatcher"

func Starter() starter.Operator {
	return &crudDispatcherStarter{}
}

var l logger.Operator
var _ starter.Operator = &crudDispatcherStarter{}

type crudDispatcherStarter struct {
	interfaceMessagesKey joiner.InterfaceKey
	interfaceStorageKey  joiner.InterfaceKey

	interfaceKey joiner.InterfaceKey
}

func (sds *crudDispatcherStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (sds *crudDispatcherStarter) Prepare(cfg *config.Config, options common.Map) error {
	sds.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	return nil
}

func (sds *crudDispatcherStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	crudOps := CrudOps{}
	crudComponents := joinerOp.InterfacesAll((*crud.Operator)(nil))
	for _, c := range crudComponents {
		if crudOp, _ := c.Interface.(crud.Operator); crudOp == nil {
			return fmt.Errorf("wrong crud.Operator in joiner.Component (%#v)", c)
		} else {
			types, err := crudOp.Types()
			if err != nil {
				return fmt.Errorf("in joiner.Component (%#v): %s", c, err)
			}
			if len(types) < 1 {
				return fmt.Errorf("in joiner.Component (%#v): no types to use this crud.Operator", c)
			}

			roles, err := crudOp.Roles()
			if err != nil {
				return fmt.Errorf("in joiner.Component (%#v): %s", c, err)
			}
			if len(roles) < 1 {
				return fmt.Errorf("in joiner.Component (%#v): no roles to use this crud.Operator", c)
			}

			for _, t := range types {
				crudOpsTyped := crudOps[t]
				if crudOpsTyped == nil {
					crudOpsTyped = map[rbac.Role]crud.Operator{}
				}

				for _, r := range roles {
					crudOpsTyped[r] = crudOp
				}
				crudOps[t] = crudOpsTyped
			}
		}
	}

	crudOp, err := New(crudOps)
	if err != nil {
		return errors.Wrap(err, "can't init *crudDispatcher{} as crud.Operator")
	}

	if err = joinerOp.Join(crudOp, sds.interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join *crudDispatcher(%#v) as crud.Operator with key '%s'", crudOp, sds.interfaceKey)
	}

	return nil
}
