package records

import (
	"github.com/pavlo67/common/common/auth"

	"github.com/pavlo67/data/components/crud"
	"github.com/pavlo67/data/components/ns"
	"github.com/pavlo67/data/components/selectors"
	"github.com/pavlo67/data/components/vcs"
)

type ID = crud.ID

type Content struct {
	Title   string `json:",omitempty" bson:",omitempty"`
	Summary string `json:",omitempty" bson:",omitempty"`
	Type    string `json:",omitempty" bson:",omitempty"`
	Data    string `json:",omitempty" bson:",omitempty"`
}

type Record struct {
	Content  `          json:",inline"    bson:",inline"`
	Embedded []Content `json:",omitempty" bson:",omitempty"`
}

type Operator interface {
	Save(Item, auth.Actor) (ID, ns.URN, vcs.History, error)
	Read(ID, auth.Actor) (*Item, error)
	Remove(ID, auth.Actor) error
	List(*selectors.Term, auth.Actor) ([]Item, error)

	SetURN(id ID) (ns.URN, error)
}

type Item struct {
	ID

	crud.Description `json:",inline" bson:",inline"`
	Record           `json:",inline" bson:",inline"`
}
