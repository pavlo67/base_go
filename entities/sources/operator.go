package sources

import (
	"time"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/joiner"

	"github.com/pavlo67/data/entities"

	"github.com/pavlo67/data/components/ns"
	"github.com/pavlo67/data/components/vcs"
)

type ID = common.IDStr

type Source struct {
	SourceURN            ns.URN              `json:",omitempty"`
	Title                string              `json:",omitempty"`
	ImporterInterfaceKey joiner.InterfaceKey `json:",omitempty"`
	Params               common.Map          `json:",omitempty"` // for Create/Update methods for ex. tags list to set them on each imported item
}

type Item struct {
	ID                   `json:",omitempty"`
	Source               `json:",inline"`
	entities.Description `json:",inline"`
}

type Stat struct {
	Start           time.Time
	Duration        time.Duration
	RecordsTotalNum int
	RecordsSavedNum int
	ErrorsNum       int
	LastError       error
}

type Operator interface {
	Save(Item, auth.Actor) (ID, ns.URN, vcs.History, error)
	Read(ID, auth.Actor) (*Item, error)
	Remove(ID, auth.Actor) error
	List(*entities.Term, auth.Actor) ([]Item, error)
	AddStat(Stat, auth.Actor) error
}
