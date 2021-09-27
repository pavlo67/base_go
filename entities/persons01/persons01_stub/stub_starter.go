package persons01_stub

import (
	"fmt"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/config"
	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/joiner"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/starter"
	"github.com/pavlo67/data/entities/persons01"
)

func Starter() starter.Operator {
	return &personsStubStarter{}
}

var l logger.Operator
var _ starter.Operator = &personsStubStarter{}

type personsStubStarter struct {
	interfaceKey joiner.InterfaceKey
	cleanerKey   joiner.InterfaceKey
}

func (pss *personsStubStarter) Name() string {
	return logger.GetCallInfo().PackageName
}

func (pss *personsStubStarter) Prepare(cfg *config.Config, options common.Map) error {

	pss.interfaceKey = joiner.InterfaceKey(options.StringDefault("interface_key", string(persons01.InterfaceKey)))
	pss.cleanerKey = joiner.InterfaceKey(options.StringDefault("cleaner_key", string(persons01.InterfaceCleanerKey)))

	return nil
}

func (pss *personsStubStarter) Run(joinerOp joiner.Operator) error {
	if l, _ = joinerOp.Interface(logger.InterfaceKey).(logger.Operator); l == nil {
		return fmt.Errorf("no logger.Operator with key %s", logger.InterfaceKey)
	}

	personsOp, personsCleanerOp, err := New(nil)
	if err != nil {
		return errors.Wrap(err, "can't init *personsStub{} as persons.Operator")
	}

	if err = joinerOp.Join(personsOp, pss.interfaceKey); err != nil {
		return errors.Wrapf(err, "can't join *personsStub{} as persons.Operator with key '%s'", pss.interfaceKey)
	}

	if err = joinerOp.Join(personsCleanerOp, pss.cleanerKey); err != nil {
		return errors.Wrapf(err, "can't join *personsStub{} as db.Cleaner with key '%s'", pss.cleanerKey)
	}

	return nil
}
