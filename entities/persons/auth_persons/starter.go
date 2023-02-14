package auth_persons

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/auth"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"
	"github.com/pavlo67/data/entities/persons"
)

const InterfaceKey joiner.InterfaceKey = "auth_persons"

var l logger.Operator

var _ starter.Operator = &authPersonsStarter{}

type authPersonsStarter struct {
	personsKey joiner.InterfaceKey

	interfaceKey joiner.InterfaceKey
}

func Starter() starter.Operator {
	return &authPersonsStarter{}
}

func (aps *authPersonsStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (aps *authPersonsStarter) Prepare(cfg *config.Config, options common.Map) error {
	aps.personsKey = joiner.InterfaceKey(options.StringDefault("persons_key", string(persons.InterfaceKey)))
	aps.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(auth.InterfaceKey)))

	return nil
}

const onRun = "on authPersonsStarter.Run(): "

func (aps *authPersonsStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	personsOp, _ := joinerOp.Interface(aps.personsKey).(persons.Operator)
	if personsOp == nil {
		return errors.Errorf(onRun+"no persons.Operator with key %s", aps.personsKey)
	}

	authOp, err := New(personsOp, 10)
	if err != nil {
		return errors.Wrap(err, onRun)
	}

	if err = joinerOp.Join(authOp, aps.interfaceKey); err != nil {
		return errors.Wrapf(err, onRun+"can't join *authPersons{} as auth.Operator with key '%s'", aps.interfaceKey)
	}

	return nil
}
