package strlib

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSpaces(t *testing.T) {
	original := "   a   b c    d "
	expected := "a b c d"
	require.Equal(t, expected, ReSpaces.ReplaceAllString(ReSpacesFin.ReplaceAllString(original, ""), " "))
}
