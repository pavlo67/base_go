package crud

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/pavlo67/common/common"

	"github.com/pavlo67/data/common/auth"

	"github.com/pavlo67/data/elements/selectors"
)

type Type common.IDStr

type Key struct {
	Type
	ID
}

type Data struct {
	Key
	Description
	Value interface{}
}

type DataRaw struct {
	Key
	Description
	Value json.RawMessage
}

type Operator interface {
	Types() ([]Type, error)

	Save(Data, auth.Actor) (*Key, error)
	Read(Key, auth.Actor) (*Data, error)
	List(Type, selectors.Options, auth.Actor) ([]Data, error)
	Remove(Key, auth.Actor) error

	// TestIfEqual(t *testing.T, expectedKey Key, checkIfEqual test.CheckIfEqual) error
}

type ChangeItemForTest func(Data, Key) (*Data, error)
type ReadValueRaw func(message json.RawMessage) (interface{}, error)

type Stat struct {
	ChldCount int64
	TotalSize int64
	CreatedAt time.Time
	UpdatedAt *time.Time

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
