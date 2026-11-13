package sets

type Pair[L any, R any] struct {
	L L
	R R
}

type Correspondence[L comparable, R comparable] []Pair[L, R]

func (corr Correspondence[L, R]) RightTo(l L) []R {
	var rs []R

	for _, pair := range corr {
		if pair.L == l && !In(rs, pair.R) {
			rs = append(rs, pair.R)
		}
	}

	return rs
}

func (corr Correspondence[L, R]) LeftTo(r R) []L {
	var ls []L

	for _, pair := range corr {
		if pair.R == r && !In(ls, pair.L) {
			ls = append(ls, pair.L)
		}
	}

	return ls
}

type MultiPair[L any, R any] struct {
	L []L
	R []R
}

type CorrespondenceClosed[L comparable, R comparable] []MultiPair[L, R]

// AlgClosure performs an algebraic closure (by transitivity) of the correspondence
func (corr Correspondence[L, R]) AlgClosure() CorrespondenceClosed[L, R] {
	if len(corr) < 1 {
		return nil
	}

	var corrLeft CorrespondenceClosed[L, R]

	for _, pair := range corr {
		accepted := false
		for i := range corrLeft {
			if corrLeft[i].L[0] == pair.L {
				corrLeft[i].R = append(corrLeft[i].R, pair.R)
				accepted = true
			}
		}
		if !accepted {
			corrLeft = append(corrLeft, MultiPair[L, R]{[]L{pair.L}, []R{pair.R}})
		}
	}

	corrClosed := corrLeft[:1]

	for _, multiPairLeft := range corrLeft[1:] {
		var accepted []int

		for i := range corrClosed {
			if Intersect(corrClosed[i].R, multiPairLeft.R) {
				accepted = append(accepted, i)
			}
		}

		if len(accepted) < 1 {
			corrClosed = append(corrClosed, multiPairLeft)
		} else {
			var corrClosedNew CorrespondenceClosed[L, R]
			multiPairJoined := multiPairLeft
			for i, multiPair := range corrClosed {
				if In(accepted, i) {
					multiPairJoined.L = Union(multiPairJoined.L, multiPair.L)
					multiPairJoined.R = Union(multiPairJoined.R, multiPair.R)
				} else {
					corrClosedNew = append(corrClosedNew, multiPair)
				}
			}
			corrClosed = append(corrClosedNew, multiPairJoined)
		}
	}

	return corrClosed
}

func (corrClosed CorrespondenceClosed[L, R]) Left() []L {
	if len(corrClosed) < 1 {
		return nil
	}

	l := corrClosed[0].L
	for _, multiPair := range corrClosed[1:] {
		l = Union(l, multiPair.L)
	}

	return l
}

func (corrClosed CorrespondenceClosed[L, R]) Right() []R {
	if len(corrClosed) < 1 {
		return nil
	}

	r := corrClosed[0].R
	for _, multiPair := range corrClosed[1:] {
		r = Union(r, multiPair.R)
	}

	return r
}
