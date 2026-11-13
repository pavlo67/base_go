package geolib

import "math"

type DMS struct {
	D, M int
	S    float64
}

func (dms DMS) Degrees() Degrees {
	return Degrees(float64(dms.D) + float64(dms.M)/60 + float64(dms.S)/3600)
}

type Degrees float64

func (degrees Degrees) DMS() DMS {
	d := int(degrees)
	dRest := math.Abs(float64(degrees) - float64(d))

	m := int(dRest * 60)
	s := (dRest - float64(m)/60) * 3600

	return DMS{d, m, s}
}

func (degrees Degrees) Angle() float64 {
	return math.Pi * float64(degrees) / 180
}

func (degrees Degrees) Canon() Degrees {
	for degrees >= 360 {
		degrees -= 360
	}
	for degrees <= -360 {
		degrees += 360
	}

	return degrees
}

func DegreesFromAngle(angle float64) Degrees {
	return Degrees(180 * float64(angle) / math.Pi)
}
