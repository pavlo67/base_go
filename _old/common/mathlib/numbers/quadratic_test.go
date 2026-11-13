package numbers

import (
	"reflect"
	"testing"
)

func TestQuadraticEquation(t *testing.T) {
	tests := []struct {
		name string
		a    float64
		b    float64
		c    float64
		want *[2]float64
	}{
		{
			name: "",
			a:    1,
			b:    0,
			c:    -1,
			want: &[2]float64{1, -1},
		},
		{
			name: "",
			a:    1,
			b:    3,
			c:    -4,
			want: &[2]float64{1, -4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := QuadraticEquation(tt.a, tt.b, tt.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QuadraticEquation() = %v, want %v", got, tt.want)
			}
		})
	}
}
