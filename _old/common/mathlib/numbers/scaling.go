package numbers

import "math"

func DividerToRangeInt(v0, vMin, vMax float64) int {
	if vMin > vMax {
		vMin, vMax = vMax, vMin
	}

	if v0 == 0 || math.IsInf(v0, 0) || math.IsNaN(v0) {
		return 0
	} else if v0 >= vMin && v0 <= vMax {
		return 1
	}

	dividerMin := v0 / vMax
	dividerMax := v0 / vMin
	if math.Ceil(dividerMin) <= math.Floor(dividerMax) {
		return int(math.Ceil(dividerMin)+math.Floor(dividerMax)) / 2
	}

	return 0
}

func ScaleToRangeFloat(v0, vMin, vMax float64) float64 {
	if vMin > vMax {
		vMin, vMax = vMax, vMin
	}

	if v0 == 0 || math.IsInf(v0, 0) || math.IsNaN(v0) {
		return 0
	} else if v0 >= vMin && v0 <= vMax {
		return 1
	}

	return (vMin + vMax) / v0
}
