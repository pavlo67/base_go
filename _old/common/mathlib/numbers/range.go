package numbers

type Range [2]float64

func (r *Range) Canon() {
	if r[0] > r[1] {
		r[0], r[1] = r[1], r[0]
	}
}
