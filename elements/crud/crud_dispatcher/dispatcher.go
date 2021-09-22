package crud_dispatcher

import (
	"fmt"

	"github.com/pavlo67/data/elements/selectors"

	"github.com/pavlo67/data/elements/crud"
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

func (crudOp crudDispatcher) Save(key crud.Key, data interface{}) (*crud.Key, error) {
	op := crudOp.crudOps[key.Type]
	if op == nil {
		return nil, fmt.Errorf("on crudDispatcher.Save(%#v): wrong key.Type", key)
	}

	return op.Save(key, data)
}

func (crudOp crudDispatcher) Read(key crud.Key) (interface{}, error) {
	op := crudOp.crudOps[key.Type]
	if op == nil {
		return nil, fmt.Errorf("on crudDispatcher.Read(%#v): wrong key.Type", key)
	}

	return op.Read(key)
}

func (crudOp crudDispatcher) List(crudType crud.Type, options selectors.Options) ([]interface{}, error) {
	op := crudOp.crudOps[crudType]
	if op == nil {
		return nil, fmt.Errorf("on crudDispatcher.List(%#v): wrong crudType", crudType)
	}

	return op.List(crudType, options)
}

func (crudOp crudDispatcher) Remove(key crud.Key) error {

	op := crudOp.crudOps[key.Type]
	if op == nil {
		return fmt.Errorf("on crudDispatcher.Remove(%#v): wrong key.Type", key)
	}

	return op.Remove(key)
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

func (crudOp crudDispatcher) CheckIfEqual(expectedKey crud.Key, expected interface{}, toCheck interface{}) error {
	op := crudOp.crudOps[expectedKey.Type]
	if op == nil {
		return fmt.Errorf("on crudDispatcher.Remove(%#v): wrong key.Type", expectedKey)
	}

	return op.CheckIfEqual(expectedKey, expected, toCheck)
}
