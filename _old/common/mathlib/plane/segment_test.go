package plane

import (
	"reflect"
	"testing"
)

func TestSegmentGoOutCircle(t *testing.T) {
	tests := []struct {
		name string
		s    Segment
		p    Point2
		r    float64
		want *Point2
	}{
		{
			name: "",
			s:    Segment{{3, 0}, {3, 10}},
			p:    Point2{0, 0},
			r:    5,
			want: &Point2{3, 4},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.GoOutCircle(tt.p, tt.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GoOutCircle() = %v, want %v", got, tt.want)
			}
		})
	}
}
