package config

import (
	"io/ioutil"

	"github.com/pavlo67/common/common/serialization"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/errors"
)

type Envs struct {
	Data      common.Map
	Marshaler serialization.Marshaler
}

var errNoEnvs = errors.New("no envs")

func (c *Envs) Raw(key string) (interface{}, error) {
	if c == nil {
		return nil, errNoEnvs
	}

	valueRaw, ok := c.Data[key]
	if !ok {
		return nil, errors.CommonError(common.NotFoundKey, common.Map{"reason": "no key in envs", "key": key})
	}

	return valueRaw, nil
}

func (c *Envs) Value(key string, target interface{}) error {
	if c == nil {
		return errNoEnvs
	}

	if value, ok := c.Data[key]; ok {
		valueRaw, err := c.Marshaler.Marshal(value)
		if err != nil {
			return errors.Wrapf(err, "can't marshal value (%s / %#v) to raw bytes", key, value)
		}

		return c.Marshaler.Unmarshal(valueRaw, target)
	}

	return errors.CommonError(common.NotFoundKey, common.Map{"reason": "no key in envs", "key": key})
}

// -----------------------------------------------------------------------------

func Get(envsFile string, marshaler serialization.Marshaler) (*Envs, error) {

	if len(envsFile) < 1 {
		return nil, errors.New("empty envs path")
	}

	dataBytes, err := ioutil.ReadFile(envsFile)
	if err != nil {
		return nil, errors.Wrapf(err, "can't read envs file from '%s'", envsFile)
	}

	envs := Envs{Marshaler: marshaler}

	if err = marshaler.Unmarshal(dataBytes, &envs.Data); err != nil {
		return nil, errors.Wrapf(err, "can't .Unmarshal('%s') from envs '%s'", dataBytes, envsFile)
	}

	return &envs, nil
}
