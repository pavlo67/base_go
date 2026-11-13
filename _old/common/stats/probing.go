package stats

import "golang.org/x/exp/constraints"

type Probe[TP, TV constraints.Integer | constraints.Float] struct {
	Param TP
	Value TV
}

type ProbedRange[TP, TV constraints.Integer | constraints.Float] struct {
	Min, Max Probe[TP, TV]
}
