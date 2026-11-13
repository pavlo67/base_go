package ziplib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestZipUnzip(t *testing.T) {
	src := "aaa.zip"
	filename := "aaa.txt"
	data := "aaa"

	_, err := ZipFiles(src, []ToZip{{[]byte(data), filename}}, 0644)
	require.NoError(t, err)

	dataUnzipped, err := UnzipFile(src, filename)
	require.NoError(t, err)
	require.Equal(t, data, string(dataUnzipped))
}
