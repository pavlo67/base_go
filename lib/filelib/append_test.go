package filelib

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

const testFileName = "append.test"

func TestAppendFile(t *testing.T) {
	os.Remove(testFileName)

	var testData string

	for i := 0; i < 5; i++ {
		a := strconv.Itoa(i)
		testData += a

		err := AppendFile(testFileName, []byte(a))
		require.NoError(t, err)
	}

	data, err := os.ReadFile(testFileName)
	require.NoError(t, err)
	require.Equal(t, testData, string(data))

}
