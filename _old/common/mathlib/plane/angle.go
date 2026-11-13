package plane

import "math"

// XToYAngle is angle in radians (-Pi < rotation <= Pi), rotation from Ox to Oy (counter clockwise) has a positive angle
type XToYAngle float64

// Position moves geometry shape to Point2 and (after!) rotates it with XToYAngle angle around itself
type Position struct {
	Point2
	XToYAngle
}

func (xToYAngle XToYAngle) Canon() XToYAngle {
	if xToYAngle > math.Pi {
		return xToYAngle - 2*math.Pi
	} else if xToYAngle <= -math.Pi {
		return xToYAngle + 2*math.Pi
	}

	return xToYAngle
}

func (xToYAngle XToYAngle) Point2(radius float64) Point2 {
	return Point2{radius * math.Cos(float64(xToYAngle)), radius * math.Sin(float64(xToYAngle))}
}

// XToYAngle lies in the range: -math.Pi < p.XToYAngleFromOy() <= math.Pi
func (p Point2) XToYAngleFromOx() XToYAngle {
	if p.X == 0 {
		if p.Y > 0 {
			return math.Pi / 2
		} else if p.Y < 0 {
			return -math.Pi / 2
		} else {
			return XToYAngle(math.NaN())
		}
	} else if p.X >= 0 {
		return XToYAngle(math.Atan(p.Y / p.X))
	} else if p.Y >= 0 {
		return XToYAngle(math.Atan(p.Y/p.X) + math.Pi)
	} else {
		return XToYAngle(math.Atan(p.Y/p.X) - math.Pi)
	}
}

// DEPRECATED
func (p Point2) AnglesDelta(p1 Point2) float64 {
	return -p.AngleFrom(p1)
}

func (p Point2) AngleFrom(p1 Point2) float64 {
	angle := p.XToYAngleFromOx() - p1.XToYAngleFromOx()
	if angle > math.Pi {
		return float64(angle - 2*math.Pi)
	} else if angle <= -math.Pi {
		return float64(angle + 2*math.Pi)
	}
	return float64(angle)
}
