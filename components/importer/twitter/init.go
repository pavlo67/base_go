package twitterimporter

import (
	"github.com/pavlo67/punctum/starter"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pavlo67/punctum/starter/joiner"
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

const InterfaceKey joiner.InterfaceKey = "importer.twitter"

func (fl *twitterComponent) Name() string {
	return string(InterfaceKey)
}

func (tc *twitterComponent) Check(conf config.Config, indexPath string) ([]joiner.Info, error) {

	confTwitter, errs := conf.Google("", nil)
	if errs != nil {
		return nil, errs.Err()
	}
	tc.key = confTwitter["twitter_key"]
	tc.keySecret = confTwitter["twitter_secret"]
	tc.token = confTwitter["twitter_token"]
	tc.tokenSecret = confTwitter["twitter_token_secret"]

	return nil, errs.Err()
}

func (tc *twitterComponent) Setup(conf config.Config, indexPath string, data map[string]string) error {
	return nil
}

func (tc *twitterComponent) Init() error {
	twitterOp := &Twitter{
		Key:         tc.key,
		KeySecret:   tc.keySecret,
		Token:       tc.token,
		TokenSecret: tc.tokenSecret,
	}

	err := joiner.JoinInterface(twitterOp, InterfaceKey)
	if err != nil {
		return errors.Wrap(err, "can't join twitter importer")
	}

	return nil
}
