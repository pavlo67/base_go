package ns

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

type TestCase struct {
	URN          URN
	URNExpected  URN
	PathExpected string
	IsNull       bool
}

func TestIdentity(t *testing.T) {
	testCases := []TestCase{
		{"", "", "", true},
		{"abc", "abc", "", false},
		{"abc/", "abc", "/", false},
		{"/abc", "/abc", "/abc", false},
		{"dumaj.org.ua/abc/123#11", "dumaj.org.ua/abc/123#11", "/abc/123", false},
		{"dumaj.org.ua//123", "dumaj.org.ua//123", "//123", false},
		{"dumaj.org.ua/a/b/c/d/123", "dumaj.org.ua/a/b/c/d/123", "/a/b/c/d/123", false},
		{"dumaj.org.ua/a/b/c/d/123#", "dumaj.org.ua/a/b/c/d/123", "/a/b/c/d/123", false},
		{"dumaj.org.ua/a/b/c/d/123####abcd", "dumaj.org.ua/a/b/c/d/123####abcd", "/a/b/c/d/123", false},
	}

	for _, tc := range testCases {
		url := tc.URN.Item()

		if tc.IsNull {
			require.Nil(t, url)
		} else {
			require.NotNil(t, url)
			log.Printf("%s --> %#v", tc.URN, url)
			require.Equal(t, tc.PathExpected, url.Path)
			require.Equal(t, tc.URNExpected, url.URN())
		}
	}
}
