package google_plus

import (
	"github.com/pavlo67/punctum/basis/config"
	"github.com/pavlo67/punctum/basis/program"
	"github.com/pavlo67/punctum/interfaces/starter"
	"github.com/pkg/errors"
)

// Starter ...
func Starter() starter.Operator {
	return &plusComponent{}
}

type plusComponent struct {
	apiKey string
	//apiID      string
	//apiSecret  string
	//pathToJSON string

	conf   config.PunctumConfig
	params map[string]string
}

const InterfaceKey = "importer.google_plus"

func (fl *plusComponent) Name() string {
	return InterfaceKey
}

func (pc *plusComponent) Check(conf config.PunctumConfig, componentKeys config.ComponentKeys, indexPath string) ([]starter.Info, error) {

	var errs []error
	pc.apiKey, errs = conf.Google("api_key", errs)
	//pc.apiID, errs = conf.Google("user_id", errs)
	//pc.apiSecret, errs = conf.Google("secret", errs)
	pc.conf = conf
	return nil, nil
}

func (pc *plusComponent) Setup(conf config.PunctumConfig, componentKeys config.ComponentKeys, indexPath string, data map[string]string) error {
	pc.conf = conf
	return nil
}

func (pc *plusComponent) Init() error {
	plusOp := &PLUS{
		ApiKey: pc.apiKey,
		//ApiID:      pc.apiID,
		//ApiSecret:  pc.apiSecret,
		//PathToJSON: pc.pathToJSON,

	}

	err := program.JoinInterface(plusOp, InterfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join G+ importer")
	}

	return nil
}
