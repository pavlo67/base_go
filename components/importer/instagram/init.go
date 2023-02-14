package instagramimporter

import (
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"
	"github.com/pkg/errors"
)

// Starter ...
func Starter() starter.Operator {
	return &instagramComponent{}
}

type instagramComponent struct {
	token string
}

const InterfaceKey joiner.InterfaceKey = "importer.instagram"

func (fl *instagramComponent) Name() string {
	return string(InterfaceKey)
}

func (ic *instagramComponent) Check(conf config.Config, indexPath string) ([]joiner.Info, error) {

	confInstagram, errs := conf.Google("", nil)
	if errs != nil {
		return nil, errs.Err()
	}
	ic.token = confInstagram["access_token"]

	return nil, errs.Err()
}

func (ic *instagramComponent) Setup(conf config.Config, indexPath string, data map[string]string) error {
	return nil
}

func (ic *instagramComponent) Init() error {
	instagramOp := &Instagram{
		Token: ic.token,
	}

	err := joiner.JoinInterface(instagramOp, InterfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join instagram importer")
	}

	return nil
}
