package plane

func (pCh PolyChain) AddProjectionPoint(pr ProjectionOnPolyChain) (PolyChain, ProjectionOnPolyChain, bool) {
	if pr.N < 0 {
		return append(PolyChain{pr.Point2}, pCh...), ProjectionOnPolyChain{Point2: pr.Point2}, true
	} else if pr.N >= len(pCh) || (pr.N == len(pCh)-1 && pr.Position > 0) {
		return append(pCh, pr.Point2), ProjectionOnPolyChain{N: len(pCh), Point2: pr.Point2}, true
	} else if pr.Position == 0 {
		// TODO??? check if pr.Point2 == pCh[pr.N]
		return pCh, pr, false
	}

	return append(pCh[:pr.N+1], append(PolyChain{pr.Point2}, pCh[pr.N+1:]...)...), ProjectionOnPolyChain{N: pr.N + 1, Point2: pr.Point2}, true
}

func (pCh PolyChain) EndProjection(start bool) ProjectionOnPolyChain {
	var pr ProjectionOnPolyChain
	if start {
		pr.Point2 = pCh[0]
	} else {
		pr.N, pr.Point2 = len(pCh)-1, pCh[len(pCh)-1]
	}

	return pr
}

func ProjectionBetween(pr0, pr1, pr ProjectionOnPolyChain) bool {
	if pr0.N > pr1.N || (pr0.N == pr1.N && pr0.Position > pr1.Position) {
		pr0, pr1 = pr1, pr0
	}

	return (pr.N > pr0.N || (pr.N == pr0.N && pr.Position >= pr0.Position)) &&
		(pr.N < pr1.N || (pr.N == pr1.N && pr.Position <= pr1.Position))
}

func (pCh PolyChain) CutWithProjection(pr ProjectionOnPolyChain) (head, tail PolyChain) {
	if pr.N < 0 {
		return nil, pCh
	} else if pr.N >= len(pCh) {
		return pCh, nil
	}

	if pr.Position == 0 {
		return pCh[:pr.N], pCh[pr.N:]
	}

	return pCh[:pr.N+1], append(PolyChain{pr.Point2}, pCh[pr.N+1:]...)
}

func (pCh PolyChain) CutWithProjections(pr0, pr1 ProjectionOnPolyChain) PolyChain {
	if pr0.N < 0 || pr0.N >= len(pCh) || pr1.N < 0 || pr1.N >= len(pCh) {
		return nil
	}

	var reversed bool
	if pr0.N > pr1.N || (pr0.N == pr1.N && pr0.Position > pr1.Position) {
		reversed, pr0, pr1 = true, pr1, pr0

	}

	if pr1.Position == 0 {
		pCh = append(PolyChain{}, pCh[:pr1.N+1]...)
	} else {
		pCh = append(append(PolyChain{}, pCh[:pr1.N+1]...), pr1.Point2)
	}

	if pr0.Position == 0 {
		pCh = append(PolyChain{}, pCh[pr0.N:]...)
	} else {
		pCh = append(PolyChain{pr0.Point2}, pCh[pr0.N+1:]...)
	}
	if reversed {
		return pCh.Reversed()
	}

	return pCh
}

//func DivideByProjection(pCh plane.PolyChain, pr plane.ProjectionOnPolyChain) []plane.PolyChain {
//	if pr.n < 0 || pr.n >= len(pCh) {
//		return nil
//	}
//
//	if pr.pos > 0 {
//		return []plane.PolyChain{
//			append(pCh[:pr.n+1], pr.Point2).Reversed(),
//			append(plane.PolyChain{pr.Point2}, pCh[pr.n+1:]...),
//		}
//	}
//
//	pChs := []plane.PolyChain{pCh[pr.n:]}
//	if pr.n > 0 {
//		pChs = append(pChs, pCh[:pr.n+1].Reversed())
//	}
//	return pChs
//}
