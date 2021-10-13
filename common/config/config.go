package config

import (
	"fmt"

	"github.com/pavlo67/common/common/config"
)

func ConfigOther(envPath, configEnv string, marshaller config.Marshaler) (config.Config, error) {

	cfgServicePath := envPath + configEnv + ".yaml"
	cfgServicePtr, err := config.Get(cfgServicePath, marshaller)
	if err != nil || cfgServicePtr == nil {
		return config.Config{}, fmt.Errorf("on config.ConfigOther(%s, %s) got %#v / %s", envPath, configEnv, cfgServicePtr, err)
	}

	return *cfgServicePtr, nil

}
