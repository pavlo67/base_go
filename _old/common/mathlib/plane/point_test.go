package plane

import (
	"math"
	"testing"

	"github.com/pavlo67/common/common/mathlib"

	"github.com/stretchr/testify/require"
)

func TestRotateByAngle(t *testing.T) {
	type args struct {
		p        Point2
		addAngle XToYAngle
	}
	tests := []struct {
		name string
		args args
		want Point2
	}{
		{
			name: "0+pi",
			args: args{p: Point2{1, 0}, addAngle: math.Pi},
			want: Point2{-1, 0},
		},
		{
			name: "pi+0",
			args: args{p: Point2{-1, 0}, addAngle: 0},
			want: Point2{-1, 0},
		},
		{
			name: "0-pi/2",
			args: args{p: Point2{1, 0}, addAngle: -math.Pi / 2},
			want: Point2{0, -1},
		},
		{
			name: "pi/4+pi/2",
			args: args{p: Point2{1, 1}, addAngle: math.Pi / 2},
			want: Point2{-1, 1},
		},
		{
			name: "3*pi/4+pi/2",
			args: args{p: Point2{-1, 1}, addAngle: math.Pi / 2},
			want: Point2{-1, -1},
		},
		{
			name: "5*pi/4-pi/2",
			args: args{p: Point2{-1, -1}, addAngle: -math.Pi / 2},
			want: Point2{-1, 1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.p.RotateByAngle(tt.args.addAngle); got.DistanceTo(tt.want) > mathlib.EPS {
				t.Errorf("RotateByAngle() = %v, wantDistance %v", got, tt.want)
			}
		})
	}
}

func TestRotateWithRatio(t *testing.T) {
	type args struct {
		p     Point2
		ratio float64
	}
	tests := []struct {
		name string
		args args
		want Point2
	}{
		{
			name: "pi/4*1",
			args: args{p: Point2{1, 1}, ratio: 1},
			want: Point2{1, 1},
		},
		{
			name: "pi/4*3",
			args: args{p: Point2{1, 1}, ratio: 3},
			want: Point2{-1, 1},
		},
		{
			name: "pi/4*5",
			args: args{p: Point2{1, 1}, ratio: 5},
			want: Point2{-1, -1},
		},
		{
			name: "pi/4*7",
			args: args{p: Point2{1, 1}, ratio: 7},
			want: Point2{1, -1},
		},

		{
			name: "5*pi/4*2",
			args: args{p: Point2{-1, -1}, ratio: 2},
			want: Point2{0, 1.4142135623730951},
		},
		{
			name: "5*pi/4*3",
			args: args{p: Point2{-1, -1}, ratio: 3},
			want: Point2{1, -1},
		},

		{
			name: "pi/2*1",
			args: args{p: Point2{0, 1}, ratio: 1},
			want: Point2{0, 1},
		},
		{
			name: "pi/2*2",
			args: args{p: Point2{0, 1}, ratio: 2},
			want: Point2{-1, 0},
		},
		{
			name: "pi/2*3",
			args: args{p: Point2{0, 1}, ratio: 3},
			want: Point2{0, -1},
		},
		{
			name: "pi/2*4",
			args: args{p: Point2{0, 1}, ratio: 4},
			want: Point2{1, 0},
		},
		{
			name: "pi*1",
			args: args{p: Point2{-1, 0}, ratio: 1},
			want: Point2{-1, 0},
		},
		{
			name: "pi*2",
			args: args{p: Point2{-1, 0}, ratio: 2},
			want: Point2{1, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.p.RotateWithRatio(tt.args.ratio); got.DistanceTo(tt.want) > mathlib.EPS {
				t.Errorf("RotateByAngle() = %v, wantDistance %v", got, tt.want)
			}
		})
	}
}

func TestTurnAroundAxis(t *testing.T) {
	type args struct {
		p    Point2
		axis Segment
	}
	tests := []struct {
		name string
		args args
		want Point2
	}{
		{
			name: "Oy/Ox",
			args: args{p: Point2{0, 1}, axis: Segment{Point2{0, 0}, Point2{1, 0}}},
			want: Point2{0, -1},
		},
		{
			name: "Ox/Oy",
			args: args{p: Point2{-1, 1}, axis: Segment{Point2{0, -5}, Point2{0, 1}}},
			want: Point2{1, 1},
		},
		{
			name: "Oy/45",
			args: args{p: Point2{0, 1}, axis: Segment{Point2{-3, -3}, Point2{0, 0}}},
			want: Point2{1, 0},
		},
		{
			name: "Oy/45-1",
			args: args{p: Point2{0, 1}, axis: Segment{Point2{-3, -4}, Point2{0, -1}}},
			want: Point2{2, -1},
		},
		{
			name: "Ox/-45",
			args: args{p: Point2{1, 0}, axis: Segment{Point2{0, 0}, Point2{1, -1}}},
			want: Point2{0, -1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.axis.TurnAroundAxis(tt.args.p); got.DistanceTo(tt.want) > mathlib.EPS {
				t.Errorf("TurnAroundAxis() = %v, wantDistance %v", got, tt.want)
			}
		})
	}
}

