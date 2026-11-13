package numbers

import (
	"fmt"
	"sort"
)

const onInterpolateByTable = "on numlib.InterpolateByTable()"

func InterpolateByTable(x float64, xy [][2]float64) (float64, error) {
	if len(xy) < 1 {
		return 0, fmt.Errorf(onInterpolateByTable + ": table is empty")
	}

	sort.Slice(xy, func(i, j int) bool { return xy[i][0] < xy[j][0] })
	for i := 0; i < len(xy); i++ {
		if x == xy[i][0] {
			return xy[i][1], nil
		} else if x < xy[i][0] {
			if i == 0 {
				return 0, fmt.Errorf(onInterpolateByTable+": x (%f) is smaller than tha first table argument (%f)", x, xy[0][0])
			}

			// TODO: if xy[i][0] == xy[i-1][0] && xy[i][1] != xy[i-1][1]
			return xy[i-1][1] + (x-xy[i-1][0])*(xy[i][1]-xy[i-1][1])/(xy[i][0]-xy[i-1][0]), nil
		}
	}

	return 0, fmt.Errorf(onInterpolateByTable+": x (%f) is greater than tha last table argument (%f)", x, xy[len(xy)-1][0])
}

const onInterpolateByTwoPoints = "on numlib.InterpolateByTwoPoints()"

func InterpolateByTwoPoints(x float64, xy [2][2]float64) (float64, error) {
	if xy[0][0] < xy[1][0] {
		if x < xy[0][0] || x > xy[1][0] {
			return 0, fmt.Errorf("onInterpolateBetweenTwoPoints: target x (%f) isn't between x0 (%f) and x1 (%f)", x, xy[0][0], xy[1][0])
		}
	} else {
		if x > xy[0][0] || x < xy[1][0] {
			return 0, fmt.Errorf("onInterpolateBetweenTwoPoints: target x (%f) isn't between x0 (%f) and x1 (%f)", x, xy[0][0], xy[1][0])
		} else if xy[1][0] == xy[0][0] {
			// TODO: if xy[1][1] != xy[0][1]
			return xy[0][1], nil
		}
	}

	return xy[0][1] + (x-xy[0][0])*(xy[1][1]-xy[0][1])/(xy[1][0]-xy[0][0]), nil
}
