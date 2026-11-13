package geolib

import (
	"math"

	"github.com/pavlo67/common/common/mathlib"
)

func DPM(lat Degrees, zoom int) float64 {

	n := math.Pow(2, float64(zoom))
	mpd := 156543.03 * math.Cos(float64(lat)*math.Pi/180) / n // resolution: meters per dot (mpd)

	if mpd <= mathlib.EPS {
		return math.Inf(1)
	}

	return 1 / mpd
}

type Tile struct {
	X, Y int
	Zoom int
}

// https://wiki.openstreetmap.org/wiki/Slippy_map_tilenames
// left-top corner of the tile

func (tile Tile) LeftTop() Point {
	n := math.Pow(2, float64(tile.Zoom))
	latAngle := math.Atan(math.Sinh(math.Pi * (1 - 2*float64(tile.Y)/n)))

	return Point{
		Lat: Degrees(latAngle * 180 / math.Pi),
		Lon: Degrees(float64(tile.X)/n*360 - 180),
	}
}

func (tile Tile) Center() Point {
	tilesN := int(math.Round(math.Pow(2, float64(tile.Zoom))))
	leftTop, rightBottom := tile.LeftTop(), Tile{(tile.X + 1) % tilesN, (tile.Y + 1) % tilesN, tile.Zoom}.LeftTop()

	return Point{
		Lat: 0.5 * (leftTop.Lat + rightBottom.Lat),
		Lon: 0.5 * (leftTop.Lon + rightBottom.Lon),
	}
}
