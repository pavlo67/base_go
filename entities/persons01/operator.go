package persons01

import (
	"github.com/pavlo67/common/common/auth"

	"github.com/pavlo67/data/elements/selectors"

	"github.com/pavlo67/data/types"
)

type ID interface{}

type Operator interface {
	Save(Item, *auth.Identity) (ID, error)
	Read(ID, *auth.Identity) (*Item, error)
	Remove(ID, *auth.Identity) error
	List(*selectors.Term, *auth.Identity) ([]Item, error)
	// Stat(*selectors.Term, *auth.Identity) (crud.StatMap, error)
}

type Item struct {
	ID
	types.Person01
}
