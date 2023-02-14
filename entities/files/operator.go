package files

import (
	"time"

	"github.com/pavlo67/common/common/auth"
)

type File struct {
	Path      string     `json:",omitempty" bson:",omitempty"`
	IsDir     bool       `json:",omitempty" bson:",omitempty"`
	Size      int64      `json:",omitempty" bson:",omitempty"`
	CreatedAt time.Time  `json:",omitempty" bson:",omitempty"`
	UpdatedAt *time.Time `json:",omitempty" bson:",omitempty"`
}

type Operator interface {
	Save(path, newFilePattern string, data []byte, identity *auth.Identity) (string, error)
	Read(path string, identity *auth.Identity) ([]byte, error)
	Remove(path string, identity *auth.Identity) error
	List(path string, depth int, identity *auth.Identity) ([]File, error)
	Stat(path string, depth int, identity *auth.Identity) (*File, error)
}
