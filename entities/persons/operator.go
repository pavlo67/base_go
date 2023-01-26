package persons

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"

	"github.com/pavlo67/data/entities/contacts"

	"github.com/pavlo67/data/components/crud"
	"github.com/pavlo67/data/components/ns"
	"github.com/pavlo67/data/components/vcs"
)

type ID = common.IDStr

type Person struct {
	Firstnames []string        `json:",omitempty" bson:",omitempty"`
	Middlename string          `json:",omitempty" bson:",omitempty"`
	Lastname   string          `json:",omitempty" bson:",omitempty"`
	Nicknames  []string        `json:",omitempty" bson:",omitempty"`
	Contacts   []contacts.Item `json:",omitempty" bson:",omitempty"`
	Info       common.Map      `json:",omitempty" bson:",omitempty"`
}

type Operator interface {
	Save(Item, auth.Actor) (ID, ns.URN, vcs.History, error)
	Read(ID, auth.Actor) (*Item, error)
	Remove(ID, auth.Actor) error
	List(*crud.Term, auth.Actor) ([]Item, error)

	SetURN(id ID) (ns.URN, error)
}

type Item struct {
	ID
	crud.Description `json:",inline"    bson:",inline"`
	Person
}
