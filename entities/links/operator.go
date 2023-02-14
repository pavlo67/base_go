package links

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/workshop/components/selectors"
)

const InterfaceKey joiner.InterfaceKey = "links"
const CleanerInterfaceKey joiner.InterfaceKey = "links.cleaner"

type LinkedInfo struct {
	ObjectID    string
	CountLinked uint
}

type Linked struct {
	LinkedType string
	LinkedID   string
	Type       string
	Tag        string
	ObjectID   string
}

type Tag string

type TagInfo struct {
	Tag
	Count uint64
}

type Item struct {
	ID     string    `bson:"id,omitempty"       json:"id,omitempty"`
	Type   string    `bson:"type,omitempty"     json:"type,omitempty"`
	Name   string    `bson:"name,omitempty"     json:"name,omitempty"`
	To     string    `bson:"to,omitempty"       json:"to,omitempty"`
	RView  common.ID `bson:"r_view,omitempty"   json:"r_view,omitempty"`
	ROwner common.ID `bson:"r_owner,omitempty"  json:"r_owner,omitempty"`
}

type Links []Item

func (Links) Tags() []string {
	return nil
}

type Operator interface {

	// SetLinks corrects link database after tagged entity is created or changed.
	SetLinks(userIS common.ID, linkedType, linkedID string, newLinks Links) ([]LinkedInfo, error)

	// Query selects all tagged entities with selector and without rights check (it should be done later).
	Query(userIS common.ID, selector *selectors.Term) ([]Linked, error)

	// QueryByTag selects all tagged entities without rights check (it should be done later).
	QueryByTag(userIS common.ID, tag string) ([]Linked, error)

	// QueryByObjectID selects all entities linked to selected object_id without rights check (it should be done later).
	QueryByObjectID(userIS common.ID, id string) ([]Linked, error)

	// QueryTags selects all tags.comp with selector accordingly to user's rights.
	QueryTags(userIS common.ID, selector *selectors.Term) ([]TagInfo, error)

	// QueryTagsByOwner selects all tags.comp accordingly to rOwner and user's rights.
	QueryTagsByOwner(userIS common.ID, rOwner common.ID) ([]TagInfo, error)

	Close() error
}
