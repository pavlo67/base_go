package records01

import (
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/data/components/crud"
	crud012 "github.com/pavlo67/data/components/crud01"

	"github.com/pavlo67/data/elements/selectors"

	"github.com/pavlo67/data/entities"
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
	crud012.Description `json:",inline"    bson:",inline"`
	entities.Record01
}
