package serialization

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type common struct {
	name    string
	prefix  string
	indent  string
	want    string
	wantErr bool
}

type testCase[T any] struct {
	common
	list []T
}

func TestJSONList(t *testing.T) {
	testInt := testCase[int]{
		list: []int{1, 2, 3, 4},
		common: common{
			prefix:  "\n",
			indent:  "  ",
			want:    "[\n  1,\n  2,\n  3,\n  4\n]",
			wantErr: false,
		},
	}
	CheckOneJSONList(t, testInt.common, testInt.list)

	testStr := testCase[string]{
		list: []string{"1", "2", "3", "4"},
		common: common{
			prefix:  "\n",
			indent:  "  ",
			want:    "[\n  \"1\",\n  \"2\",\n  \"3\",\n  \"4\"\n]",
			wantErr: false,
		},
	}
	CheckOneJSONList(t, testStr.common, testStr.list)

	testAny := testCase[interface{}]{
		list: []interface{}{"1", "2", 3, []int{4, 5, 6}},
		common: common{
			prefix:  "\n",
			indent:  "  ",
			want:    "[\n  \"1\",\n  \"2\",\n  3,\n  [4,5,6]\n]",
			wantErr: false,
		},
	}
	CheckOneJSONList(t, testAny.common, testAny.list)

}

func CheckOneJSONList[T any](t *testing.T, tt common, list []T) {
	t.Run(tt.name, func(t *testing.T) {
		got, err := JSONList(list, tt.prefix, tt.indent)
		require.Equalf(t, tt.wantErr, err != nil, "tt.wantErr = %t, err = %s", tt.wantErr, err)
		require.Equal(t, tt.want, string(got))
	})

}
