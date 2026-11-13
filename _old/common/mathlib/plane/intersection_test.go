package plane

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type IntersectionTestCase struct {
	S1, S2       Segment
	Intersection *Point2
}

func TestSegmentsIntersection(t *testing.T) {
	testCases := []IntersectionTestCase{
		{Segment{Point2{1, 0}, Point2{1, 1}},
			Segment{Point2{0, 1.5}, Point2{2, 1.5}}, nil},
		{Segment{Point2{1, 2}, Point2{1, 1}},
			Segment{Point2{1.25, 1.5}, Point2{2, 1.5}}, nil},

		{Segment{Point2{1, 1}, Point2{1, 2}},
			Segment{Point2{3, 0}, Point2{2, 1}}, nil},

		{Segment{Point2{2, 2}, Point2{-1, -1}},
			Segment{Point2{-3, 1}, Point2{3, -1}}, &Point2{0, 0}},
		{Segment{Point2{2, 2}, Point2{3, 3}},
			Segment{Point2{-3, 1}, Point2{3, -1}}, nil},
		{Segment{Point2{3, 3}, Point2{2, 2}},
			Segment{Point2{-3, 1}, Point2{3, -1}}, nil},
		{Segment{Point2{0, 0}, Point2{1, 1}},
			Segment{Point2{1, 1}, Point2{2, 2}}, nil},

		{Segment{Point2{0, 0}, Point2{0, 1}},
			Segment{Point2{0, 0.1}, Point2{0, 2}}, &Point2{0, 0.1}},

		{Segment{Point2{0, 0}, Point2{1, 0}},
			Segment{Point2{-2, 0}, Point2{-1, 0}}, nil},
		{Segment{Point2{1, 1}, Point2{1, 2}},
			Segment{Point2{1, 1}, Point2{2, 1}}, &Point2{1, 1}},
		{Segment{Point2{1, 1}, Point2{1, 2}},
			Segment{Point2{2, 1}, Point2{2, 2}}, nil},
		{Segment{Point2{1, 1}, Point2{1, 2}},
			Segment{Point2{0, 1.5}, Point2{2, 1.5}}, &Point2{1, 1.5}},
	}

	for i, testCase := range testCases {
		t.Log(i)
		// t.Logf("%#v", testCase)
		intersection := testCase.S1.Intersection(testCase.S2)
		if testCase.Intersection == nil {
			require.Nil(t, intersection)
		} else {
			require.NotNil(t, intersection)
			require.Equal(t, *testCase.Intersection, *intersection)
		}
	}
}

func TestSegmentsIntersection1(t *testing.T) {
	tests := []struct {
		name string
		s    Segment
		s1   Segment
		want *Point2
	}{
		{"", Segment{{1, 0}, {1, 1}}, Segment{{0, 1.5}, {2, 1.5}}, nil},
		{"", Segment{{1, 2}, {1, 1}}, Segment{{1.25, 1.5}, {2, 1.5}}, nil},
		{"", Segment{{1, 1}, {1, 2}}, Segment{{3, 0}, {2, 1}}, nil},
		{"", Segment{{2, 2}, {-1, -1}}, Segment{{-3, 1}, {3, -1}}, &Point2{0, 0}},
		{"", Segment{{2, 2}, {3, 3}}, Segment{{-3, 1}, {3, -1}}, nil},
		{"", Segment{{3, 3}, {2, 2}}, Segment{{-3, 1}, {3, -1}}, nil},
		{"", Segment{{0, 0}, {1, 1}}, Segment{{1, 1}, {2, 2}}, nil},
		{"", Segment{{0, 0}, {0, 1}}, Segment{{0, 0.1}, {0, 2}}, &Point2{0, 0.1}},
		{"", Segment{{0, 0}, {1, 0}}, Segment{{-2, 0}, {-1, 0}}, nil},
		{"", Segment{{1, 1}, {1, 2}}, Segment{{1, 1}, {2, 1}}, &Point2{1, 1}},
		{"", Segment{{1, 1}, {1, 2}}, Segment{{2, 1}, {2, 2}}, nil},
		{"", Segment{{1, 1}, {1, 2}}, Segment{{0, 1.5}, {2, 1.5}}, &Point2{1, 1.5}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPCross := tt.s.Intersection(tt.s1)
			if !reflect.DeepEqual(gotPCross, tt.want) {
				t.Errorf("SegmentsIntersection() gotPCross = %v, want %v", gotPCross, tt.want)
			}
		})
	}
}

func TestLinesIntersection(t *testing.T) {
	testCases := []IntersectionTestCase{
		{Segment{Point2{1, 1}, Point2{1, 2}},
			Segment{Point2{1, 1}, Point2{2, 1}}, &Point2{1, 1}},

		{Segment{Point2{1, 1}, Point2{1, 2}},
			Segment{Point2{2, 1}, Point2{2, 2}}, nil},

		{Segment{Point2{1, 1}, Point2{1, 2}},
			Segment{Point2{0, 1.5}, Point2{2, 1.5}}, &Point2{1, 1.5}},

		{Segment{Point2{1, 2}, Point2{1, 1}},
			Segment{Point2{1.25, 1.5}, Point2{2, 1.5}}, &Point2{1, 1.5}},

		{Segment{Point2{1, 1}, Point2{1, 2}},
			Segment{Point2{3, 0}, Point2{2, 1}}, &Point2{1, 2}},

		{Segment{Point2{2, 2}, Point2{-1, -1}},
			Segment{Point2{-3, 1}, Point2{3, -1}}, &Point2{0, 0}},

		{Segment{Point2{2, 2}, Point2{3, 3}},
			Segment{Point2{-3, 1}, Point2{3, -1}}, &Point2{0, 0}},

		{Segment{Point2{3, 3}, Point2{2, 2}},
			Segment{Point2{-3, 1}, Point2{3, -1}}, &Point2{0, 0}},

		{Segment{Point2{0, 0}, Point2{1, 1}},
			Segment{Point2{1, 1}, Point2{2, 2}}, nil},

		{Segment{Point2{0, 0}, Point2{0, 1}},
			Segment{Point2{0, 0.1}, Point2{0, 2}}, nil},

		{Segment{Point2{0, 0}, Point2{1, 0}},
			Segment{Point2{-2, 0}, Point2{-1, 0}}, nil},
	}

	for i, testCase := range testCases {
		t.Log(i)
		intersection := testCase.S1.LinesIntersection(testCase.S2)
		if testCase.Intersection == nil {
			require.Nil(t, intersection)
		} else {
			require.NotNil(t, intersection)
			require.Equal(t, *testCase.Intersection, *intersection)
		}
	}
}
