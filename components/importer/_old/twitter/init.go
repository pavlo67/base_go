package twitterimporter

import (
	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/config"
	"github.com/pavlo67/punctum/basis/program"
	"github.com/pavlo67/punctum/interfaces/starter"
	"github.com/pkg/errors"
)

// Starter ...
func Starter() starter.Operator {
	return &twitterComponent{}
}

type twitterComponent struct {
	key         string
	keySecret   string
	token       string
	tokenSecret string
}

const InterfaceKey = "importer.twitter"

func (fl *twitterComponent) Name() string {
	return InterfaceKey
}

func (tc *twitterComponent) Check(conf config.PunctumConfig, componentKeys config.ComponentKeys, indexPath string) ([]starter.Info, error) {

	var errs []error
	tc.key, errs = conf.Twitter("twitter_key", errs)
	tc.keySecret, errs = conf.Twitter("twitter_secret", errs)
	tc.token, errs = conf.Twitter("twitter_token", errs)
	tc.tokenSecret, errs = conf.Twitter("twitter_token_secret", errs)

	return nil, basis.JoinErrors(errs)
}

func (tc *twitterComponent) Setup(conf config.PunctumConfig, componentKeys config.ComponentKeys, indexPath string, data map[string]string) error {
	return nil
}

func (tc *twitterComponent) Init() error {
	twitterOp := &Twitter{
		Key:         tc.key,
		KeySecret:   tc.keySecret,
		Token:       tc.token,
		TokenSecret: tc.tokenSecret,
	}

	err := program.JoinInterface(twitterOp, InterfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join twitter importer")
	}

	return nil
}
