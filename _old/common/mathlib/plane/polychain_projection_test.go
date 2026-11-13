package plane

import (
	"math"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/mathlib"
)

func TestProjectionsOnPolyChainProbe(t *testing.T) {
	polyChain := PolyChain{{221.29971138416823, 290.62290413201606}, {238, 268.5}, {262, 252}}
	point := Point2{237, 266}

	got := point.ProjectionsOnPolyChain(polyChain, 10)

	t.Logf("%#v", got)
}

func TestProjectionsOnPolyChain(t *testing.T) {
	tests := []struct {
		name        string
		polyChain   PolyChain
		p           Point2
		distanceMax float64
		expected    []ProjectionOnPolyChainDirected
	}{
		{
			name:        "",
			polyChain:   PolyChain{{0, 2}, {1, 1}, {2, 2}}, // , {3, 1}
			p:           Point2{1, 0},
			distanceMax: 10,
			expected: []ProjectionOnPolyChainDirected{
				{Distance: 1, ProjectionOnPolyChain: ProjectionOnPolyChain{N: 1, Point2: Point2{1, 1}}}},
		},
		{
			name:        "",
			polyChain:   PolyChain{{0, 2}, {1, 1}, {2, 2}, {3, 1}}, //
			p:           Point2{1, 0},
			distanceMax: 10,
			expected: []ProjectionOnPolyChainDirected{
				{Distance: 1, ProjectionOnPolyChain: ProjectionOnPolyChain{N: 1, Point2: Point2{1, 1}}},
				{Distance: 3 * math.Sqrt(2) / 2,
					ProjectionOnPolyChain: ProjectionOnPolyChain{N: 2, Position: math.Sqrt(2) / 2, Point2: Point2{2.5, 1.5}}},
			},
		},
		{
			name:        "",
			polyChain:   PolyChain{{X: 567.758571909734, Y: 327.33667901650387}, {X: 588, Y: 381}}, //
			p:           Point2{X: 574, Y: 355},
			distanceMax: 3,
			expected:    []ProjectionOnPolyChainDirected{},
		},
		{
			name:        "",
			polyChain:   PolyChain{{X: 416.5, Y: 420}, {X: 422, Y: 413}, {X: 414.5, Y: 402}}, //
			p:           Point2{X: 427, Y: 412},
			distanceMax: 3,
			expected:    []ProjectionOnPolyChainDirected{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.ProjectionsOnPolyChain(tt.polyChain, tt.distanceMax)
			CheckProjections(t, tt.expected, got)
		})
	}
}

// PolyChain{Point2{X:556, Y:355}, Point2{X:559, Y:355}, Point2{X:562, Y:355}, Point2{X:565, Y:355}, Point2{X:568,
//	Y:355}, Point2{X:571, Y:355}, Point2} / 3)

func CheckProjections(t *testing.T, expected, got []ProjectionOnPolyChainDirected) {
	require.Equalf(t, len(expected), len(got), "expected: %v, got: %#v", expected, got)
	for i, e := range expected {
		g := got[i]
		require.Equalf(t, e.N, g.N, "#d: %v vs %v", i, e, g)
		require.Truef(t, math.Abs(e.Position-g.Position) <= mathlib.EPS, "#d: %v vs %v", i, e, g)
		require.Truef(t, math.Abs(e.Distance-g.Distance) <= mathlib.EPS, "#d: %v vs %v", i, e, g)
		require.Truef(t, math.Abs(e.Angle-g.Angle) <= mathlib.EPS, "#d: %v vs %v", i, e, g)
		require.Truef(t, math.Abs(e.X-g.X) <= mathlib.EPS, "#d: %v vs %v", i, e, g)
		require.Truef(t, math.Abs(e.Y-g.Y) <= mathlib.EPS, "#d: %v vs %v", i, e, g)

	}
}

func TestCutWithProjections(t *testing.T) {
	tests := []struct {
		name string
		pCh  PolyChain
		pr0  ProjectionOnPolyChain
		pr1  ProjectionOnPolyChain
		want PolyChain
	}{
		{
			name: "",
			pCh:  PolyChain{{X: 2.3, Y: 458}, {X: 91, Y: 427.5}},
			pr0:  ProjectionOnPolyChain{N: 1, Position: 0, Point2: Point2{X: 91, Y: 427.5}},
			pr1:  ProjectionOnPolyChain{N: 0, Position: 9.7, Point2: Point2{X: 11.5, Y: 455}},
			want: PolyChain{{X: 91, Y: 427.5}, {11.5, 455}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pCh.CutWithProjections(tt.pr0, tt.pr1); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CutWithProjections() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDistanceToSegment(t *testing.T) {

	tests := []struct {
		name    string
		p       Point2
		segment Segment
		want    float64
	}{
		{
			name:    "",
			p:       Point2{203.77788799006555, 564.5811856102348},
			segment: Segment{{230.79193808962214, 591.2694925293472}, {145.99999999999994, 507.50000000000006}},
			want:    0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := tt.p.DistanceToSegment(tt.segment)
			if got != tt.want {
				t.Errorf("DistanceToSegment() got = %v, want %v", got, tt.want)
			}
		})
	}
}
