package files

// operating files or "records about files" in database

import (
	"github.com/pavlo67/common/common/errors"
	"time"
)

var err = errors.Wrap(nil)

type File struct {
	Path     string    `json:",omitempty" bson:",omitempty"` // unique key
	IsDir    bool      `json:",omitempty" bson:",omitempty"`
	Size     uint64    `json:",omitempty" bson:",omitempty"`
	CTime    time.Time `json:",omitempty" bson:",omitempty"`
	MTime    time.Time `json:",omitzero" bson:",omitempty"`
	MimeType string    `json:",omitempty" bson:",omitempty"`
}

type Item struct {
	File      `          json:",inline" bson:",inline"`
	CreatedAt time.Time `json:",omitzero" bson:",omitempty"`
	UpdatedAt time.Time `json:",omitzero" bson:",omitempty"`
}

type Operator interface {
	// creates new or replaces existing Item's record
	Save(file File) error
	Read(path string) (*Item, error)
	Remove(path string) error
	List(path string, depth int) ([]Item, error)
}
