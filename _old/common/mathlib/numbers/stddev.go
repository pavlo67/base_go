package numbers

import "math"

func StdDev(data []float64) float64 {
	if len(data) <= 1 {
		return 0
	}

	sum := 0.
	for _, d := range data {
		sum += d
	}
	avg := sum / float64(len(data))

	dev := 0.
	for _, d := range data {
		dev += (d - avg) * (d - avg)
	}

	return math.Sqrt(dev / float64(len(data)))
}
