package geolib

import (
	"fmt"

	"github.com/pavlo67/common/common/mathlib/plane"
)

type XYRange [2]int

type XYRanges struct {
	Zoom   int
	XT, YT XYRange
}

func (xyRanges XYRanges) Key() string {
	return fmt.Sprintf("z%d_x%d_x%d_y%d_y%d", xyRanges.Zoom, xyRanges.XT[0], xyRanges.XT[1], xyRanges.YT[0], xyRanges.YT[1])
}

func (xyRanges XYRanges) Canon() XYRanges {
	xMin, xMax, yMin, yMax := xyRanges.XT[0], xyRanges.XT[1], xyRanges.YT[0], xyRanges.YT[1]
	if xMax < xMin {
		xMin, xMax = xMax, xMin
	}
	if yMax < yMin {
		yMin, yMax = yMax, yMin
	}

	return XYRanges{xyRanges.Zoom, XYRange{xMin, xMax}, XYRange{yMin, yMax}}
}

//func (xyRanges XYRanges) Area() Area {
//	xyRanges = xyRanges.Canon()
//
//	// TODO!!! use +1 in XYRanges max itself
//
//	return Area{
//		Tile{X: xyRanges.XT[0], Y: xyRanges.YT[0], Zoom: xyRanges.Zoom}.LeftTop(),
//		Tile{X: xyRanges.XT[1] + 1, Y: xyRanges.YT[1] + 1, Zoom: xyRanges.Zoom}.LeftTop(),
//	}
//}

func (xyRanges XYRanges) LeftTop() Point {

	xyRanges = xyRanges.Canon()

	return Tile{X: xyRanges.XT[0], Y: xyRanges.YT[0], Zoom: xyRanges.Zoom}.LeftTop()
}

func XYRangesAround(geoPoint Point, zoom int, hsX, hsY float64) XYRanges {

	tileMin, tileMax := geoPoint.MovedAt(plane.Point2{-hsX, hsY}).Tile(zoom), geoPoint.MovedAt(plane.Point2{hsX, -hsY}).Tile(zoom)

	if tileMax.X < tileMin.X {
		tileMin.X, tileMax.X = tileMax.X, tileMin.X
	}
	if tileMax.Y < tileMin.Y {
		tileMin.Y, tileMax.Y = tileMax.Y, tileMin.Y
	}

	return XYRanges{Zoom: zoom, XT: XYRange{tileMin.X, tileMax.X}, YT: XYRange{tileMin.Y, tileMax.Y}}

}

func PointInRanges(geoPoint Point, xyRanges XYRanges, tileSide int) plane.Point2 {

	tile := geoPoint.Tile(xyRanges.Zoom)
	leftTop := tile.LeftTop()
	rightBottom := Tile{tile.X + 1, tile.Y + 1, tile.Zoom}.LeftTop()

	x := float64(tileSide) * (float64(tile.X-xyRanges.XT[0]) + float64(geoPoint.Lon-leftTop.Lon)/float64(rightBottom.Lon-leftTop.Lon))
	y := float64(tileSide) * (float64(tile.Y-xyRanges.YT[0]) + float64(geoPoint.Lat-leftTop.Lat)/float64(rightBottom.Lat-leftTop.Lat))

	return plane.Point2{x, y}
}
