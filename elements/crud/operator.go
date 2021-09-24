package crud

import (
	"sort"
	"time"

	"github.com/pavlo67/common/common"
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
			return statList[i].Type <= statList[j].Type || (statList[i].Type == statList[j].Type && statList[i].ID <= statList[j].ID)
		})
	}

	return statList
}
