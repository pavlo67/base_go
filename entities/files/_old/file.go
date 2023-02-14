package files

import (
	"time"

	"github.com/pavlo67/punctum/auth"
	"github.com/pavlo67/punctum/confidenter/rights"
	"github.com/pavlo67/punctum/notebook/links"
)

type Data struct {
	Name     string          `bson:"name"                json:"name"`
	MIMEType string          `bson:"mimetype,omitempty"  json:"mimetype,omitempty"`
	Links    []links.Item    `bson:"links,omitempty"     json:"links,omitempty"`
	RView    auth.ID         `bson:"r_view"              json:"r_view"`
	ROwner   auth.ID         `bson:"r_owner"             json:"r_owner"`
	Managers rights.Managers `bson:"managers,omitempty"  json:"managers,omitempty"`
	GlobalIS string          `bson:"global_is,omitempty" json:"global_is,omitempty"`
}

type Item struct {
	Data `bson:",inline" json:",inline"`

	Size     int64     `bson:"size,omitempty"     json:"size,omitempty"`
	IsDir    bool      `bson:"is_dir,omitempty"   json:"is_dir,omitempty"`
	Modified time.Time `bson:"modified,omitempty" json:"modified,omitempty"`
	Content  []byte    `bson:"content,omitempty"  json:"content,omitempty"`
}
