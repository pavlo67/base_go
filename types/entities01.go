package types

import (
	"time"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/data/elements/contacts"
)

// person ------------------------------------------------------------

type Person01 struct {
	Firstnames  []string        `json:",omitempty" bson:",omitempty"`
	Middlename  string          `json:",omitempty" bson:",omitempty"`
	Lastname    string          `json:",omitempty" bson:",omitempty"`
	Nicknames   []string        `json:",omitempty" bson:",omitempty"`
	Contacts    []contacts.Item `json:",omitempty" bson:",omitempty"`
	Info        common.Map      `json:",omitempty" bson:",omitempty"`
	Description Description01   `json:",inline"    bson:",inline"`
}

// record ------------------------------------------------------------

type Content01 struct {
	Title   string `json:",omitempty" bson:",omitempty"`
	Summary string `json:",omitempty" bson:",omitempty"`
	Type    string `json:",omitempty" bson:",omitempty"`
	Data    string `json:",omitempty" bson:",omitempty"`
}

type Record01 struct {
	Content01   `              json:",inline"    bson:",inline"`
	Embedded    []Content01   `json:",omitempty" bson:",omitempty"`
	Description Description01 `json:",inline"    bson:",inline"`
}

// file --------------------------------------------------------------

type File01 struct {
	Path      string     `json:",omitempty" bson:",omitempty"`
	IsDir     bool       `json:",omitempty" bson:",omitempty"`
	Size      int64      `json:",omitempty" bson:",omitempty"`
	CreatedAt time.Time  `json:",omitempty" bson:",omitempty"`
	UpdatedAt *time.Time `json:",omitempty" bson:",omitempty"`
}
