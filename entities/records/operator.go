package records

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/data/entities"

	"github.com/pavlo67/data/components/ns"
)

type ID = common.IDStr

type Content struct {
	Title   string `json:",omitempty"`
	Summary string `json:",omitempty"`
	Type    string `json:",omitempty"`
	Data    string `json:",omitempty"`
}

type Record struct {
	SourceURN            ns.URN `   json:",omitempty"`
	Content              `          json:",inline"   `
	Additions            []Content `json:",omitempty"`
	entities.Description `json:",inline"`
}

type Item struct {
	ID             `json:",omitempty"`
	Record         `json:",inline"`
	entities.Point `json:",inline"`
}

type Operator interface {
	Add(Record, auth.Actor) (ID, ns.URN, error)
	Update(Record, ID, auth.Actor) error
	Read(ID, auth.Actor) (*Item, error)
	Remove(ID, auth.Actor) error
	List(*entities.Term, auth.Actor) ([]Item, error)

	// setURN(id ID) (ns.URN, error)

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
