package test

import (
	"log"
	"os"
	"testing"

	"github.com/pavlo67/punctum/interfaces/importer/rss"
	"github.com/pavlo67/punctum/interfaces/importer/test_scenario"
)

func TestMain(m *testing.M) {
	if _, ok := os.LookupEnv("TEST_ENVIRONMENT"); !ok {
		log.Fatalln("No test environment!!!")
	}
	os.Exit(m.Run())
}

var testCases = []importer_test.ImporterTestCase{
	{
		Operator: &rss.RSS{},
		Fount:    "https://rss.unian.net/site/news_ukr.rss",
	},
}

func TestRSS(t *testing.T) {
	importer_test.TestImporterWithCases(t, testCases)
}
