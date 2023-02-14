package records

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/data/entities"

	"github.com/pavlo67/data/components/ns"
	"github.com/pavlo67/data/components/vcs"
)

type ID = common.IDStr

type Content struct {
	Title   string `json:",omitempty" bson:",omitempty"`
	Summary string `json:",omitempty" bson:",omitempty"`
	Type    string `json:",omitempty" bson:",omitempty"`
	Data    string `json:",omitempty" bson:",omitempty"`
}

type Record struct {
	SourceURN ns.URN `   json:",omitempty" bson:",omitempty"`
	Content   `          json:",inline"    bson:",inline"`
	Additions []Content `json:",omitempty" bson:",omitempty"`
}

type Item struct {
	ID
	entities.Description `json:",inline" bson:",inline"`
	Record               `json:",inline" bson:",inline"`
}

type Operator interface {
	Save(Item, auth.Actor) (ID, ns.URN, vcs.History, error)
	Read(ID, auth.Actor) (*Item, error)
	Remove(ID, auth.Actor) error
	List(*entities.Term, auth.Actor) ([]Item, error)

	SetURN(id ID) (ns.URN, error)

	//List(*selectors.Term, *auth.Identity) ([]Item, error)
	//Stat(*selectors.Term, *auth.Identity) (db.StatMap, error)
	//Tags(*selectors.Term, *auth.Identity) (tags.StatMap, error)
	//
	//AddParent(ts []tags.Item, id ID) ([]tags.Item, error)

}

//func ReadWithChildren(recordsOp OperatorCRUD, id ID, identity *auth.Identity) (*Item, []Item, error) {
//	r, err := recordsOp.Read(id, identity)
//	if err != nil {
//		return r, nil, err
//	}
//
//	selectorParent := selectors.Term{
//		Key:    HasParent,
//		Values: []string{string(id)},
//	}
//
//	children, err := recordsOp.List(&selectorParent, identity)
//	return r, children, err
//}
