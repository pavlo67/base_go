package geolib

import (
	"math"

	"github.com/pavlo67/common/common/mathlib/plane"
)

// Bearing if a "right-angle" in degrees
type Bearing Degrees

func (bearing Bearing) Canon() Bearing {
	for bearing >= 360 {
		bearing -= 360
	}
	for bearing < 0 {
		bearing += 360
	}

	return bearing
}

func (bearing Bearing) CanonTo180() Bearing {
	for bearing > 180 {
		bearing -= 360
	}
	for bearing <= -180 {
		bearing += 360
	}

	return bearing
}

func BearingFromPoint(point plane.Point2) Bearing {
	bearingDegrees := 90 - (180 * point.XToYAngleFromOx() / math.Pi)

	for bearingDegrees >= 360 {
		bearingDegrees -= 360
	}
	for bearingDegrees < 0 {
		bearingDegrees += 360
	}

	return Bearing(bearingDegrees)
}

func BearingFromXToYAngleFromOy(rotation plane.XToYAngle) Bearing {
	bearingDegrees := -(180 * rotation / math.Pi)

	for bearingDegrees >= 360 {
		bearingDegrees -= 360
	}
	for bearingDegrees < 0 {
		bearingDegrees += 360
	}

	return Bearing(bearingDegrees)
}

// XToYAngleFromOy is calculated from Oy (as well as Bearing itself)
func (bearing Bearing) XToYAngleFromOy() plane.XToYAngle {
	angle := plane.XToYAngle(-bearing * math.Pi / 180)
	if angle <= -math.Pi {
		return angle + 2*math.Pi
	} else if angle > math.Pi {
		return angle - 2*math.Pi
	}

	return angle
}

func (bearing Bearing) Point(radius float64) plane.Point2 {
	return (math.Pi * (0.5 - plane.XToYAngle(bearing)/180)).Point2(radius)
}
