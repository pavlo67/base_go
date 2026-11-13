package geolib

import (
	"math"
	"testing"

	"github.com/pavlo67/common/common/mathlib"
)

func TestDegrees_DMS(t *testing.T) {
	tests := []struct {
		name    string
		degrees Degrees
		want    DMS
	}{
		{
			name:    "",
			degrees: 10.5,
			want:    DMS{10, 30, 0},
		},
		{
			name:    "",
			degrees: 10.55,
			want:    DMS{10, 33, 0},
		},
		{
			name:    "",
			degrees: 10.555,
			want:    DMS{10, 33, 18},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.degrees.DMS(); !(got.D == tt.want.D && got.M == tt.want.M && math.Abs(got.S-tt.want.S) <= mathlib.EPS) {
				t.Errorf("DMS() = %v, want %v", got, tt.want)
			}
		})
	}
}
