package sets

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCorrespondence_AlgClosure(t *testing.T) {
	type testCase[L comparable, R comparable] struct {
		name string
		corr Correspondence[L, R]
		want CorrespondenceClosed[L, R]
	}
	tests := []testCase[int, string]{
		{
			name: "1",
			corr: nil,
			want: nil,
		},
		{
			name: "1",
			corr: Correspondence[int, string]{{1, "a"}},
			want: CorrespondenceClosed[int, string]{{[]int{1}, []string{"a"}}},
		},
		{
			name: "1",
			corr: Correspondence[int, string]{{1, "a"}, {2, "a"}, {1, "b"}, {3, "b"}, {4, "c"}},
			want: CorrespondenceClosed[int, string]{{[]int{1, 2, 3}, []string{"a", "b"}}, {[]int{4}, []string{"c"}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.corr.AlgClosure()
			require.Equal(t, len(tt.want), len(got))

		GOT:
			for _, g := range got {
				for _, w := range tt.want {
					if !(Equal(g.L, w.L) && Equal(g.R, w.R)) {
						continue
					}
					continue GOT
				}
				t.Errorf("AlgClosure() = %v, want %v, error resulting: %v", got, tt.want, g)
			}
		})
	}
}
