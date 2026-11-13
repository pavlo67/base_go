package encrlib

import (
	"errors"
)

var ErrBadCryptype = errors.New("bad cryptography type")

type Cryptype string

const (
	SHA256  Cryptype = "SHA256"
	Provos  Cryptype = "Provos"
	NoCrypt Cryptype = ""
)

//const CryptypePreferred Cryptype = SHA256

//func PasswordValidation(password string, minLength int) (string, error) {
//	password = strings.TrimSpace(password)
//
//	if minLength > 0 && len(password) < minLength {
//		return "", fmt.Errorf("закороткий пароль, повинно бути не менше %d символів", minLength)
//	}
//	return password, nil
//}
