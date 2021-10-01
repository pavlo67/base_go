package records01

import (
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/data/entities"

	"github.com/pavlo67/data/elements/selectors"
)

type ID interface{}

type Operator interface {
	Save(Item, *auth.Identity) (ID, error)
	Read(ID, *auth.Identity) (*Item, error)
	Remove(ID, *auth.Identity) error
	List(*selectors.Term, *auth.Identity) ([]Item, error)
}

type Item struct {
	ID
	entities.Record01
}
