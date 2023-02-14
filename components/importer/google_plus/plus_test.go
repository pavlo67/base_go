package google_plus

import (
	"errors"
	"log"
	"testing"

	"github.com/pavlo67/punctum/basis/filelib"
	"github.com/pavlo67/punctum/starter/config"
)

func setParams(t *testing.T) []importer_test.ImporterTestCase {

	conf, err := config.Get(filelib.CurrentPath() + "../../../../punctum/cfg.json5")
	if err != nil {
		log.Fatal(err)
	}
	if conf == nil {
		log.Fatal(errors.New("no config data"))
	}

	confGoogle, errs := conf.Google("collector", nil)
	if errs != nil {
		t.Fatal(errs)
	}

	key := confGoogle["api_key"]

	var testCases = []importer_test.ImporterTestCase{
		{
			Operator: &googlePlus{
				ApiKey: key,
				//ApiID: id,
				//ApiSecret: secret,
				//PathToJSON: path,

			},
			Fount: "https://www.googleapis.com/plus/flow_v1/people/103228082707112449686/activities/public",
			DBKey: "",
		},
	}
	return testCases
}

func TestHTML(t *testing.T) {
	importer_test.TestImporterWithCases(t, setParams(t))
}
