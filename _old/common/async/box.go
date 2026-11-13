package async

import "time"

type ValueTimed[S any] struct {
	Value S
	Time  time.Time
}

type Box[S any] interface {
	Set(S, time.Time) error
	Get() *ValueTimed[S]
}
