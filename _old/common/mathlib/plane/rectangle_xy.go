package plane

import (
	"image"
	"math"
)

type RectangleXY struct {
	Point2
	HalfSideX, HalfSideY float64
}

func (rectXY RectangleXY) Union(s RectangleXY) RectangleXY {
	if s.HalfSideX < 0 || s.HalfSideY < 0 {
		return rectXY
	} else if rectXY.HalfSideX < 0 || rectXY.HalfSideY < 0 {
		return s
	}
	minX := min(rectXY.X-rectXY.HalfSideX, s.X-s.HalfSideX)
	minY := min(rectXY.Y-rectXY.HalfSideY, s.Y-s.HalfSideY)
	maxX := max(rectXY.X+rectXY.HalfSideX, s.X+s.HalfSideX)
	maxY := max(rectXY.Y+rectXY.HalfSideY, s.Y+s.HalfSideY)

	return RectangleXY{
		Point2:    Point2{0.5 * (minX + maxX), 0.5 * (minY + maxY)},
		HalfSideX: 0.5 * (maxX - minX),
		HalfSideY: 0.5 * (maxY - minY),
	}
}

func RectangleAround(margin float64, pts ...Point2) RectangleXY {
	if len(pts) < 1 {
		return RectangleXY{}
	}

	minX, maxX, minY, maxY := pts[0].X, pts[0].X, pts[0].Y, pts[0].Y

	for _, p := range pts[1:] {
		if p.X <= minX {
			minX = p.X
		} else if p.X > maxX {
			maxX = p.X
		}
		if p.Y <= minY {
			minY = p.Y
		} else if p.Y > maxY {
			maxY = p.Y
		}
	}

	return RectangleXY{
		Point2{0.5 * (minX + maxX), 0.5 * (minY + maxY)},
		margin + 0.5*(maxX-minX),
		margin + 0.5*(maxY-minY)}
}

func (rectXY RectangleXY) Contains(p2 Point2) bool {
	minX, maxX, minY, maxY := rectXY.X-rectXY.HalfSideX, rectXY.X+rectXY.HalfSideX, rectXY.Y-rectXY.HalfSideY, rectXY.Y+rectXY.HalfSideY
	if minX > maxX {
		minX, maxX = maxX, minX
	}
	if minY > maxY {
		minY, maxY = maxY, minY
	}

	return p2.X >= minX && p2.X <= maxX && p2.Y >= minY && p2.Y <= maxY
}

//func (r RectangleXY) AverageWeighted(s RectangleXY, ratio float64) RectangleXY {
//	if s.HalfSideX < 0 || s.HalfSideY < 0 {
//		return r
//	} else if r.HalfSideX < 0 || r.HalfSideY < 0 {
//		return s
//	}
//
//	ratio = math.Abs(ratio)
//
//	minX := (r.X - r.HalfSideX + ratio*(s.X-s.HalfSideX)) / (1 + ratio)
//	minY := (r.Y - r.HalfSideY + ratio*(s.Y-s.HalfSideY)) / (1 + ratio)
//	maxX := (r.X + r.HalfSideX + ratio*(s.X+s.HalfSideX)) / (1 + ratio)
//	maxY := (r.Y + r.HalfSideY + ratio*(s.Y+s.HalfSideY)) / (1 + ratio)
//
//	return RectangleXY{
//		Point2:    Point2{0.5 * (minX + maxX), 0.5 * (minY + maxY)},
//		HalfSideX: 0.5 * (maxX - minX),
//		HalfSideY: 0.5 * (maxY - minY),
//	}
//}

//func (r RectangleXY) UnionWeighted(s RectangleXY, ratio float64) RectangleXY {
//	if s.HalfSideX < 0 || s.HalfSideY < 0 {
//		return r
//	} else if r.HalfSideX < 0 || r.HalfSideY < 0 {
//		return s
//	}
//
//	ratio = math.Abs(ratio)
//
//	minX := r.X - r.HalfSideX
//	if minXNew := s.X - s.HalfSideX; minXNew < minX {
//		minX = (minX + ratio*minXNew) / (1 + ratio)
//	}
//
//	minY := r.Y - r.HalfSideY
//	if minYNew := s.Y - s.HalfSideY; minYNew < minY {
//		minY = (minY + ratio*minYNew) / (1 + ratio)
//	}
//
//	maxX := r.X + r.HalfSideX
//	if maxXNew := s.X + s.HalfSideX; maxXNew > maxX {
//		maxX = (maxX + ratio*maxXNew) / (1 + ratio)
//	}
//
//	maxY := r.Y + r.HalfSideY
//	if maxYNew := s.Y + s.HalfSideY; maxYNew > maxY {
//		maxY = (maxY + ratio*maxYNew) / (1 + ratio)
//	}
//
//	return RectangleXY{
//		Point2:    Point2{0.5 * (minX + maxX), 0.5 * (minY + maxY)},
//		HalfSideX: 0.5 * (maxX - minX),
//		HalfSideY: 0.5 * (maxY - minY),
//	}
//}

func (rectXY RectangleXY) Intersects(s RectangleXY) bool {
	return math.Abs(rectXY.X-s.X) <= rectXY.HalfSideX+s.HalfSideX && math.Abs(rectXY.Y-s.Y) <= rectXY.HalfSideY+s.HalfSideY
}

func (rectXY RectangleXY) ImageRectangle() image.Rectangle {
	return image.Rectangle{
		Min: image.Point{int(math.Round(rectXY.X - rectXY.HalfSideX)), int(math.Round(rectXY.Y - rectXY.HalfSideY))},
		Max: image.Point{int(math.Round(rectXY.X + rectXY.HalfSideX)), int(math.Round(rectXY.Y + rectXY.HalfSideY))},
	}
}

func RectangleXYFromImageRectangle(rect image.Rectangle) RectangleXY {
	return RectangleXY{
		Point2:    Point2{0.5 * float64(rect.Min.X+rect.Max.X), 0.5 * float64(rect.Min.Y+rect.Max.Y)},
		HalfSideX: 0.5 * math.Abs(float64(rect.Max.X-rect.Min.X)),
		HalfSideY: 0.5 * math.Abs(float64(rect.Max.Y-rect.Min.Y)),
	}
}
