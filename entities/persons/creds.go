package persons

//import (
//	"fmt"
//	"strings"
//
//	"github.com/GehirnInc/crypt"
//	_ "github.com/GehirnInc/crypt/sha256_crypt"
//
//	"github.com/pavlo67/common/common/auth"
//	"github.com/pavlo67/common/common/errors"
//
//	"github.com/pavlo67/data/types"
//)
//
//var crypter = crypt.SHA256.New()
//
//func hash(value string) (string, error) {
//	if value := strings.TrimSpace(value); value == "" {
//		return "", errors.New("no value to hash")
//	}
//
//	var salt []byte // TODO: generate salt
//	hash, err := crypter.Generate([]byte(value), salt)
//	if err != nil {
//		return "", errors.Wrap(err, fmt.Sprintf("can't crypt.Generate(%s, %s)", value, salt))
//	}
//
//	return hash, nil
//}
//
//func (item *types.Person01) CredsByKey(key auth.CredsType) interface{} {
//	if item == nil {
//		return nil
//	}
//
//	return item.creds[key]
//}
//
//func (item *types.Person01) Creds() auth.Creds {
//	if item == nil {
//		return nil
//	}
//
//	return item.creds
//}
//
//const onSetCreds = "on persons.Item.SetCreds()"
//
//func (item *types.Person01) SetCreds(creds auth.Creds) error {
//	if item == nil {
//		return fmt.Errorf(onSetCreds + ": no item to set creds")
//	} else if item.creds == nil {
//		item.creds = auth.Creds{}
//	}
//
//	for key, value := range creds {
//		if key != auth.CredsPassword {
//			item.creds[key] = value
//			continue
//		}
//
//		if value == "" {
//			item.creds[auth.CredsPasshash] = ""
//		} else if passhash, err := hash(value); err != nil {
//			return fmt.Errorf(onSetCreds+": %s", err)
//		} else {
//			item.creds[auth.CredsPasshash] = passhash
//		}
//	}
//
//	// log.Printf("set creds: %#v --> %#v", creds, item.creds)
//
//	return nil
//}
//
//func (item *types.Person01) CheckCreds(key auth.CredsType, value string) bool {
//	if item == nil {
//		return false
//	}
//
//	if key != auth.CredsPassword {
//		return item.creds[key] == value
//	}
//
//	// log.Printf("check password (%s) on passhash (%s) --> %s", value, item.creds[auth.CredsPasshash], crypter.Verify(item.creds[auth.CredsPasshash], []byte(value)))
//
//	return crypter.Verify(item.creds[auth.CredsPasshash], []byte(value)) == nil
//}
//
//func (item *types.Person01) GetCredsStr(key auth.CredsType) string {
//	if item == nil {
//		return ""
//	}
//
//	return item.creds[key]
//}
