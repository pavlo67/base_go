package entities

import (
	"encoding/json"
	"sort"
	"time"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/rbac"

	"github.com/pavlo67/data/components/vcs"
)

type Type common.IDStr

type Key struct {
	Type
	ID common.IDStr
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

type OperatorCRUD interface {
	Types() ([]Type, error)
	Roles() (rbac.Roles, error)

	Save(Data, auth.Actor) (*Key, vcs.History, error)
	Read(Key, auth.Actor) (*Data, error)
	List(Type, Options, auth.Actor) ([]Data, error)
	Remove(Key, auth.Actor) error
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
