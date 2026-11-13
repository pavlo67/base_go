package plane

import "math"

type ProjectionOnPolyChain struct {
	N        int
	Position float64
	Point2
}

type ProjectionOnPolyChainDirected struct {
	Distance float64
	Angle    float64
	ProjectionOnPolyChain
}

func (p Point2) ProjectionOnLineSegment(ls Segment) (distance, projectionPosition float64) {
	d0, d1, d := p.DistanceSquare(ls[0]), p.DistanceSquare(ls[1]), ls[0].DistanceSquare(ls[1])
	var reversed bool
	if d1 < d0 {
		d0, d1 = d1, d0
		reversed = true
	}
	if d0+d < d1 {
		return math.NaN(), math.NaN()
	}

	c0 := (d0 + d - d1) / (2 * math.Sqrt(d))

	if reversed {
		return math.Sqrt(d0 - c0*c0), math.Sqrt(d) - c0
	} else {
		return math.Sqrt(d0 - c0*c0), c0
	}
}

func (p Point2) DistanceToPolyChain(pCh PolyChain) (float64, ProjectionOnPolyChain) {

	if len(pCh) < 1 {
		return math.NaN(), ProjectionOnPolyChain{N: -1, Position: math.NaN(), Point2: Point2{math.NaN(), math.NaN()}}
	} else if len(pCh) == 1 {
		return p.DistanceTo(pCh[0]), ProjectionOnPolyChain{Point2: pCh[0]}
	}

	minDist := math.Inf(1)
	var pr ProjectionOnPolyChain
	var n int
	var pPr Point2

	// POINTS:
	for i, pI := range pCh[:len(pCh)-1] {
		dist, position := p.DistanceToSegment(Segment{pI, pCh[i+1]})

		// fmt.Printf("{%g %g} {%g %g} {%g %g} --> %g\n", p.X, p.Y, pI.X, pI.Y, pCh[i+1].X, pCh[i+1].Y, dist)

		if dist >= minDist {
			continue
		}

		if segmentLength := pI.DistanceTo(pCh[i+1]); segmentLength <= 0 {
			pPr, n, position = pI, i, 0
		} else if position >= segmentLength {
			pPr, n, position = pCh[i+1], i+1, 0
		} else {
			dx, dy := pCh[i+1].X-pI.X, pCh[i+1].Y-pI.Y
			pPr, n = Point2{pI.X + dx*position/segmentLength, pI.Y + dy*position/segmentLength}, i
		}

		minDist, pr = dist, ProjectionOnPolyChain{n, position, pPr}

		//if dist < minDist {
		//	minDist, projections = dist, []ProjectionOnPolyChain{{n, position, pPr}}
		//} else {
		//	for _, pr := range projections {
		//		if n == pr.n && position == pr.pos {
		//			continue POINTS
		//		}
		//	}
		//
		//	projections = append(projections, ProjectionOnPolyChain{n, position, pPr})
		//}
	}

	return minDist, pr
}

func (p Point2) ProjectionOnPolyChain(pCh PolyChain) (float64, ProjectionOnPolyChain) {

	if len(pCh) < 1 {
		return math.NaN(), ProjectionOnPolyChain{N: -1, Position: math.NaN(), Point2: Point2{math.NaN(), math.NaN()}}
	} else if len(pCh) == 1 {
		return p.DistanceTo(pCh[0]), ProjectionOnPolyChain{Point2: pCh[0]}
	}

	minDist := math.Inf(1)
	pr := ProjectionOnPolyChain{
		N:        -1,
		Position: math.NaN(),
		Point2:   Point2{math.NaN(), math.NaN()},
	}
	var n int
	var pPr Point2

	// POINTS:
	for i, pI := range pCh[:len(pCh)-1] {
		dist, position := p.ProjectionOnLineSegment(Segment{pI, pCh[i+1]})
		if dist >= minDist || math.IsNaN(dist) {
			continue
		}

		if segmentLength := pI.DistanceTo(pCh[i+1]); segmentLength <= 0 {
			pPr, n, position = pI, i, 0
		} else if position >= segmentLength {
			pPr, n, position = pCh[i+1], i+1, 0
		} else {
			dx, dy := pCh[i+1].X-pI.X, pCh[i+1].Y-pI.Y
			pPr, n = Point2{pI.X + dx*position/segmentLength, pI.Y + dy*position/segmentLength}, i
		}

		minDist, pr = dist, ProjectionOnPolyChain{n, position, pPr}
	}

	return minDist, pr
}

func (p Point2) ProjectionsOnPolyChain(polyChain PolyChain, distanceMax float64) []ProjectionOnPolyChainDirected {
	if len(polyChain) < 1 {
		return nil

	} else if len(polyChain) == 1 {
		if distance := p.DistanceTo(polyChain[0]); distance <= distanceMax {
			return []ProjectionOnPolyChainDirected{{
				Distance:              distance,
				ProjectionOnPolyChain: ProjectionOnPolyChain{Point2: polyChain[0]},
			}}
		}
		return nil

	}

	var projections []ProjectionOnPolyChainDirected

	for n := 0; n < len(polyChain)-1; n++ {
		pn := polyChain[n]
		dist, position := p.DistanceToSegment(Segment{pn, polyChain[n+1]})

		// log.Print(dist, distanceMax)

		if dist > distanceMax {
			continue
		}

		segmLen := pn.DistanceTo(polyChain[n+1])
		if position < segmLen {
			projections = append(projections, ProjectionOnPolyChainDirected{dist, 0, ProjectionOnPolyChain{
				n, position, Point2{pn.X + (polyChain[n+1].X-pn.X)*position/segmLen, pn.Y + (polyChain[n+1].Y-pn.Y)*position/segmLen}}})
		} else if n == len(polyChain)-2 {
			projections = append(projections, ProjectionOnPolyChainDirected{dist, 0, ProjectionOnPolyChain{
				n + 1, 0, polyChain[n+1]}})
		}
	}

	return projections
}
