package plane

import "math"

func (s Segment) TurnAroundAxis(p2 Point2) Point2 {
	axisDX, axisDY := s[1].X-s[0].X, s[1].Y-s[0].Y
	var axisDerivative float64

	if axisDX == 0 {
		if axisDY == 0 {
			// TODO!!! be careful, it's a very non-standard case
			return Point2{s[0].X*2 - p2.X, s[0].Y*2 - p2.Y}
		}
		return Point2{s[0].X*2 - p2.X, p2.Y}
	} else if axisDerivative = axisDY / axisDX; math.IsInf(axisDerivative, 0) {
		return Point2{s[0].X*2 - p2.X, p2.Y}
	} else if axisDY == 0 || math.IsInf(1/axisDerivative, 0) {
		return Point2{p2.X, s[0].Y*2 - p2.Y}
	}

	pIntersection := s.LinesIntersection(Segment{p2, Point2{p2.X + 1, p2.Y - 1/axisDerivative}})
	return Point2{pIntersection.X*2 - p2.X, pIntersection.Y*2 - p2.Y}

}

func (s Segment) TurnAroundAxisMultiple(p2s ...Point2) []Point2 {
	axisDX, axisDY := s[1].X-s[0].X, s[1].Y-s[0].Y
	var axisDerivative float64

	p2sTurned := make([]Point2, len(p2s))

	if axisDX == 0 {
		if axisDY == 0 {
			// TODO!!! be careful, it's a very non-standard case
			for i, p2 := range p2s {
				p2sTurned[i] = Point2{s[0].X*2 - p2.X, s[0].Y*2 - p2.Y}
			}
		} else {
			for i, p2 := range p2s {
				p2sTurned[i] = Point2{s[0].X*2 - p2.X, p2.Y}
			}
		}
	} else if axisDerivative = axisDY / axisDX; math.IsInf(axisDerivative, 0) {
		for i, p2 := range p2s {
			p2sTurned[i] = Point2{s[0].X*2 - p2.X, p2.Y}
		}
	} else if axisDY == 0 || math.IsInf(1/axisDerivative, 0) {
		for i, p2 := range p2s {
			p2sTurned[i] = Point2{p2.X, s[0].Y*2 - p2.Y}
		}
	} else {
		for i, p2 := range p2s {
			pIntersection := s.LinesIntersection(Segment{p2, Point2{p2.X + 1, p2.Y - 1/axisDerivative}})
			p2sTurned[i] = Point2{pIntersection.X*2 - p2.X, pIntersection.Y*2 - p2.Y}
		}
	}

	return p2sTurned
}