func TestAngle2(t *testing.T) {
	type args struct {
		v0 Point2
		v1 Point2
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{args: args{v0: Point2{X: 1, Y: 0}, v1: Point2{X: 1, Y: 0}}, want: 0},
		{args: args{v0: Point2{X: 1, Y: 0}, v1: Point2{X: 1, Y: 1}}, want: math.Pi / 4},
		{args: args{v0: Point2{X: 1, Y: 0}, v1: Point2{X: 0, Y: 1}}, want: math.Pi / 2},
		{args: args{v0: Point2{X: 1, Y: 0}, v1: Point2{X: -1, Y: 1}}, want: 3 * math.Pi / 4},
		{args: args{v0: Point2{X: 1, Y: 0}, v1: Point2{X: -1, Y: 0}}, want: math.Pi},
		{args: args{v0: Point2{X: 1, Y: 0}, v1: Point2{X: -1, Y: -1}}, want: -3 * math.Pi / 4},
		{args: args{v0: Point2{X: 1, Y: 0}, v1: Point2{X: 0, Y: -1}}, want: -math.Pi / 2},
		{args: args{v0: Point2{X: 1, Y: 0}, v1: Point2{X: 1, Y: -1}}, want: -math.Pi / 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.v0.AnglesDelta(tt.args.v1); math.Abs(got-tt.want) > mathlib.EPS {
				t.Errorf("AnglesDelta() = %v, wantDistance %v", got, tt.want)
			}
		})
	}
}

func TestRotation(t *testing.T) {
	type args struct {
		p Point2
	}
	tests := []struct {
		name string
		args args
		want XToYAngle
	}{
		{args: args{p: Point2{X: 0, Y: 0}}, want: XToYAngle(math.NaN())},
		{args: args{p: Point2{X: 1, Y: 0}}, want: 0},
		{args: args{p: Point2{X: 1, Y: 1}}, want: math.Pi / 4},
		{args: args{p: Point2{X: 0, Y: 1}}, want: math.Pi / 2},
		{args: args{p: Point2{X: -1, Y: 1}}, want: 3 * math.Pi / 4},
		{args: args{p: Point2{X: -1, Y: 0}}, want: math.Pi},
		{args: args{p: Point2{X: 1, Y: -1}}, want: -math.Pi / 4},
		{args: args{p: Point2{X: 0, Y: -1}}, want: -math.Pi / 2},
		{args: args{p: Point2{X: -1, Y: -1}}, want: -3 * math.Pi / 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.p.XToYAngleFromOx(); math.Abs(float64(got-tt.want)) > mathlib.EPS {
				t.Errorf("XToYAngleFromOy() = %v, wantDistance %v", got, tt.want)
			}
		})
	}
}

//func TestAngle1(t *testing.T) {
//	type args struct {
//		p Point2
//	}
//	tests := []struct {
//		name string
//		args args
//		wantDistance float64
//	}{
//		{args: args{p: Point2{Position: 0, Y: 0}}, wantDistance: math.NaN()},
//		{args: args{p: Point2{Position: 1, Y: 0}}, wantDistance: 0},
//		{args: args{p: Point2{Position: 1, Y: 1}}, wantDistance: math.Pi / 4},
//		{args: args{p: Point2{Position: 0, Y: 1}}, wantDistance: math.Pi / 2},
//		{args: args{p: Point2{Position: -1, Y: 1}}, wantDistance: 3 * math.Pi / 4},
//		{args: args{p: Point2{Position: -1, Y: 0}}, wantDistance: math.Pi},
//		{args: args{p: Point2{Position: 1, Y: -1}}, wantDistance: -math.Pi / 4},
//		{args: args{p: Point2{Position: 0, Y: -1}}, wantDistance: -math.Pi / 2},
//		{args: args{p: Point2{Position: -1, Y: -1}}, wantDistance: -3 * math.Pi / 4},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := Angle1(tt.args.p); math.Abs(got-tt.wantDistance) > Eps {
//				t.Errorf("XToYAngleFromOy() = %v, wantDistance %v", got, tt.wantDistance)
//			}
//		})
//	}
//}

func TestDistanceToLine(t *testing.T) {
	tests := []struct {
		name string
		p    Point2
		line Segment
		want float64
	}{
		{
			name: "1",
			p:    Point2{38, 57.5},
			line: Segment{Point2{11.666666666666666, 38.666666666666664}, Point2{165.33333333333331, 38.5}},
			want: 18.861883338267475,
		},
		{
			name: "2",
			p:    Point2{25.333333333333332, 4.666666666666667},
			line: Segment{Point2{11.666666666666666, 38.666666666666664}, Point2{165.33333333333331, 38.5}},
			want: 33.98515716183312,
		},
		{
			name: "3",
			p:    Point2{25.333333333333332, 4.5},
			line: Segment{Point2{11.666666666666666, 38.5}, Point2{165.33333333333331, 38.5}},
			want: 34,
		},
		{
			name: "4",
			p:    Point2{25.3, 4.5},
			line: Segment{Point2{11, 38.5}, Point2{11, 48.5}},
			want: 14.3,
		},
		{
			name: "5",
			p:    Point2{1, 1},
			line: Segment{Point2{2, 2}, Point2{2, 4}},
			want: 1,
		},
		{
			name: "5",
			p:    Point2{1, 1},
			line: Segment{Point2{2, 2}, Point2{0, 0}},
			want: 0,
		},
		{
			name: "5",
			p:    Point2{1, 1},
			line: Segment{Point2{2, 2}, Point2{2, 2}},
			want: math.NaN(),
		},
		{
			name: "6",
			p:    Point2{0, 1},
			line: Segment{{}, {2, 2}},
			want: math.Sqrt(2) / 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.p.DistanceToLine(tt.line)
			if math.IsNaN(tt.want) {
				require.Truef(t, math.IsNaN(got), "wanted: %f, gotten: %f", tt.want, got)
			} else {
				require.Truef(t, math.Abs(got-tt.want) < mathlib.EPS, "wanted: %f, gotten: %f", tt.want, got)
			}
		})
	}
}
