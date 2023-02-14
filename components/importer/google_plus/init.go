package google_plus

import (
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"
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

	conf   config.Config
	params map[string]string
}

const InterfaceKey joiner.InterfaceKey = "importer.google_plus"

func (fl *plusComponent) Name() string {
	return string(InterfaceKey)
}

func (pc *plusComponent) Check(conf config.Config, indexPath string) ([]joiner.Info, error) {

	confGoogle, errs := conf.Google("", nil)
	if errs != nil {
		return nil, errs.Err()
	}
	pc.apiKey = confGoogle["api_key"]

	pc.conf = conf
	return nil, nil
}

func (pc *plusComponent) Setup(conf config.Config, indexPath string, data map[string]string) error {
	pc.conf = conf
	return nil
}

func (pc *plusComponent) Init() error {
	plusOp := &googlePlus{
		ApiKey: pc.apiKey,
		//ApiID:      pc.apiID,
		//ApiSecret:  pc.apiSecret,
		//PathToJSON: pc.pathToJSON,

	}

	err := joiner.JoinInterface(plusOp, InterfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join G+ importer")
	}

	return nil
}
