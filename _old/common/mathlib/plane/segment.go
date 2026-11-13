package plane

import (
	"log"
	"math"

	"github.com/pavlo67/common/common/mathlib/numbers"

	"github.com/pavlo67/common/common/mathlib"
)

type Segment [2]Point2

func (s *Segment) Vector() Point2 {
	if s == nil {
		return Point2{}
	}
	return Point2{s[1].X - s[0].X, s[1].Y - s[0].Y}
}

func (s Segment) Middle() Point2 {
	return Point2{(s[1].X + s[0].X) / 2, (s[1].Y + s[0].Y) / 2}
}

func (segment Segment) Paired(distanceToRight float64) Segment {
	if segment[0] == segment[1] || distanceToRight == 0 {
		return segment
	}

	direction := segment[1].Sub(segment[0])
	dirRadius := direction.Radius()
	dirDistance := Point2{direction.X * distanceToRight / dirRadius, direction.Y * distanceToRight / dirRadius}

	dirToTheSide := dirDistance.RotateByAngle(XToYAngle(math.Pi * 0.5))
	return Segment{segment[0].Add(dirToTheSide), segment[1].Add(dirToTheSide)}
}

// we check the ray directed from ls[1] to infinity
func (ls Segment) ProjectionOnPolyChain(pCh PolyChain, distanceMax float64) *ProjectionOnPolyChainDirected {
	//log.Printf("%v + %v", pCh, ls)

	if len(pCh) < 1 {
		return nil

	}

	segmLen := ls[0].DistanceTo(ls[1])

	if len(pCh) == 1 {
		distance0, distance1 := pCh[0].DistanceTo(ls[0]), pCh[0].DistanceTo(ls[1])
		if math.Abs(segmLen+distance1-distance0) <= mathlib.EPS {
			return &ProjectionOnPolyChainDirected{Distance: distance1, ProjectionOnPolyChain: ProjectionOnPolyChain{Point2: pCh[0]}}
		}
		return nil
	}

	var pr *ProjectionOnPolyChainDirected

	for n, pn := range pCh[:len(pCh)-1] {
		if p := ls.DivideByLine(pn, pCh[n+1]); p != nil {

			//log.Print(*p)

			if distance1 := p.DistanceTo(ls[1]); distance1 <= distanceMax && (pr == nil || distance1 < pr.Distance) {

				distance0 := p.DistanceTo(ls[0])

				//log.Print(distance0, distance1, segmLen)

				if math.Abs(segmLen+distance1-distance0) <= mathlib.EPS {
					position := p.DistanceTo(pn)
					if position >= pn.DistanceTo(pCh[n+1]) {
						pr = &ProjectionOnPolyChainDirected{
							Distance:              distance1,
							ProjectionOnPolyChain: ProjectionOnPolyChain{N: n + 1, Point2: *p}}
					} else {
						pr = &ProjectionOnPolyChainDirected{
							Distance:              distance1,
							ProjectionOnPolyChain: ProjectionOnPolyChain{N: n, Position: position, Point2: *p}}
					}
				}
			}
		}
	}

	return pr
}

func (s Segment) GoOutCircle(p Point2, r float64) *Point2 {
	if s[0].DistanceTo(p) > r || s[1].DistanceTo(p) <= r {
		return nil
	}

	x0, y0, dx, dy := s[0].X-p.X, s[0].Y-p.Y, s[1].X-s[0].X, s[1].Y-s[0].Y
	a, b, c := dx*dx+dy*dy, 2*(dx*x0+dy*y0), x0*x0+y0*y0-r*r

	roots := numbers.QuadraticEquation(a, b, c)
	if roots == nil {
		return nil
	} else if roots[0] >= 0 && roots[0] <= 1 {
		return &Point2{s[0].X + dx*roots[0], s[0].Y + dy*roots[0]}
	} else if roots[1] >= 0 && roots[1] <= 1 {
		return &Point2{s[0].X + dx*roots[1], s[0].Y + dy*roots[1]}
	}

	log.Printf("wrong roots: %v for s = %v, p = %v, r = %g / on s.GoOutCircle()", *roots, s, p, r)

	return nil
}

