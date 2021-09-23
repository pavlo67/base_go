package persons

import (
	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/data/types"

	"github.com/pavlo67/data/elements/selectors"
)

type ID common.IDStr

const HasEmail selectors.Key = "has_email"
const HasNickname selectors.Key = "has_nickname"

type Operator interface {
	Save(Item, *auth.Identity) (ID, error)
	Read(ID, *auth.Identity) (*Item, error)
	Remove(ID, *auth.Identity) error
	List(*selectors.Term, *auth.Identity) ([]Item, error)
	// Stat(*selectors.Term, *auth.Identity) (crud.StatMap, error)
}

type Item struct {
	ID
	types.Person01
}

//const UserKeyFieldName = "key"
//const EmailFieldName = "email"
//const NicknameFieldName = "nickname"
//const VerifiedFieldName = "verified"
//
//type Item struct {
//	auth.User `bson:",omitempty" json:",omitempty"`
//
//	Allowed  bool           `bson:",omitempty" json:",omitempty"`
//	ToVerify []Verification `bson:",omitempty" json:",omitempty"`
//	History  crud.History   `bson:",omitempty" json:",omitempty"`
//}
//
//type Verification struct {
//	CredsType auth.CredsType `bson:",omitempty" json:",omitempty"`
//	Value     string         `bson:",omitempty" json:",omitempty"`
//	Open      bool           `bson:",omitempty" json:",omitempty"`
//	History   crud.History   `bson:",omitempty" json:",omitempty"`
//}
//
//type Operator interface {
//	Save(Item, *crud.SaveOptions) (identity.Key, error)
//	Remove(identity.Key, *crud.RemoveOptions) error
//
//	Read(identity.Key, *crud.GetOptions) (*Item, error)
//	List(*selectors.Term, *crud.GetOptions) ([]Item, error)
//	Count(*selectors.Term, *crud.GetOptions) (uint64, error)
//
//	CheckPassword(password, passHash string) bool
//
//	Allow() error
//	SetVerification(auth.CredsType, string, bool) error
//	Verify(auth.CredsType, string, common.Errors) error
//}

//func (item *Item) UnfoldFromJSON(id auth.ID, rolesBytes, credsBytes, emailBytes, infoBytes, tagsBytes, urnBytes, historyBytes []byte) error {
//	if item == nil {
//		return errors.New("nil persons.Item to be unfolded")
//	}
//
//	item.Identity.ID = id
//	if len(rolesBytes) > 0 {
//		if err := json.Unmarshal(rolesBytes, &item.Roles); err != nil {
//			return errors.Wrapf(err, "can't unmarshal .Roles (%s)", rolesBytes)
//		}
//	}
//
//	var creds auth.Creds
//	if len(credsBytes) > 0 {
//		if err := json.Unmarshal(credsBytes, &creds); err != nil {
//			return errors.Wrapf(err, "can't unmarshal .creds (%s)", credsBytes)
//		}
//	}
//	if len(emailBytes) > 0 {
//		if creds == nil {
//			creds = auth.Creds{}
//		}
//
//		creds[auth.CredsEmail] = string(emailBytes)
//	}
//	item.SetCreds(creds)
//
//	return item.ItemDescription.UnfoldFromJSON(infoBytes, tagsBytes, urnBytes, historyBytes)
//}
//
//func (item *Item) FoldIntoJSON() (rolesBytes, credsBytes, emailBytes, infoBytes, tagsBytes, historyBytes, urnBytes []byte, err error) {
//	if item == nil {
//		return nil, nil, nil, nil, nil, nil, nil, errors.New("nil persons.Item to be folded")
//	}
//
//	rolesBytes = []byte{} // to satisfy NOT NULL constraint
//	if len(item.Roles) > 0 {
//		if rolesBytes, err = json.Marshal(item.Roles); err != nil {
//			return nil, nil, nil, nil, nil, nil, nil, errors.Wrapf(err, "can't marshal .Roles (%#v)", item.Roles)
//		}
//	}
//
//	creds := item.Creds()
//
//	if email := strings.TrimSpace(creds[auth.CredsEmail]); email != "" {
//		emailBytes = []byte(email)
//	}
//
//	delete(creds, auth.CredsEmail)
//
//	credsBytes = []byte{} // to satisfy NOT NULL constraint
//	if len(creds) > 0 {
//		if credsBytes, err = json.Marshal(creds); err != nil {
//			return nil, nil, nil, nil, nil, nil, nil, errors.Wrapf(err, "can't marshal creds (%#v)", creds)
//		}
//	}
//
//	// TODO!!! append to item.History
//	if infoBytes, tagsBytes, urnBytes, historyBytes, err = item.ItemDescription.FoldIntoJSON(); err != nil {
//		return nil, nil, nil, nil, nil, nil, nil, err
//	}
//
//	return rolesBytes, credsBytes, emailBytes, infoBytes, tagsBytes, historyBytes, urnBytes, nil
//}
