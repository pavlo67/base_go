package crud

import (
	"sort"
	"testing"
	"time"

	"github.com/pavlo67/common/common/auth"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/data/elements/selectors"
)

type Type common.IDStr

type Key struct {
	Type
	ID interface{}
}

type ChangeItem func(interface{}, Key) (interface{}, error)

type Operator interface {
	Types() ([]Type, error)

	Save(Key, interface{}, *auth.Identity) (*Key, error)
	Read(Key, *auth.Identity) (interface{}, error)
	List(Type, selectors.Options, *auth.Identity) ([]interface{}, error)
	Remove(Key, *auth.Identity) error

	TestIfEqual(t *testing.T, expectedKey Key, expected interface{}, toCheck interface{}) error
}

type Stat struct {
	ChldCount int64
	TotalSize int64
	CreatedAt time.Time

	Key `json:"-" bson:"-"`
}

type StatMap map[Key]Stat

func (ts StatMap) List(sortBy string) []Stat {
	var statList []Stat
	for key, stat := range ts {
		stat.Key = key
		statList = append(statList, stat)
	}

	switch sortBy {
	case "count":
		sort.Slice(statList, func(i, j int) bool { return statList[i].ChldCount >= statList[j].ChldCount })
	case "size":
		sort.Slice(statList, func(i, j int) bool { return statList[i].TotalSize <= statList[j].TotalSize })
	default:
		sort.Slice(statList, func(i, j int) bool {
			return statList[i].Type <= statList[j].Type
		})
	}

	return statList
}