func (s Segment) Intersection(s1 Segment) (pCross *Point2) {
	if s[1].X < s[0].X {
		s[0], s[1] = s[1], s[0]
	}
	if s1[1].X < s1[0].X {
		s1[0], s1[1] = s1[1], s1[0]
	}
	if s1[0].X < s[0].X {
		s, s1 = s1, s
	}
	if s1[0].X-s[1].X >= mathlib.EPS {
		return nil
	}

	r := Point2{s[1].X - s[0].X, s[1].Y - s[0].Y}
	l := Point2{s1[1].X - s1[0].X, s1[1].Y - s1[0].Y}

	cr := Cross(r, l)

	if cr > -mathlib.EPS && cr < mathlib.EPS {
		// vertical segments
		if math.Abs(Cross(r, Point2{0, 1})) < mathlib.EPS {
			if s[1].X-s1[0].X >= mathlib.EPS {
				return nil
			}

			if s[1].Y < s[0].Y {
				s = Segment{s[1], s[0]}
			}
			if s1[1].Y < s1[0].Y {
				s1 = Segment{s1[1], s1[0]}
			}

			if s[0].Y < s1[0].Y {
				if s[1].Y >= s1[0].Y {
					return &s1[0]
				}
			} else if s[0].Y > s1[0].Y {
				if s1[1].Y >= s[0].Y {
					return &s[0]
				}
			} else {
				return &s[0]
			}
			return nil
		}

		// compare s1[0].Y and corresponding point on s
		k := r.Y / r.X
		s01Y := s[0].Y + k*(s1[0].X-s[0].X)
		if math.Abs(s01Y/math.Sqrt(1+k*k)) >= mathlib.EPS {
			return nil
		}

		return &s1[0]
		//} else if math.Abs(s[0].X-s[1].X) < mathlib.EPS {
		//
		//	s, s1 = s1, s
		//	r = Point2{s[1].X - s[0].X, s[1].Y - s[0].Y}
		//	l = Point2{s1[1].X - s1[0].X, s1[1].Y - s1[0].Y}
		//	cr = Cross(r, l)

	}

	q := Point2{s1[0].X - s[0].X, s1[0].Y - s[0].Y}
	t := Cross(q, l) / cr
	xIntersect := s[0].X + t*r.X
	yIntersect := s[0].Y + t*r.Y

	if math.Abs(s[0].X-s[1].X) < mathlib.EPS {
		s0Y, s1Y := s[0].Y, s[1].Y
		if s0Y > s1Y {
			s0Y, s1Y = s1Y, s0Y
		}

		if yIntersect < s0Y || yIntersect > s1Y || xIntersect < s1[0].X || xIntersect > s1[1].X {
			return nil
		}

	} else {
		s10Y, s11Y := s1[0].Y, s1[1].Y
		if s10Y > s11Y {
			s10Y, s11Y = s11Y, s10Y
		}

		if xIntersect < s[0].X || xIntersect > s[1].X || yIntersect < s10Y || yIntersect > s11Y {
			return nil
		}

	}

	return &Point2{xIntersect, yIntersect}
}

func (s Segment) DistanceTo(s1 Segment) float64 {
	if pCross := s.Intersection(s1); pCross != nil {
		return 0
	}

	d00, _ := s[0].DistanceToSegment(s1)
	d01, _ := s[1].DistanceToSegment(s1)
	d10, _ := s1[0].DistanceToSegment(s)
	d11, _ := s1[1].DistanceToSegment(s)

	return min(d00, d01, d10, d11)
}

func (s Segment) AngleAbs(s1 Segment) float64 {
	angle := math.Abs(s.Vector().AnglesDelta(s1.Vector()))
	for angle > math.Pi {
		angle -= math.Pi
	}

	return angle
}

func (s Segment) AngleLinesAbs(s1 Segment) float64 {
	angle := math.Abs(s.Vector().AnglesDelta(s1.Vector()))
	for angle > math.Pi/2 {
		angle -= math.Pi / 2
	}

	return angle
}
