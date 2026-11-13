package plane

import "math"

func Center(points ...Point2) Point2 {
	if len(points) < 1 {
		return Point2{math.NaN(), math.NaN()}
	}
	var x, y float64
	for _, element := range points {
		x += element.X
		y += element.Y
	}

	n := float64(len(points))

	return Point2{X: x / n, Y: y / n}
}

func TriangleArea(p0, p1, p2 Point2) float64 {
	a, b, c := p0.DistanceTo(p1), p1.DistanceTo(p2), p2.DistanceTo(p0)

	//log.Printf("p0 (%v), p1 (%v), p2 (%v) --> a (%f), b (%f), c (%f)", p0, p1, p2, a, b, c)

	p := (a + b + c) / 2

	return math.Sqrt(p * (p - a) * (p - b) * (p - c))
}
