package persons

import (
	"fmt"
	"strings"

	_ "github.com/GehirnInc/crypt/sha256_crypt"

	"github.com/GehirnInc/crypt"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/errors"
)

var crypter = crypt.SHA256.New()

func hash(value string) (string, error) {
	if value := strings.TrimSpace(value); value == "" {
		return "", errors.New("no value to hash")
	}

	var salt []byte // TODO: generate salt
	hash, err := crypter.Generate([]byte(value), salt)
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("can't crypt.Generate(%s, %s)", value, salt))
	}

	return hash, nil
}

func (item *Item) CredsByKey(key auth.CredsType) interface{} {
	if item == nil {
		return nil
	}

	return item.creds[key]
}

func (item *Item) Creds() auth.Creds {
	if item == nil {
		return nil
	}

	return item.creds
}

const onSetCreds = "on persons.Item.SetCreds()"

func (item *Item) SetCreds(creds auth.Creds) error {
	if item == nil {
		return fmt.Errorf(onSetCreds + ": no item to set creds")
	} else if item.creds == nil {
		item.creds = auth.Creds{}
	}

	for key, value := range creds {
		if key != auth.CredsPassword {
			item.creds[key] = value
			continue
		}

		if value == "" {
			item.creds[auth.CredsPasshash] = ""
		} else if passhash, err := hash(value); err != nil {
			return fmt.Errorf(onSetCreds+": %s", err)
		} else {
			item.creds[auth.CredsPasshash] = passhash
		}
	}

	// log.Printf("set creds: %#v --> %#v", creds, item.creds)

	return nil
}

func (item *Item) CheckCreds(key auth.CredsType, value string) bool {
	if item == nil {
		return false
	}

	if key != auth.CredsPassword {
		return item.creds[key] == value
	}

	// log.Printf("check password (%s) on passhash (%s) --> %s", value, item.creds[auth.CredsPasshash], crypter.Verify(item.creds[auth.CredsPasshash], []byte(value)))

	return crypter.Verify(item.creds[auth.CredsPasshash], []byte(value)) == nil
}

func (item *Item) GetCredsStr(key auth.CredsType) string {
	if item == nil {
		return ""
	}

	return item.creds[key]
}

//type IdentityForMarshallingCreds struct {
//	ID       ID         `json:",omitempty"`
//	Nickname string     `json:",omitempty"`
//	Roles    rbac.Roles `json:",omitempty"`
//	Creds    common.Map `json:",omitempty"`
//}

//func (identity Identity) MarshalJSONCreds() ([]byte, error) {
//	return json.Marshal(identity.creds)
//}
//
//func (identity *Identity) UnmarshalJSONCreds(jsonBytes []byte) error {
//	if len(jsonBytes) < 1 {
//		return nil
//	}
//
//	return json.Unmarshal(jsonBytes, &identity.creds)
//}
