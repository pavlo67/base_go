package encrlib

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/pkg/errors"
)

const onNewRSAPrivateKey = "on encrlib.NewRSAPrivateKey()"

func NewRSAPrivateKey(pathToStore string) (*rsa.PrivateKey, error) {
	if pathToStore != "" {
		if _, err := os.Stat(pathToStore); !os.IsNotExist(err) {
			keyJSON, err := ioutil.ReadFile(pathToStore)
			if err != nil {
				return nil, errors.Wrapf(err, onNewRSAPrivateKey+": can't read file (%s)", pathToStore)
			}

			var privateKey rsa.PrivateKey
			err = json.Unmarshal(keyJSON, &privateKey)
			if err != nil {
				return nil, errors.Wrapf(err, onNewRSAPrivateKey+": can't .json.Unmarshal file (%s --> %s)", pathToStore, keyJSON)
			}

			return &privateKey, nil
		}
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, errors.Wrap(err, onNewRSAPrivateKey)
	}

	if privateKey == nil {
		return nil, errors.New(onNewRSAPrivateKey + ": nil key was generated")
	}

	if pathToStore != "" {
		keyJSON, err := json.Marshal(privateKey)
		if err != nil {
			return nil, errors.Wrapf(err, onNewRSAPrivateKey+": can't .json.Marshal key (%#v)", privateKey)
		}

		if err = ioutil.WriteFile(pathToStore, keyJSON, 0644); err != nil {
			return nil, errors.Wrapf(err, onNewRSAPrivateKey+": can't write file (%s)", pathToStore)
		}
	}

	return privateKey, nil
}
