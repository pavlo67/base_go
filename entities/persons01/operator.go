package persons01

import (
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/data/components/crud"
	"github.com/pavlo67/data/entities"

	"github.com/pavlo67/data/elements/selectors"
)

type ID crud.ID

type Operator interface {
	Save(Item, *auth.Identity) (ID, error)
	Read(ID, *auth.Identity) (*Item, error)
	Remove(ID, *auth.Identity) error
	List(*selectors.Term, *auth.Identity) ([]Item, error)
	// Stat(*selectors.Term, *auth.Identity) (crud.StatMap, error)
}

type Item struct {
	ID
	crud.Description `json:",inline"    bson:",inline"`
	entities.Person01
}
