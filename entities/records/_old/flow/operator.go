package flow_old

import (
	"time"

	"github.com/pavlo67/punctum/basis/program"
	"github.com/pavlo67/punctum/interfaces"
	"github.com/pavlo67/punctum/interfaces/confidenter"
	"github.com/pavlo67/punctum/interfaces/crud"
)

const InterfaceKey = "flow"

type Item struct {
	ID         string
	RView      confidenter.IdentityString
	ROwner     confidenter.IdentityString
	FountIS    confidenter.IdentityString
	FountURL   string
	OriginalID string
	Original   interface{}
	URL        string
	Title      string
	Summary    string
	Content    string
	ImportedTo string
	CreatedAt  time.Time
	Media      *ItemMedia
}

type ItemPicture struct {
	ImageUrl string `json:"image_url"`
	HREFUrl  string `json:"href_url"`
}

type ItemMedia struct {
	HashTags []string      `json:"hash_tag"`
	Pictures []ItemPicture `json:"pictures"`
}

func (i Item) IdentityString() confidenter.Identity {
	return confidenter.Identity{
		Domain: program.Domain(),
		Path:   "flow.item",
		ID:     i.ID,
	}
}

type Operator interface {
	Create(identity *confidenter.Identity, item Item) (confidenter.Identity, error)

	Read(identity *confidenter.Identity, itemIS confidenter.IdentityString) (*Item, error)

	ReadAll(identity *confidenter.Identity, options *crud.ReadAllOptions, selector interfaces.Selector) ([]Item, int64, error)

	Update(identity *confidenter.Identity, itemIS confidenter.IdentityString, item Item) (crud.Result, error)

	Delete(identity *confidenter.Identity, itemIS confidenter.IdentityString) (crud.Result, error)

	IsNew(item Item) (bool, error)

	ImportTo(identity *confidenter.Identity, id int64, importIS string) error

	Close()
}
