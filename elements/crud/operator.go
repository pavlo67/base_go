package crud

import (
	"time"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/data/elements/ns"
	"github.com/pavlo67/data/elements/selectors"
)

type Type common.IDStr
type ID common.IDStr

type Key struct {
	Type
	ID
}

type Operator interface {
	Types() ([]Type, error)

	Save(Key, interface{}) (*Key, error)
	Read(Key) (interface{}, error)
	List(Type, selectors.Options) ([]interface{}, error)
	Remove(Key) error

	CheckIfEqual(expectedKey Key, expected interface{}, toCheck interface{}) error
}

type Stat struct {
	NSS       ns.NSS
	ChldCount int64
	TotalSize int64
	CreatedAt time.Time
}

type StatMap map[string]Stat
