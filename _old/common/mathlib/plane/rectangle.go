package plane

type Rectangle struct {
	RectangleXY
	XToYAngle
}

func (rect Rectangle) Sides() [2]Segment {

	p00 := Point2{rect.HalfSideX, rect.HalfSideY}.RotateByAngle(rect.XToYAngle)
	p01 := Point2{-rect.HalfSideX, rect.HalfSideY}.RotateByAngle(rect.XToYAngle)

	return [2]Segment{{rect.RectangleXY.Point2.Add(p00), rect.RectangleXY.Point2.Add(p01)}, {rect.RectangleXY.Point2.Sub(p00), rect.RectangleXY.Point2.Sub(p01)}}
}

func (rect Rectangle) Contains(p2 Point2) bool {
	p2Rotated := p2.RotateAround(rect.RectangleXY.Point2, -rect.XToYAngle)
	return p2Rotated.X >= rect.RectangleXY.Point2.X-rect.HalfSideX && p2Rotated.X <= rect.RectangleXY.Point2.X+rect.HalfSideX &&
		p2Rotated.Y >= rect.RectangleXY.Point2.Y-rect.HalfSideY && p2Rotated.Y <= rect.RectangleXY.Point2.Y+rect.HalfSideY
}

func (rect Rectangle) Outer(margin float64) Rectangle {
	return Rectangle{
		RectangleXY: RectangleXY{
			Point2:    rect.RectangleXY.Point2,
			HalfSideX: rect.HalfSideX + margin,
			HalfSideY: rect.HalfSideY + margin,
		},
		XToYAngle: rect.XToYAngle,
	}
}

//func (rect Rectangle) Intersection(pCh PolyChain) PolyChain {
//
//	log.Fatal("on Rectangle.Intersects()")
//
//	//for _, p := range pCh {
//	//	p2Rot := RotateByAngle(p, -rect.XToYAngleFromOy)
//	//	if p2Rot.XT >= rect.Min.XT && p2Rot.XT <= rect.MaxIn.XT && p2Rot.YT >= rect.Min.YT && p2Rot.YT <= rect.MaxIn.YT {
//	//		return true
//	//	}
//	//}
//	return nil
//}
//
//func (rect Rectangle) IntersectionArea(rect1 Rectangle) float64 {
//
//	log.Fatal("on Rectangle.IntersectionArea()")
//
//	return 0
//}

func (rect Rectangle) Points() (p00, p01, p10, p11 Point2) {
	if rect.HalfSideX == 0 && rect.HalfSideY == 0 {
		return rect.RectangleXY.Point2, rect.RectangleXY.Point2, rect.RectangleXY.Point2, rect.RectangleXY.Point2
	}

	p00Fixed, p01Fixed := Point2{-rect.HalfSideX, -rect.HalfSideY}, Point2{-rect.HalfSideX, rect.HalfSideY}
	p00_, p01_ := p00Fixed.RotateByAngle(rect.XToYAngle), p01Fixed.RotateByAngle(rect.XToYAngle)

	return Point2{p00_.X + rect.RectangleXY.Point2.X, p00_.Y + rect.RectangleXY.Point2.Y},
		Point2{p01_.X + rect.RectangleXY.Point2.X, p01_.Y + rect.RectangleXY.Point2.Y},
		Point2{-p01_.X + rect.RectangleXY.Point2.X, -p01_.Y + rect.RectangleXY.Point2.Y},
		Point2{-p00_.X + rect.RectangleXY.Point2.X, -p00_.Y + rect.RectangleXY.Point2.Y}
}

func (rect Rectangle) OuterXY(margin float64) RectangleXY {
	p00, p01, p10, p11 := rect.Points()
	return RectangleAround(margin, p00, p01, p10, p11)
}
