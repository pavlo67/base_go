package types

import (
	"time"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/rbac"
)

type File01 struct {
	Path        string        `json:",omitempty" bson:",omitempty"`
	IsDir       bool          `json:",omitempty" bson:",omitempty"`
	Size        int64         `json:",omitempty" bson:",omitempty"`
	CreatedAt   time.Time     `json:",omitempty" bson:",omitempty"`
	Description Description01 `json:",inline"    bson:",inline"`
}

type Person01 struct {
	Nickname    string     `   json:",omitempty" bson:",omitempty"`
	Info        common.Map `   json:",omitempty" bson:",omitempty"`
	rbac.Roles  `              json:",omitempty" bson:",omitempty"`
	auth.Creds  `              json:",omitempty" bson:",omitempty"`
	Description Description01 `json:",inline"    bson:",inline"`
}

type Content01 struct {
	Title   string `json:",omitempty" bson:",omitempty"`
	Summary string `json:",omitempty" bson:",omitempty"`
	Type    string `json:",omitempty" bson:",omitempty"`
	Data    string `json:",omitempty" bson:",omitempty"`
}

type Record struct {
	Content01   `              json:",inline"    bson:",inline"`
	Embedded    []Content01   `json:",omitempty" bson:",omitempty"`
	Description Description01 `json:",inline"    bson:",inline"`
}
