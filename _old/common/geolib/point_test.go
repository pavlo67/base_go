package geolib

import (
	"math"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/mathlib/plane"
)

func TestPoint_MovedBeared(t *testing.T) {
	tests := []struct {
		name    string
		p       Point
		bearing Bearing
		moving  plane.Point2
	}{
		{
			name:    "",
			p:       Point{},
			bearing: 0,
			moving:  plane.Point2{400, 0},
		},

		{
			name:    "",
			p:       Point{45, 38},
			bearing: 33,
			moving:  plane.Point2{400, 200},
		},

		{
			name:    "",
			p:       Point{45, 38},
			bearing: -90,
			moving:  plane.Point2{1400, 200},
		},

		{
			name:    "",
			p:       Point{45, 38},
			bearing: -180,
			moving:  plane.Point2{2400, 200},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pMoved := tt.p.MovedBeared(tt.bearing, tt.moving)

			movingGeo := tt.moving.RotateByAngle(tt.bearing.XToYAngleFromOy())
			// movingGeoXToYAngleFromOx := movingGeo.XToYAngleFromOx()
			pMovedDirect := tt.p.MovedAt(movingGeo)

			// t.Logf("tt.p: %+v, movingGeo: %v, pMoved: %+v, pMovedDirect: %+v", tt.p, movingGeo, pMoved, pMovedDirect)

			require.Truef(t, math.Abs(tt.p.DistanceTo(pMoved)-movingGeo.Radius()) < DistanceEps,
				"tt.p.DistanceTo(pMoved): %f, movingGeo.Radius(): %f", tt.p.DistanceTo(pMoved), movingGeo.Radius())

			require.Truef(t, math.Abs(pMoved.DistanceTo(pMovedDirect)) < DistanceEps,
				"pMoved.DistanceTo(pMovedDirect): %f", pMoved.DistanceTo(pMovedDirect))

			//if !math.IsNaN(float64(movingGeoXToYAngleFromOx)) {
			//	require.Truef(t, math.Abs(float64(stepGeoXToYAngleFromOy-geoPointInitial.BearingTo(*geoPoint).XToYAngle().Canon())) < DistanceEps,
			//		"stepGeoXToYAngleFromOy: %f, geoPointInitial.BearingTo(*geoPoint): %f, geoPointInitial.BearingTo(*geoPoint).XToYAngle(): %f",
			//		stepGeoXToYAngleFromOy, geoPointInitial.BearingTo(*geoPoint), geoPointInitial.BearingTo(*geoPoint).XToYAngle().Canon())
			//}

			movingReturn := plane.Point2{}.Sub(tt.moving)
			geoPointReturn := pMoved.MovedBeared(tt.bearing, movingReturn)

			require.Truef(t, tt.p.DistanceTo(geoPointReturn) < DistanceEps,
				"tt.p.DistanceTo(geoPointReturn): %f", tt.p.DistanceTo(geoPointReturn))

		})
	}
}
