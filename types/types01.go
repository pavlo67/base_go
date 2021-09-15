package types

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/rbac"

	"github.com/pavlo67/data/elements/item"
)

// persons -------------------------------------------------------------

type Person struct {
	Nickname         string `    json:",omitempty" bson:",omitempty"`
	rbac.Roles       `           json:",omitempty" bson:",omitempty"`
	auth.Creds       `           json:",omitempty" bson:",omitempty"`
	Info             common.Map `json:",omitempty" bson:",omitempty"`
	item.Description `           json:",inline"    bson:",inline"`
}

// records -------------------------------------------------------------

type Content struct {
	Title   string `json:",omitempty" bson:",omitempty"`
	Summary string `json:",omitempty" bson:",omitempty"`
	Type    string `json:",omitempty" bson:",omitempty"`
	Data    string `json:",omitempty" bson:",omitempty"`
}

type Record struct {
	Content          `          json:",inline" bson:",inline"`
	Embedded         []Content `json:",omitempty" bson:",omitempty"`
	item.Description `          json:",inline" bson:",inline"`
}
