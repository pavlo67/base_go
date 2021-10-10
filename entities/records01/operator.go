package records01

import (
	"github.com/pavlo67/data/common/auth"

	"github.com/pavlo67/data/elements/selectors"

	"github.com/pavlo67/data/entities"

	"github.com/pavlo67/data/components/crud"
)

type ID = crud.ID

type Operator interface {
	Save(Item, auth.Actor) (ID, error)
	Read(ID, auth.Actor) (*Item, error)
	Remove(ID, auth.Actor) error
	List(*selectors.Term, auth.Actor) ([]Item, error)
}

type Item struct {
	ID
	crud.Description `json:",inline"    bson:",inline"`
	entities.Record01
}
