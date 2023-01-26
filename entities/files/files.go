package files

import (
	"time"
)

// record ------------------------------------------------------------

// file --------------------------------------------------------------

type File struct {
	Path      string     `json:",omitempty" bson:",omitempty"`
	IsDir     bool       `json:",omitempty" bson:",omitempty"`
	Size      int64      `json:",omitempty" bson:",omitempty"`
	CreatedAt time.Time  `json:",omitempty" bson:",omitempty"`
	UpdatedAt *time.Time `json:",omitempty" bson:",omitempty"`
}
