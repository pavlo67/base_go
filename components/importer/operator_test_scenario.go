package importer

import (
	"log"
	"testing"

	"github.com/pavlo67/data/components/ns"

	"github.com/stretchr/testify/require"
)

type ImporterTestCase struct {
	Operator Operator
	Source   ns.URN
}

func TestImporterWithCases(t *testing.T, testCases []ImporterTestCase) {
	for _, tc := range testCases {
		series, err := tc.Operator.Get(tc.Source, nil)
		require.NoError(t, err)
		require.NotNil(t, series)
		require.True(t, len(series.Records) > 0)

		for _, item := range series.Records {
			log.Printf("%#v", item)
		}
	}
}
