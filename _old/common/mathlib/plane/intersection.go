package plane

import (
	"github.com/pavlo67/common/common/mathlib"
)

type Intersection struct {
	Angle  float64
	X0, X1 float64
}

func Cross(p1, p2 Point2) float64 {
	return p1.X*p2.Y - p1.Y*p2.X
}

func (s Segment) LinesIntersection(s1 Segment) *Point2 {
	// https://stackoverflow.com/questions/7446126/opencv-2d-line-intersection-helper-function
	// https://stackoverflow.com/questions/563198/how-do-you-detect-where-two-line-segments-intersect

	r := Point2{s[1].X - s[0].X, s[1].Y - s[0].Y}
	l := Point2{s1[1].X - s1[0].X, s1[1].Y - s1[0].Y}

	cr := Cross(r, l)
	if cr > -mathlib.EPS && cr < mathlib.EPS {
		// lines are about parallel (they may be collinear!)
		return nil
	}

	q := Point2{s1[0].X - s[0].X, s1[0].Y - s[0].Y}
	t := Cross(q, l) / cr

	return &Point2{s[0].X + t*r.X, s[0].Y + t*r.Y}

}

func (s Segment) DivideByLine(p0, p1 Point2) *Point2 {
	pIntersect := s.LinesIntersection(Segment{p0, p1})
	if pIntersect == nil {
		return nil
	}

	if p1.X > p0.X {
		if pIntersect.X > p0.X && pIntersect.X < p1.X {
			return pIntersect
		}
	} else if p1.X < p0.X {
		if pIntersect.X > p1.X && pIntersect.X < p0.X {
			return pIntersect
		}
	} else if p1.Y > p0.Y {
		if pIntersect.Y > p0.Y && pIntersect.Y < p1.Y {
			return pIntersect
		}
	} else {
		if pIntersect.Y > p1.Y && pIntersect.Y < p0.Y {
			return pIntersect
		}
	}

	// log.Print(p0, p1, pIntersect)

	return nil
}

type PolyChainsIntersection struct {
	Point2
	N0, N1 int
}

func PolyChainsIntersectionAny(pCh0, pCh1 PolyChain) *PolyChainsIntersection {
	for i0 := 1; i0 < len(pCh0); i0++ {
		s0 := Segment{pCh0[i0-1], pCh0[i0]}
		for i1 := 1; i1 < len(pCh1); i1++ {
			if p := s0.Intersection(Segment{pCh1[i1-1], pCh1[i1]}); p != nil {
				return &PolyChainsIntersection{*p, i0 - 1, i1 - 1}
			}
		}
	}

	return nil
}
