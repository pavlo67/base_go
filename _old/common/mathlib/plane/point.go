package plane

import (
	"image"
	"math"
)

type Point2 struct {
	X, Y float64
}

func (p Point2) ImagePoint() image.Point {
	return image.Point{int(math.Round(p.X)), int(math.Round(p.Y))}
}

// TODO!!!! be careful: all single point angles are calculated in range -pi < angle <= pi

func (p Point2) VectorTo(p1 Point2) Point2 {
	return Point2{X: p1.X - p.X, Y: p1.Y - p.Y}
}

func (p Point2) Sub(p1 Point2) Point2 {
	return Point2{p.X - p1.X, p.Y - p1.Y}
}

func (p Point2) Add(p1 Point2) Point2 {
	return Point2{p.X + p1.X, p.Y + p1.Y}
}

func (p Point2) Radius() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

func (p Point2) OnRay(r float64) Point2 {
	r0 := p.Radius()
	if r0 == 0 {
		return Point2{}
	}

	return Point2{p.X * r / r0, p.Y * r / r0}
}

func (p Point2) DistanceTo(p1 Point2) float64 {
	return math.Sqrt((p.X-p1.X)*(p.X-p1.X) + (p.Y-p1.Y)*(p.Y-p1.Y))
}

func (p Point2) DistanceSquare(p1 Point2) float64 {
	return (p.X-p1.X)*(p.X-p1.X) + (p.Y-p1.Y)*(p.Y-p1.Y)
}

func (p Point2) DistanceToSegment(s Segment) (distance, projectionPosition float64) {
	d0, d1, d := p.DistanceSquare(s[0]), p.DistanceSquare(s[1]), s[0].DistanceSquare(s[1])
	var reversed bool
	if d1 < d0 {
		d0, d1 = d1, d0
		reversed = true
	}
	if d0+d <= d1 {
		if reversed {
			return math.Sqrt(d0), math.Sqrt(d)
		} else {
			return math.Sqrt(d0), 0
		}
	}

	c0 := (d0 + d - d1) / (2 * math.Sqrt(d))

	if reversed {
		distance, projectionPosition = math.Sqrt(d0-c0*c0), max(0, math.Sqrt(d)-c0)
	} else {
		distance, projectionPosition = math.Sqrt(d0-c0*c0), max(0, c0)
	}

	if math.IsNaN(distance) {
		return 0, projectionPosition
	}

	return distance, projectionPosition
}

func (p Point2) RotateAround(center Point2, angle XToYAngle) Point2 {
	pToCenter := Point2{p.X - center.X, p.Y - center.Y}
	pToCenterRotated := pToCenter.RotateByAngle(angle)

	return Point2{pToCenterRotated.X + center.X, pToCenterRotated.Y + center.Y}
}

func (p Point2) RotateByAngle(addAngle XToYAngle) Point2 {
	angle := p.XToYAngleFromOx()
	r := math.Sqrt(p.X*p.X + p.Y*p.Y)

	return Point2{r * math.Cos(float64(angle+addAngle)), r * math.Sin(float64(angle+addAngle))}
}

func (p Point2) RotateWithRatio(ratio float64) Point2 {
	angle := p.XToYAngleFromOx()
	r := math.Sqrt(p.X*p.X + p.Y*p.Y)

	return Point2{r * math.Cos(float64(angle)*ratio), r * math.Sin(float64(angle)*ratio)}
}

// TODO!!!! be careful: axis angles are calculated in range -pi/2 < angle <= pi/2
// TODO!!!! be careful:	if math.Abs(angle - math.Pi / 2) <= Eps then deltaFromCenter == DX, else == DY

//func CenterAxis(axis Segment) (angle, deltaFromCenter float64) {
//	axisDX, axisDY := axis.End.Position-axis.Begin.Position, axis.End.Y-axis.Begin.Y
//
//	if axisDX == 0 {
//		if axisDY == 0 {
//			// TODO!!! be careful: it's a convention only
//			return 0, axis.Begin.Y
//		}
//		return math.Pi / 2, axis.Begin.Position
//	}
//
//	axisDerivative := axisDY / axisDX
//	if math.IsInf(axisDerivative, 0) {
//		return math.Pi / 2, axis.Begin.Position
//	}
//
//	angle = math.Atan(axisDerivative)
//	if math.Abs(angle-math.Pi/2) <= Eps {
//		return angle, axis.Begin.Position - axis.Begin.Y/axisDerivative
//	}
//
//	return angle, axis.Begin.Y - axis.Begin.Position*axisDerivative
//}
//
//func SetAlongOx(points []Point2, axis Segment) ([]Point2, image.Rect) {
//	angle, deltaFromCenter := CenterAxis(axis)
//	pointsMoved := make([]Point2, len(points))
//
//	if math.Abs(angle-math.Pi/2) <= Eps {
//		for i, p := range points {
//			pointsMoved[i] = RotateByAngle(Point2{p.Position - deltaFromCenter, p.Y}, -angle)
//		}
//	} else {
//		for i, p := range points {
//			pointsMoved[i] = RotateByAngle(Point2{p.Position, p.Y - deltaFromCenter}, -angle)
//		}
//	}
//
//	return pointsMoved
//}

func (p Point2) DistanceToLine(line Segment) float64 {
	p2 := line.Vector()
	pIntersect := line.LinesIntersection(Segment{p, Point2{p.X - p2.Y, p.Y + p2.X}})
	if pIntersect == nil {
		// TODO!!! be careful: it's impossible
		return math.NaN()
	}

	return math.Sqrt((pIntersect.X-p.X)*(pIntersect.X-p.X) + (pIntersect.Y-p.Y)*(pIntersect.Y-p.Y))
}

func Diameter(pts []Point2) float64 {
	distMax := 0.
	for i, p0 := range pts {
		for _, p1 := range pts[i+1:] {
			if dist := p0.DistanceTo(p1); dist > distMax {
				distMax = dist
			}
		}
	}

	return distMax
}
