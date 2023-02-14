package htmlimporter

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"github.com/pavlo67/punctum/collector/importer/test_scenario"
)

func TestMain(m *testing.M) {
	if err := os.Setenv("ENV", "test"); err != nil {
		log.Fatalln("No test environment!!!")
	}
	os.Exit(m.Run())
}

func setParams() []importer_test.ImporterTestCase {

	var p = ImportParams{
		AcceptableTags:       []string{"html", "body", "title", "div"},
		ImportSeparateRegexp: "</div>",
	}
	pJSON, _ := json.Marshal(p)

	var testCases = []importer_test.ImporterTestCase{
		{
			Operator: &ImporterHTML{},
			Fount:    "https://www.unian.ua/world/10044620-druzhina-trampa-molodshogo-podala-na-rozluchennya.html",
			DBKey:    string(pJSON),
		},
	}
	return testCases
}

func TestHTML(t *testing.T) {
	importer_test.TestImporterWithCases(t, setParams())
}
