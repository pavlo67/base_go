package sources

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/crud"
	"github.com/pavlo67/common/common/identity"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/selectors"
)

const InterfaceKey joiner.InterfaceKey = "sources"
const CollectionDefault = "sources"

type Item struct {
	ID      common.ID `bson:"_id,omitempty" json:",omitempty"`
	Key     identity.Key
	URL     string
	Title   string
	Type    joiner.InterfaceKey
	Params  common.Map    // for Create/Update methods for ex. tags list to set them on each imported item
	History []crud.Action `bson:",omitempty" json:",omitempty"`
}

type Operator interface {
	Save(Item, *crud.SaveOptions) (common.ID, error)
	Remove(common.ID, *crud.RemoveOptions) error
	Read(common.ID, *crud.GetOptions) (*Item, error)
	List(*selectors.Term, *crud.GetOptions) ([]Item, error)

	AddHistory(common.ID, crud.History, *crud.SaveOptions) (crud.History, error)
}
