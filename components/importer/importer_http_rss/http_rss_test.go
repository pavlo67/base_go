package importer_http_rss

import (
	"testing"

	"github.com/pavlo67/common/common/instruments/importer"
)

var testCases = []importer.ImporterTestCase{
	{
		Operator: &rss{},
		Source:   "https://rss.unian.net/site/news_ukr.rss",
	},
}

func TestRSS(t *testing.T) {
	importer.TestImporterWithCases(t, testCases)
}
