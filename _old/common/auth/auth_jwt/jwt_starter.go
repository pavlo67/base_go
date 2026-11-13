package auth_jwt

import (
	"fmt"
	"strings"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"
)

const InterfaceKey joiner.InterfaceKey = "auth_jwt"

func Starter() starter.Operator {
	return &identity_jwtStarter{}
}

var l logger.Operator
var _ starter.Operator = &identity_jwtStarter{}

type identity_jwtStarter struct {
	interfaceKey joiner.InterfaceKey
	// interfaceSetCredsKey  joiner.InterfaceKey
	keyPath string
}

func (ss *identity_jwtStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ss *identity_jwtStarter) Run(cfg *config.Envs, options common.Map, joinerOp joiner.Operator, l_ logger.OperatorJ) error {
	l = l_

	var cfgAuthJWT common.Map
	if err := cfg.Value("auth_jwt", &cfgAuthJWT); err != nil {
		return err
	}

	// var errs basis.multipleErrors
	ss.keyPath = strings.TrimSpace(cfgAuthJWT.StringDefault("key_path", ""))
	if ss.keyPath == "" {
		ss.keyPath = "./"
	} else if ss.keyPath[len(ss.keyPath)-1] != '/' {
		ss.keyPath += "/"
	}

	ss.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))
	// ss.interfaceSetCredsKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(InterfaceKey)))

	authOp, err := New(ss.keyPath + "jwt.key")
	if err != nil || authOp == nil {
		return errors.CommonError(err, fmt.Sprintf("got %#v", authOp))
	}

	if err = joinerOp.Join(authOp, ss.interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join auth_jwt as auth.Operator with key '%s'", ss.interfaceKey)
	}

	//if ss.interfaceKey != ss.interfaceSetCredsKey {
	//	if err = joinerOp.Join(authOp, ss.interfaceSetCredsKey); err != nil {
	//		return errors.Wrapf(err, "can't join auth_jwt as auth.Operator with key '%s'", ss.interfaceSetCredsKey)
	//	}
	//}

	return nil
}
