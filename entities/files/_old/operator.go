package files

import (
	"github.com/pavlo67/partes/crud"

	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/starter/joiner"
)

const InterfaceKey joiner.InterfaceKey = "files"

const LinkType = "file"
const RepoSchema = "repo:"

type Operator interface {
	Create(userIS auth.ID, file *Item) (string, error)

	Read(userIS auth.ID, name string) (*Item, error)

	ReadList(userIS auth.ID, options *content.ListOptions) ([]Item, uint64, error)

	Update(userIS auth.ID, file *Item) (crud.Result, error)

	Delete(userIS auth.ID, name string) (crud.Result, error)

	Close()
}
