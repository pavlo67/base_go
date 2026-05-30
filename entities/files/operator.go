package files

// operating files or "records about files" in database

import (
	"time"
)

type Data struct {
	Path     string    `json:",omitempty" bson:",omitempty"` // unique key
	IsDir    bool      `json:",omitempty" bson:",omitempty"`
	Size     uint64    `json:",omitempty" bson:",omitempty"`
	CTime    time.Time `json:",omitempty" bson:",omitempty"`
	MTime    time.Time `json:",omitzero"  bson:",omitempty"`
	CRC      *int64    `json:",omitempty" bson:",omitempty"`
	MimeType string    `json:",omitempty" bson:",omitempty"`
}

type Item struct {
	Data      `          json:",inline" bson:",inline"`
	CreatedAt time.Time `json:",omitzero" bson:",omitempty"`
	UpdatedAt time.Time `json:",omitzero" bson:",omitempty"`
}

type Operator interface {
	// creates new or replaces existing Item's record
	Save(data Data) error
	Read(path string) (*Item, error)
	Remove(path string, forceRecursion bool) error
	List(path string, depth int) ([]Item, error)
}
