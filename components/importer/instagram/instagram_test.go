package instagramimporter

import (
	"log"
	"os"
	"testing"

	"github.com/pkg/errors"

	"github.com/pavlo67/punctum/basis/filelib"
	"github.com/pavlo67/punctum/starter/config"
)

func TestMain(m *testing.M) {
	if err := os.Setenv("ENV", "test"); err != nil {
		log.Fatalln("No test environment!!!")
	}
	os.Exit(m.Run())
}

func TestInstagram(t *testing.T) {
	//t.Skip()
	conf, err := config.Get(filelib.CurrentPath() + "../../../../punctum/cfg.json5")
	if err != nil {
		log.Fatal(err)
	}
	if conf == nil {
		log.Fatal(errors.New("no config data"))
	}

	confInstagram, errs := conf.Instagram("collector", nil)
	if errs != nil {
		t.Fatal(errs)
	}

	id := confInstagram["client_id"]
	secret := confInstagram["client_secret"]
	token := confInstagram["access_token"]

	var testCases = []importer_test.ImporterTestCase{

		{
			Operator: &Instagram{
				ID:     id,
				Secret: secret,
				Token:  token,
			},
			Fount: "https://www.instagram.com/onoff69/?hl=ru",
		},
	}

	importer_test.TestImporterWithCases(t, testCases)
}
