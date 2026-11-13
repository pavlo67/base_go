package auth_stub

import (
	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"
)

func Starter() starter.Operator {
	return &authstubStarter{}
}

var l logger.Operator
var _ starter.Operator = &authstubStarter{}

type authstubStarter struct {
	defaultActors []auth.Actor
	interfaceKey  joiner.InterfaceKey
}

func (ahs *authstubStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (ahs *authstubStarter) Run(cfg *config.Envs, options common.Map, joinerOp joiner.Operator, l_ logger.OperatorJ) error {
	l = l_

	if err := cfg.Value("actors", &ahs.defaultActors); err != nil {
		return err
	}

	ahs.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(auth.InterfaceKey)))

	authOp, err := New(ahs.defaultActors)
	if err != nil {
		return errors.Wrap(err, "can't init *authstub{} as auth.Operator")
	}

	if err = joinerOp.Join(authOp, ahs.interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join *authstub{} as auth.Operator with key '%s'", ahs.interfaceKey)
	}

	return nil
}
