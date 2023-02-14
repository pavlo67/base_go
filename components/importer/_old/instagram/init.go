package instagramimporter

import (
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/config"
	"github.com/pavlo67/punctum/basis/program"
	"github.com/pavlo67/punctum/interfaces/starter"
	"github.com/pkg/errors"
)

// Starter ...
func Starter() starter.Operator {
	return &instagramComponent{}
}

type instagramComponent struct {
	token string
}

const InterfaceKey = "importer.instagram"

func (fl *instagramComponent) Name() string {
	return InterfaceKey
}

func (ic *instagramComponent) Check(conf config.PunctumConfig, componentKeys config.ComponentKeys, indexPath string) ([]starter.Info, error) {

	var errs []error
	ic.token, errs = conf.Instagram("access_token", errs)

	return nil, basis.JoinErrors(errs)
}

func (ic *instagramComponent) Setup(conf config.PunctumConfig, componentKeys config.ComponentKeys, indexPath string, data map[string]string) error {
	return nil
}

func (ic *instagramComponent) Init() error {
	instagramOp := &Instagram{
		Token: ic.token,
	}

	err := program.JoinInterface(instagramOp, InterfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join instagram importer")
	}

	return nil
}
