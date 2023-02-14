package twitterimporter

import (
	"log"
	"os"
	"testing"

	"github.com/pavlo67/punctum/basis/filelib"
	"github.com/pavlo67/punctum/starter/config"
	"github.com/pkg/errors"
)

func TestMain(m *testing.M) {
	if err := os.Setenv("ENV", "test"); err != nil {
		log.Fatalln("No test environment!!!")
	}
	os.Exit(m.Run())
}

func TestTwitter(t *testing.T) {
	//t.Skip()
	conf, err := config.Get(filelib.CurrentPath() + "../../../../punctum/cfg.json5")
	if err != nil {
		log.Fatal(err)
	}
	if conf == nil {
		log.Fatal(errors.New("no config data"))
	}

	confTwitter, errs := conf.Twitter("", nil)
	if errs != nil {
		t.Fatal(errs)
	}

	key := confTwitter["key"]
	keySecret := confTwitter["secret"]
	token := confTwitter["token"]
	tokenSecret := confTwitter["token_secret"]

	var testCases = []importer_test.ImporterTestCase{

		{
			Operator: &Twitter{
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
