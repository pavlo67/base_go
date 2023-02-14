package test

import (
	"log"
	"os"
	"testing"

	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/config"

	"github.com/pavlo67/punctum/interfaces/importer/instagram"
	"github.com/pavlo67/punctum/interfaces/importer/test_scenario"
)

func TestMain(m *testing.M) {
	if _, ok := os.LookupEnv("TEST_ENVIRONMENT"); !ok {
		log.Fatalln("No test environment!!!")
	}
	os.Exit(m.Run())
}

func TestInstagram(t *testing.T) {
	//t.Skip()
	conf, err := config.Get(basis.CurrentPath() + "../../../../../punctum/cfg.json5")
	if err != nil {
		log.Fatal(err)
	}
	if conf == nil {
		log.Fatal(errors.New("no config data"))
	}

	var errs []error
	var id, secret, token string
	id, errs = conf.Instagram("client_id", errs)
	secret, errs = conf.Instagram("client_secret", errs)
	token, errs = conf.Instagram("access_token", errs)

	var testCases = []importer_test.ImporterTestCase{

		{
			Operator: &instagramimporter.Instagram{
				ID:     id,
				Secret: secret,
				Token:  token,
			},
			Fount: "https://www.instagram.com/onoff69/?hl=ru",
		},
	}

	importer_test.TestImporterWithCases(t, testCases)
}
