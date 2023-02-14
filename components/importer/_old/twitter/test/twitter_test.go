package test

import (
	"log"
	"os"
	"testing"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/config"
	"github.com/pavlo67/punctum/interfaces/importer/test_scenario"
	"github.com/pavlo67/punctum/interfaces/importer/twitter"
	"github.com/pkg/errors"
)

func TestMain(m *testing.M) {
	if _, ok := os.LookupEnv("TEST_ENVIRONMENT"); !ok {
		log.Fatalln("No test environment!!!")
	}
	os.Exit(m.Run())
}

func TestTwitter(t *testing.T) {
	//t.Skip()
	conf, err := config.Get(basis.CurrentPath() + "../../../../../punctum/cfg.json5")
	if err != nil {
		log.Fatal(err)
	}
	if conf == nil {
		log.Fatal(errors.New("no config data"))
	}

	var errs []error
	var key, keySecret, token, tokenSecret string
	key, errs = conf.Twitter("twitter_key", errs)
	keySecret, errs = conf.Twitter("twitter_secret", errs)
	token, errs = conf.Twitter("twitter_token", errs)
	tokenSecret, errs = conf.Twitter("twitter_token_secret", errs)

	var testCases = []importer_test.ImporterTestCase{

		{
			Operator: &twitterimporter.Twitter{
				Key:         key,
				KeySecret:   keySecret,
				Token:       token,
				TokenSecret: tokenSecret,
			},
			//Fount:    "https://mobile.twitter.com/realdonaldtrump",
			Fount: "https://mobile.twitter.com/dmitrosel007",
		},
	}

	importer_test.TestImporterWithCases(t, testCases)
}
