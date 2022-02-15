package records01

import (
	"github.com/pavlo67/common/common/auth"

	"github.com/pavlo67/data/entities"
	"github.com/pavlo67/data/entities/crud01"

	"github.com/pavlo67/data/components/crud"
	"github.com/pavlo67/data/components/selectors"
	"github.com/pavlo67/data/components/vcs"
)

type ID = crud.ID

type Operator interface {
	Save(Item, auth.Actor) (ID, vcs.History, error)
	Read(ID, auth.Actor) (*Item, error)
	Remove(ID, auth.Actor) error
	List(*selectors.Term, auth.Actor) ([]Item, error)
}

type Item struct {
	ID                 `json:",omitempty" bson:",omitempty"`
	crud01.Description `json:",inline"    bson:",inline"`
	entities.Record01
}
