package test

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/pavlo67/punctum/interfaces/importer/html"
	"github.com/pavlo67/punctum/interfaces/importer/test_scenario"
)

func TestMain(m *testing.M) {
	if _, ok := os.LookupEnv("TEST_ENVIRONMENT"); !ok {
		log.Fatalln("No test environment!!!")
	}
	os.Exit(m.Run())
}

func setParams() []importer_test.ImporterTestCase {

	var p = htmlimporter.ImportParams{
		AcceptableTags:       []string{"html", "body", "title", "div"},
		ImportSeparateRegexp: "</div>",
	}
	pJSON, _ := json.Marshal(p)

	var testCases = []importer_test.ImporterTestCase{
		{
			Operator: &htmlimporter.ImporterHTML{},
			Fount:    "https://www.unian.ua/world/10044620-druzhina-trampa-molodshogo-podala-na-rozluchennya.html",
			DBKey:    string(pJSON),
		},
	}
	return testCases
}

func TestHTML(t *testing.T) {
	importer_test.TestImporterWithCases(t, setParams())
}
