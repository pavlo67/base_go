package numbers

import (
	"math"
)

func QuadraticEquation(a, b, c float64) *[2]float64 {

	discriminant := (b * b) - (4 * a * c)

	if discriminant >= 0 {
		return &[2]float64{(-b + math.Sqrt(discriminant)) / (2 * a), (-b - math.Sqrt(discriminant)) / (2 * a)}
	}

	// imaginary = math.Sqrt(-discriminant) / (2 * a)

	return nil
}
