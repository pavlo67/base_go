package crud_dispatcher

import (
	"fmt"

	"github.com/pavlo67/common/common/rbac"

	"github.com/pavlo67/common/common/auth"

	"github.com/pavlo67/data/elements/selectors"

	"github.com/pavlo67/data/components/crud"
)

var _ crud.Operator = &crudDispatcher{}

type CrudOps map[crud.Type]map[rbac.Role]crud.Operator

type crudDispatcher struct {
	dispatchedOps CrudOps
}

const onNew = "on crudDispatcher.New()"

func New(dispatchedOps CrudOps) (crud.Operator, error) {
	crudOp := crudDispatcher{dispatchedOps: dispatchedOps}

	return &crudOp, nil

}

func (crudOp crudDispatcher) dispatchedOp(crudType crud.Type, actor auth.Actor) (crud.Operator, error) {
	if actor.Identity == nil {
		return nil, fmt.Errorf("no actor.Identity to get crud.Operator for actor (%#v)", actor)
	}

	crudOpsTyped := crudOp.dispatchedOps[crudType]
	for _, role := range actor.Identity.Roles {
		if crudOp := crudOpsTyped[role]; crudOp != nil {
			return crudOp, nil
		}
	}

	return nil, fmt.Errorf("no crud.Operator for type/actor (%s / %#v)", crudType, actor)
}

// operator ----------------------------------------------------------------------------------------------------------------

var _ crud.Operator = &crudDispatcher{}

func (crudOp crudDispatcher) Save(data crud.Data, actor auth.Actor) (*crud.Key, error) {

	op, err := crudOp.dispatchedOp(data.Key.Type, actor)
	if err != nil {
		return nil, fmt.Errorf("on crudDispatcher.Save(%#v): %s", err)
	}

	return op.Save(data, actor)
}

func (crudOp crudDispatcher) Read(key crud.Key, actor auth.Actor) (*crud.Data, error) {
	op, err := crudOp.dispatchedOp(key.Type, actor)
	if err != nil {
		return nil, fmt.Errorf("on crudDispatcher.Read(%#v): %s", err)
	}

	return op.Read(key, actor)
}

func (crudOp crudDispatcher) List(crudType crud.Type, options selectors.Options, actor auth.Actor) ([]crud.Data, error) {
	op, err := crudOp.dispatchedOp(crudType, actor)
	if err != nil {
		return nil, fmt.Errorf("on crudDispatcher.List(%#v): %s", err)
	}

	return op.List(crudType, options, actor)
}

func (crudOp crudDispatcher) Remove(key crud.Key, actor auth.Actor) error {

	op, err := crudOp.dispatchedOp(key.Type, actor)
	if err != nil {
		return fmt.Errorf("on crudDispatcher.Remove(%#v): %s", err)
	}

	return op.Remove(key, actor)
}

func (crudOp crudDispatcher) Types() ([]crud.Type, error) {
	var types []crud.Type

	for t, op := range crudOp.dispatchedOps {
		if op == nil {
			// return nil, fmt.Errorf("on crudDispatcher.Types(): wrong .Operator for type %s", t)
			l.Errorf("on crudDispatcher.Types(): nil .Operator for type %s", t)
			continue
		}
		types = append(types, t)
	}

	return types, nil
}

const onRoles = "on crudDispatcher.Roles()"

func (crudOp *crudDispatcher) Roles() (rbac.Roles, error) {
	var rolesAll rbac.Roles

	for t, crudOpsTyped := range crudOp.dispatchedOps {
		for r, op := range crudOpsTyped {
			if op == nil {
				// return nil, fmt.Errorf("on crudDispatcher.Types(): wrong .Operator for type %s", t)
				l.Errorf(onRoles+": nil .Operator for type %s with role %s", t, r)
				continue
			}

			for _, rPrev := range rolesAll {
				if rPrev == r {
					continue
				}
			}
			rolesAll = append(rolesAll, r)
		}
	}

	return rolesAll, nil
}
