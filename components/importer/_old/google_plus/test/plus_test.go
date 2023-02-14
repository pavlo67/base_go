package test

import (
	"errors"
	"log"
	"os"
	"testing"

	"github.com/pavlo67/punctum/basis"
	"github.com/pavlo67/punctum/basis/config"

	"github.com/pavlo67/punctum/interfaces/importer/google_plus"
	"github.com/pavlo67/punctum/interfaces/importer/test_scenario"
)

func TestMain(m *testing.M) {
	if _, ok := os.LookupEnv("TEST_ENVIRONMENT"); !ok {
		log.Fatalln("No test environment!!!")
	}
	os.Exit(m.Run())
}

func setParams() []importer_test.ImporterTestCase {

	conf, err := config.Get(basis.CurrentPath() + "../../../../../punctum/cfg.json5")
	if err != nil {
		log.Fatal(err)
	}
	if conf == nil {
		log.Fatal(errors.New("no config data"))
	}

	var errs []error
	//var id, secret, key, path string
	key, errs := conf.Google("api_key", errs)
	//id, errs = conf.Google("user_id", errs)
	//secret, errs = conf.Google("secret", errs)
	//path, errs = conf.Google("path_to_json", errs)
	var testCases = []importer_test.ImporterTestCase{
		{
			Operator: &google_plus.PLUS{
				ApiKey: key,
				//ApiID: id,
				//ApiSecret: secret,
				//PathToJSON: path,

			},
			Fount: "https://www.googleapis.com/plus/v1/people/103228082707112449686/activities/public",
			DBKey: "",
		},
	}
	return testCases
}

func TestHTML(t *testing.T) {
	importer_test.TestImporterWithCases(t, setParams())
}
