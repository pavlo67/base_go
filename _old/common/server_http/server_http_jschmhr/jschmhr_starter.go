package server_http_jschmhr

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/server_http"
	"github.com/pavlo67/common/common/starter"
)

func Starter() starter.Operator {
	return &server_http_jschmhrStarter{}
}

var l logger.Operator
var _ starter.Operator = &server_http_jschmhrStarter{}

type server_http_jschmhrStarter struct{}

func (shjs *server_http_jschmhrStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (shjs *server_http_jschmhrStarter) Run(cfg *config.Envs, options common.Map, joinerOp joiner.Operator, l_ logger.OperatorJ) error {
	l = l_
	interfaceKey := joiner.InterfaceKey(options.StringDefault("interface_key", string(server_http.InterfaceKey)))

	configKey := options.StringDefault("config_key", "server_http")
	var config server_http.ConfigStarter
	if err := cfg.Value(configKey, &config); err != nil {
		return err
	}

	// TODO!!! customize it
	var secretENVs []string

	srvOp, err := New(config.Port, config.TLSCertFile, config.TLSKeyFile, secretENVs)
	if err != nil {
		return errors.Wrap(err, "on server_http_jschmhr.New()")
	}

	if err = joinerOp.Join(srvOp, interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join *serverHTTPJschmhr{} as server_http.Operator with key '%s'", interfaceKey)
	}

	return nil

}
