package plane

import (
	"math"
	"reflect"
	"testing"

	"github.com/pavlo67/common/common/mathlib"
)

func TestPolyChainAveragedProbe(t *testing.T) {
	pCh0 := PolyChain{
		{342, 162.5}, {364, 207.5},
	}
	pCh1 := PolyChain{
		{335, 151}, {406.5, 302}, {403.7629562043795, 375.0206204379562},
	}

	//[{341.3650461698297 162.800656946074} {363.29489523856535 207.82993665827638}]
	//[{334.4224510771934 151.46712409917882} {337.37437836467245 156.29806754937331}]
	//[{341.3650461698297 162.800656946074} {406.5 302} {403.7629562043795 375.0206204379562}]

	ok, pCh0Averaged, pCh1Averaged := pCh0.AverageWithAnother(pCh1, 9.870530984139577, false)

	t.Log(ok, "\n", pCh0Averaged, "\n", pCh1Averaged)
}

func TestAveragePolyChains(t *testing.T) {
	tests := []struct {
		name                   string
		polyChain0             PolyChain
		polyChain1             PolyChain
		minDistance            float64
		wantOk                 bool
		wantPolyChain0Averaged PolyChain
		wantPolyChain1Rest     []PolyChain
		connectEnds            bool
	}{
		{
			name:                   "",
			polyChain0:             PolyChain{{472, 513}, {648, 197}},
			polyChain1:             PolyChain{{673, 13}, {648, 197}, {472, 513}},
			minDistance:            11,
			wantOk:                 true,
			wantPolyChain0Averaged: nil,
			wantPolyChain1Rest:     []PolyChain{{{472, 513}, {648, 197}, {673, 13}}},
			connectEnds:            true,
		},

		{
			name:                   "",
			polyChain0:             PolyChain{{0, 0}, {0, 2}, {0, 4}},
			polyChain1:             PolyChain{{0, 4}, {0, 6}},
			minDistance:            0,
			wantOk:                 true,
			wantPolyChain0Averaged: PolyChain{{0, 0}, {0, 2}, {0, 4}},
			wantPolyChain1Rest:     []PolyChain{{{0, 4}, {0, 6}}},
			connectEnds:            false,
		},
		{
			name:                   "",
			polyChain0:             PolyChain{{0, 0}, {0, 2}, {0, 4}},
			polyChain1:             PolyChain{{0, 3}, {0, 4}, {0, 6}},
			minDistance:            0,
			wantOk:                 true,
			wantPolyChain0Averaged: PolyChain{{0, 0}, {0, 2}, {0, 3}, {0, 4}},
			wantPolyChain1Rest:     []PolyChain{{{0, 4}, {0, 6}}},
			connectEnds:            false,
		},
		{
			name:                   "",
			polyChain0:             PolyChain{{0, 0}, {0.05, 2}, {0, 4}},
			polyChain1:             PolyChain{{0.1, 3}, {0.1, 4}, {0.1, 6}},
			minDistance:            0.1,
			wantOk:                 true,
			wantPolyChain0Averaged: PolyChain{{0, 0}, {0.05, 2}, {0.07500000000000001, 3}, {0.05, 4}},
			wantPolyChain1Rest:     []PolyChain{{{0.05, 4}, {0.1, 6}}},
			connectEnds:            false,
		},

		{
			name:                   "",
			polyChain0:             PolyChain{{0, 0}, {0, 2}, {0, 4}},
			polyChain1:             PolyChain{{0, 4}, {0, 6}},
			minDistance:            0,
			wantOk:                 true,
			wantPolyChain0Averaged: PolyChain{{0, 0}, {0, 2}, {0, 4}},
			wantPolyChain1Rest:     []PolyChain{{{0, 4}, {0, 6}}},
		},
		{
			name:                   "",
			polyChain0:             PolyChain{{0, 0}, {0, 2}, {0, 4}},
			polyChain1:             PolyChain{{0, 3}, {0, 4}, {0, 6}},
			minDistance:            0,
			wantOk:                 true,
			wantPolyChain0Averaged: PolyChain{{0, 0}, {0, 2}, {0, 3}, {0, 4}},
			wantPolyChain1Rest:     []PolyChain{{{0, 4}, {0, 6}}},
		},
		{
			name:                   "",
			polyChain0:             PolyChain{{0, 0}, {0.05, 2}, {0, 4}},
			polyChain1:             PolyChain{{0.1, 3}, {0.1, 4}, {0.1, 6}},
			minDistance:            0.1,
			wantOk:                 true,
			wantPolyChain0Averaged: PolyChain{{0, 0}, {0.05, 2}, {0.07500000000000001, 3}, {0.05, 4}},
			wantPolyChain1Rest:     []PolyChain{{{0.05, 4}, {0.1, 6}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOk, gotAveraged, gotRest := tt.polyChain0.AverageWithAnother(tt.polyChain1, tt.minDistance, tt.connectEnds)

			if gotOk != tt.wantOk {
				t.Errorf("AverageWithAnother() gotOk = %t, wantOk %t", gotOk, tt.wantOk)
			}
			if !reflect.DeepEqual(gotAveraged, tt.wantPolyChain0Averaged) {
				t.Errorf("AverageWithAnother() gotAveraged = %v, wantAveraged %v", gotAveraged, tt.wantPolyChain0Averaged)
			}
			if !reflect.DeepEqual(gotRest, tt.wantPolyChain1Rest) {
				t.Errorf("AverageWithAnother() gotRest = %v, wantRest %v", gotRest, tt.wantPolyChain1Rest)
			}
		})
	}
}

func ComparePolyChains(pCh0, pCh1 PolyChain) bool {
	if len(pCh0) != len(pCh1) {
		return false
	}
	for i, p := range pCh0 {
		if p != pCh1[i] {
			return false
		}
	}

	return true
}
func TestCutPolyChain(t *testing.T) {
	type args struct {
		polyChain PolyChain
		endI      int
		axis      Segment
	}

	tests := []struct {
		name string
		args args
		want []Point2
	}{
		{
			name: "",
			args: args{
				polyChain: []Point2{{-1, -1}, {1, 1}, {1, -1}},
				endI:      0,
				axis:      Segment{Point2{X: 0, Y: 2}, Point2{X: 1, Y: 2}},
			},
			want: []Point2{{-1, -1}, {1, 1}, {1, -1}},
		},

		{
			name: "",
			args: args{
				polyChain: []Point2{{-1, -1}, {1, 1}, {1, -1}},
				endI:      0,
				axis:      Segment{Point2{X: 0, Y: 1}, Point2{X: 1, Y: 1}},
			},
			want: []Point2{{-1, -1}, {1, 1}, {1, -1}},
		},

		{
			name: "",
			args: args{
				polyChain: []Point2{{-1, -1}, {1, 1}, {1, -1}},
				endI:      0,
				axis:      Segment{Point2{X: 0, Y: 0}, Point2{X: 1, Y: 0}},
			},

			want: []Point2{{-1, -1}, {0, 0}, {1, 0}, {1, -1}},
		},

		{
			name: "",
			args: args{
				polyChain: PolyChain{{-1, -1}, {1, 1}, {1, -1}},
				endI:      1,
				axis:      Segment{Point2{X: 0, Y: 0}, Point2{X: 1, Y: 0}},
			},

			want: []Point2{{1, 1}, {1, 0}, {0, 0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.polyChain.Cut(tt.args.endI, tt.args.axis); !ComparePolyChains(got, tt.want) {
				t.Errorf("CutPolyChain() = %v, wantPolyChain %v", got, tt.want)
			}
		})
	}
}

func TestApproximatePolyChain(t *testing.T) {
	tests := []struct {
		name         string
		pts2         PolyChain
		maxDeviation float64
		want         PolyChain
	}{
		{
			name:         "0",
			pts2:         nil,
			maxDeviation: 1,
			want:         nil,
		},
		{
			name: "1", pts2: PolyChain{{}}, maxDeviation: 1, want: PolyChain{{}},
		},
		{
			name: "2", pts2: PolyChain{{}, {1, 1}}, maxDeviation: 1, want: PolyChain{{}, {1, 1}},
		},
		{
			name: "3", pts2: PolyChain{{}, {0, 1}, {2, 2}}, maxDeviation: 0.8,
			want: PolyChain{{}, {2, 2}},
		},
		{
			name: "4", pts2: PolyChain{{}, {0, 1}, {2, 2}}, maxDeviation: 0.5,
			want: PolyChain{{}, {0, 1}, {2, 2}},
		},
		{
			name: "5", pts2: PolyChain{{}, {1, 1}, {2, 1}, {3, 1}, {4, 0}}, maxDeviation: 0.9,
			want: PolyChain{{}, {2, 1}, {4, 0}},
		},
		{
			name: "5", pts2: PolyChain{{}, {-1, 0}, {4, 0}}, maxDeviation: 0.9,
			want: PolyChain{{}, {-1, 0}, {4, 0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pts2.Approximate(tt.maxDeviation); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ApproximatePolyChain() = %v, wantDistance %v", got, tt.want)
			}
		})
	}
}

//func TestAveragePolyChains(t *testing.T) {
//	tests := []struct {
//		name                   string
//		polyChain0             PolyChain
//		polyChain1             PolyChain
//		minDistance            float64
//		wantOk                 bool
//		wantPolyChain0Averaged PolyChain
//		wantPolyChain1Rest     []PolyChain
//	}{
//		{
//			name:                   "",
//			polyChain0:             PolyChain{{0, 0}, {0, 2}, {0, 4}},
//			polyChain1:             PolyChain{{0, 4}, {0, 6}},
//			minDistance:            0,
//			wantOk:                 true,
//			wantPolyChain0Averaged: PolyChain{{0, 0}, {0, 2}, {0, 4}},
//			wantPolyChain1Rest:     []PolyChain{{{0, 4}, {0, 6}}},
//		},
//		{
//			name:                   "",
//			polyChain0:             PolyChain{{0, 0}, {0, 2}, {0, 4}},
//			polyChain1:             PolyChain{{0, 3}, {0, 4}, {0, 6}},
//			minDistance:            0,
//			wantOk:                 true,
//			wantPolyChain0Averaged: PolyChain{{0, 0}, {0, 2}, {0, 3}, {0, 4}},
//			wantPolyChain1Rest:     []PolyChain{{{0, 4}, {0, 6}}},
//		},
//		{
//			name:                   "",
//			polyChain0:             PolyChain{{0, 0}, {0.05, 2}, {0, 4}},
//			polyChain1:             PolyChain{{0.1, 3}, {0.1, 4}, {0.1, 6}},
//			minDistance:            0.1,
//			wantOk:                 true,
//			wantPolyChain0Averaged: PolyChain{{0, 0}, {0.05, 2}, {0.07500000000000001, 3}, {0.05, 4}},
//			wantPolyChain1Rest:     []PolyChain{{{0.05, 4}, {0.1, 6}}},
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			gotOk, gotAveraged, gotRest := AverageWithAnother(tt.polyChain0, tt.polyChain1, tt.minDistance)
//
//			if gotOk != tt.wantOk {
//				t.Errorf("AverageWithAnother() gotOk = %t, wantOk %t", gotOk, tt.wantOk)
//			}
//			if !reflect.DeepEqual(gotAveraged, tt.wantPolyChain0Averaged) {
//				t.Errorf("AverageWithAnother() gotAveraged = %v, wantAveraged %v", gotAveraged, tt.wantPolyChain0Averaged)
//			}
//			if !reflect.DeepEqual(gotRest, tt.wantPolyChain1Rest) {
//				t.Errorf("AverageWithAnother() gotRest = %v, wantRest %v", gotRest, tt.wantPolyChain1Rest)
//			}
//		})
//	}
//}

func TestDistanceToPolyChain(t *testing.T) {
	tests := []struct {
		name           string
		polyChain      PolyChain
		p              Point2
		wantDistance   float64
		wantProjection ProjectionOnPolyChain
	}{
		{
			name:           "",
			polyChain:      PolyChain{{0.1, 3}, {0.1, 4}, {0.1, 6}},
			p:              Point2{0, 4},
			wantDistance:   0.1,
			wantProjection: ProjectionOnPolyChain{N: 1, Position: 0, Point2: Point2{0.1, 4}},
		},
		{
			name:           "",
			polyChain:      PolyChain{{0, 4}, {0, 6}},
			p:              Point2{0, 0},
			wantDistance:   4,
			wantProjection: ProjectionOnPolyChain{N: 0, Position: 0, Point2: Point2{0, 4}},
		},
		{
			name:           "",
			polyChain:      PolyChain{{0, 4}, {0, 6}},
			p:              Point2{0, 2},
			wantDistance:   2,
			wantProjection: ProjectionOnPolyChain{N: 0, Position: 0, Point2: Point2{0, 4}},
		},
		{
			name:           "",
			polyChain:      PolyChain{{0, 4}, {0, 6}},
			p:              Point2{0, 4},
			wantDistance:   0,
			wantProjection: ProjectionOnPolyChain{N: 0, Position: 0, Point2: Point2{0, 4}},
		},
		{
			name:           "",
			polyChain:      PolyChain{{0, 0}, {0, 2}, {0, 4}},
			p:              Point2{0, 4},
			wantDistance:   0,
			wantProjection: ProjectionOnPolyChain{N: 2, Position: 0, Point2: Point2{0, 4}},
		},
		{
			name:           "",
			polyChain:      PolyChain{{0, 0}, {0, 2}, {0, 4}},
			p:              Point2{0, 6},
			wantDistance:   2,
			wantProjection: ProjectionOnPolyChain{N: 2, Position: 0, Point2: Point2{0, 4}},
		},
		{
			name:           "",
			polyChain:      PolyChain{{0, 0}, {0, 2}, {0, 4}},
			p:              Point2{1, 3},
			wantDistance:   1,
			wantProjection: ProjectionOnPolyChain{N: 1, Position: 1, Point2: Point2{0, 3}},
		},
		{
			name:           "",
			polyChain:      PolyChain{{0, 0}, {0, 2}, {0, 4}},
			p:              Point2{1, 2},
			wantDistance:   1,
			wantProjection: ProjectionOnPolyChain{N: 1, Position: 0, Point2: Point2{0, 2}},
		},
		{
			name:           "",
			polyChain:      PolyChain{{0, 0}, {0, 2}, {2, 2}},
			p:              Point2{1, 1},
			wantDistance:   1,
			wantProjection: ProjectionOnPolyChain{N: 0, Position: 1, Point2: Point2{0, 1}}, // , {N: 1, Position: 1, Point2: Point2{1, 2}}},
		},
		{
			name:           "",
			polyChain:      PolyChain{{0, 0}, {0, 2}, {2, 2}},
			p:              Point2{1, 1.5},
			wantDistance:   0.5,
			wantProjection: ProjectionOnPolyChain{N: 1, Position: 1, Point2: Point2{1, 2}},
		},
		{
			name:           "",
			polyChain:      PolyChain{{0, 0}, {0, 2}, {2, 2}},
			p:              Point2{0.5, 1},
			wantDistance:   0.5,
			wantProjection: ProjectionOnPolyChain{N: 0, Position: 1, Point2: Point2{0, 1}},
		},
		{
			name:           "",
			polyChain:      PolyChain{{0, 0}, {0, 2}, {2, 2}},
			p:              Point2{2, 0},
			wantDistance:   2,
			wantProjection: ProjectionOnPolyChain{N: 0, Position: 0, Point2: Point2{0, 0}}, // {N: 2, Position: 0, Point2: Point2{2, 2}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDistance, gotProjections := tt.p.DistanceToPolyChain(tt.polyChain)

			if math.Abs(gotDistance-tt.wantDistance) > mathlib.EPS {
				t.Errorf("TestDistanceToPolyChain() gotDistance = %f, wantDistance %f", gotDistance, tt.wantDistance)
			}
			if !reflect.DeepEqual(gotProjections, tt.wantProjection) {
				t.Errorf("TestDistanceToPolyChain() gotProjections = %v, wantProjection %v", gotProjections, tt.wantProjection)
			}
		})
	}
}

func TestDistanceToLineSegment(t *testing.T) {
	tests := []struct {
		name         string
		p            Point2
		ls           Segment
		wantDistance float64
		wantPosition float64
	}{
		{
			name: "", ls: Segment{{0, 0}, {4, 0}}, p: Point2{-5, 0},
			wantDistance: 5, wantPosition: 0,
		},
		{
			name: "", ls: Segment{{0, 0}, {4, 0}}, p: Point2{-3, -4},
			wantDistance: 5, wantPosition: 0,
		},
		{
			name: "", ls: Segment{{0, 0}, {4, 0}}, p: Point2{0, 4},
			wantDistance: 4, wantPosition: 0,
		},
		{
			name: "", ls: Segment{{0, 0}, {4, 0}}, p: Point2{1, 4},
			wantDistance: 4, wantPosition: 1,
		},
		{
			name: "", ls: Segment{{0, 0}, {4, 0}}, p: Point2{2, 4},
			wantDistance: 4, wantPosition: 2,
		},
		{
			name: "", ls: Segment{{0, 0}, {4, 0}}, p: Point2{3, 4},
			wantDistance: 4, wantPosition: 3,
		},
		{
			name: "", ls: Segment{{0, 0}, {4, 0}}, p: Point2{4, 0},
			wantDistance: 0, wantPosition: 4,
		},
		{
			name: "", ls: Segment{{0, 0}, {4, 0}}, p: Point2{2, 0},
			wantDistance: 0, wantPosition: 2,
		},
		{
			name: "",
			p:    Point2{0, 4}, ls: Segment{{0, 0}, {0, 2}},
			wantDistance: 2, wantPosition: 2,
		},
		{
			name: "",
			p:    Point2{0, 6}, ls: Segment{{0, 0}, {0, 2}},
			wantDistance: 4, wantPosition: 2,
		},
		{
			name: "",
			p:    Point2{0, 4}, ls: Segment{{0, 2}, {0, 4}},
			wantDistance: 0, wantPosition: 2,
		},
		{
			name: "",
			p:    Point2{0, 6}, ls: Segment{{0, 2}, {0, 4}},
			wantDistance: 2, wantPosition: 2,
		},
		{
			name: "",
			p:    Point2{1, 3}, ls: Segment{{0, 2}, {0, 4}},
			wantDistance: 1, wantPosition: 1,
		},
		{
			name: "",
			p:    Point2{1, 3.5}, ls: Segment{{0, 2}, {0, 4}},
			wantDistance: 1, wantPosition: 1.5,
		},
		{
			name: "",
			p:    Point2{1, 2.5}, ls: Segment{{0, 2}, {0, 4}},
			wantDistance: 1, wantPosition: 0.5,
		},
		{
			name: "",
			p:    Point2{0, 3}, ls: Segment{{0, 2}, {0, 4}},
			wantDistance: 0, wantPosition: 1,
		},
		{
			name: "",
			p:    Point2{0, 2}, ls: Segment{{0, 2}, {0, 4}},
			wantDistance: 0, wantPosition: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDistance, gotPosition := tt.p.DistanceToSegment(tt.ls)
			if math.Abs(gotDistance-tt.wantDistance) > mathlib.EPS {
				t.Errorf("DistanceToSegment() gotDistance = %v, wantDistance %v", gotDistance, tt.wantDistance)
			}
			if math.Abs(gotPosition-tt.wantPosition) > mathlib.EPS {
				t.Errorf("DistanceToSegment() gotPosition = %v, wantDistance %v", gotPosition, tt.wantPosition)
			}
		})
	}
}

func TestFilterPolyChain(t *testing.T) {
	tests := []struct {
		name      string
		pCh       PolyChain
		pChWanted PolyChain
	}{
		{
			name:      "",
			pCh:       PolyChain{{1, 1}, {1, 1}, {1, 0}, {1, 0}, {1, 0}, {1, 0}, {1, 0}},
			pChWanted: PolyChain{{1, 1}, {1, 0}},
		},
		{
			name:      "",
			pCh:       PolyChain{{1, 2}, {1, 1}, {1, 0}, {1, 0}, {1, 0}, {3, 1}, {1, 0}, {1, 0}},
			pChWanted: PolyChain{{1, 2}, {1, 1}, {1, 0}, {3, 1}, {1, 0}},
		},
		{
			name:      "",
			pCh:       PolyChain{{1, 2}, {1, 1}, {1, 0}, {3, 1}, {1, 0}},
			pChWanted: PolyChain{{1, 2}, {1, 1}, {1, 0}, {3, 1}, {1, 0}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pCh.Filter(); !reflect.DeepEqual(got, tt.pChWanted) {
				t.Errorf("FilterPolyChain() = %v, pChWanted %v", got, tt.pChWanted)
			}
		})
	}
}
