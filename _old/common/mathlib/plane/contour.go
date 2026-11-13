package plane

type Contour PolyChain

func (c Contour) Rotations() []XToYAngle {
	rotations := make([]XToYAngle, len(c))

	for i, p := range c {
		rotations[i] = c[(i+1)%len(c)].Sub(p).XToYAngleFromOx()
	}

	return rotations
}

func (c Contour) Distances() []float64 {
	distances := make([]float64, len(c))

	for i, p := range c {
		distances[i] = p.DistanceTo(c[(i+1)%len(c)])
	}

	return distances
}

func (c Contour) Approximate(deviationMax float64) Contour {
	if len(c) <= 2 {
		return c
	}

	distances := c.Distances()
	var i0, i1 int
	var d0, d1 float64

	for i, d := range distances {
		d2 := d + distances[(i+len(distances)-1)%len(distances)]
		if d2 >= d0 {
			d1, i1 = d0, i0
			d0, i0 = d2, i
		} else if d2 >= d1 {
			d1, i1 = d2, i
		}
	}

	if i1 < i0 {
		i0, i1 = i1, i0
	}

	pCh0 := PolyChain(c[i0 : i1+1])
	pCh1 := append(PolyChain(c[i1:]), PolyChain(c[:i0+1])...)

	pCh0Approx, pCh1Approx := pCh0.Approximate(deviationMax), pCh1.Approximate(deviationMax)

	return Contour(append(pCh0Approx, pCh1Approx[1:len(pCh1Approx)-1]...))
}
