package crud_dispatcher

import (
	"fmt"

	"github.com/pavlo67/common/common/auth"

	"github.com/pavlo67/data/elements/selectors"

	"github.com/pavlo67/data/components/crud"
)

var _ crud.Operator = &crudDispatcher{}

type CrudOps map[crud.Type]crud.Operator

type crudDispatcher struct {
	crudOps    CrudOps
	MessagesOp crud.Operator
	StorageOp  crud.Operator
}

const onNew = "on crudDispatcher.New()"

func New(crudOps CrudOps) (crud.Operator, error) {
	crudOp := crudDispatcher{
		crudOps: crudOps,
	}

	return &crudOp, nil

}

// operator ----------------------------------------------------------------------------------------------------------------

var _ crud.Operator = &crudDispatcher{}

func (crudOp crudDispatcher) Save(data crud.Data, identity *auth.Identity) (*crud.Key, error) {
	op := crudOp.crudOps[data.Key.Type]
	if op == nil {
		return nil, fmt.Errorf("on crudDispatcher.Save(%#v): wrong key.Type", data.Key)
	}

	return op.Save(data, identity)
}

func (crudOp crudDispatcher) Read(key crud.Key, identity *auth.Identity) (*crud.Data, error) {
	op := crudOp.crudOps[key.Type]
	if op == nil {
		return nil, fmt.Errorf("on crudDispatcher.Read(%#v): wrong key.Type", key)
	}

	return op.Read(key, identity)
}

func (crudOp crudDispatcher) List(crudType crud.Type, options selectors.Options, identity *auth.Identity) ([]crud.Data, error) {
	op := crudOp.crudOps[crudType]
	if op == nil {
		return nil, fmt.Errorf("on crudDispatcher.List(%#v): wrong crudType", crudType)
	}

	return op.List(crudType, options, identity)
}

func (crudOp crudDispatcher) Remove(key crud.Key, identity *auth.Identity) error {

	op := crudOp.crudOps[key.Type]
	if op == nil {
		return fmt.Errorf("on crudDispatcher.Remove(%#v): wrong key.Type", key)
	}

	return op.Remove(key, identity)
}

func (crudOp crudDispatcher) Types() ([]crud.Type, error) {
	var types []crud.Type

	for t, op := range crudOp.crudOps {
		if op == nil {
			return nil, fmt.Errorf("on crudDispatcher.Types(): wrong .Operator for type %s", t)
		}
		types = append(types, t)
	}

	return types, nil
}
