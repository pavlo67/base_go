package geolib

import (
	"math"
	"testing"

	"github.com/pavlo67/common/common/mathlib/plane"
)

func TestPointInRanges(t *testing.T) {
	tests := []struct {
		name     string
		geoPoint Point
		moving   plane.Point2
		zoom     int
		tileSide int
	}{
		{
			name:     "",
			geoPoint: Point{38, 41},
			moving:   plane.Point2{100, 100},
			zoom:     18,
			tileSide: 256,
		},
		{
			name:     "",
			geoPoint: Point{38, 41},
			moving:   plane.Point2{100, -100},
			zoom:     18,
			tileSide: 256,
		},
		{
			name:     "",
			geoPoint: Point{38, 41},
			moving:   plane.Point2{500, 500},
			zoom:     18,
			tileSide: 256,
		},
		{
			name:     "",
			geoPoint: Point{38, 41},
			moving:   plane.Point2{500, -500},
			zoom:     18,
			tileSide: 256,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			geoPoint1 := tt.geoPoint.MovedAt(tt.moving)

			xyRanges := XYRangesAround(tt.geoPoint, tt.zoom, tt.moving.X, tt.moving.Y)

			p := PointInRanges(tt.geoPoint, xyRanges, tt.tileSide)
			p1 := PointInRanges(geoPoint1, xyRanges, tt.tileSide)

			dpm := DPM(tt.geoPoint.Lat, tt.zoom)

			moving1X, moving1Y := float64(p1.X-p.X)/dpm, float64(p.Y-p1.Y)/dpm

			if math.Abs(moving1X-tt.moving.X) > DistanceEps {
				t.Errorf("moving1X: %f, tt.moving.X: %f", moving1X, tt.moving.X)
			}
			if math.Abs(moving1X-tt.moving.X) > DistanceEps {
				t.Errorf("moving1Y: %f, tt.moving.Y: %f", moving1Y, tt.moving.Y)
			}

		})
	}
}
