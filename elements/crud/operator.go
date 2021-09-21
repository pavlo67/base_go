package crud

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/data/elements/selectors"
)

type ID common.IDStr

type Operator interface {
	Save(ID, interface{}) (ID, error)
	Read(ID) (interface{}, error)
	List(selectors.Options) ([]interface{}, error)
	Delete(ID) error

	CheckIfEqual(expected interface{}, expectedID ID, toCheck interface{}) error
}
