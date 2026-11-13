package async

import (
	"time"
)

// BoxUntypedStub ------------------------------------------------------------------------------------------

var _ BoxUntyped = NewBoxUntypedStub[int](1, 1)

func NewBoxUntypedStub[S any](c0, c1 Code) *BoxUntypedStub[S] {
	return &BoxUntypedStub[S]{c0: c0, c1: c1}
}

type BoxUntypedStub[S any] struct {
	c0 Code
	c1 Code
}

func (bus *BoxUntypedStub[S]) Code() (_, _ Code) {
	if bus == nil {
		return 0, 0
	}

	return bus.c0, bus.c1
}

func (bus *BoxUntypedStub[S]) Set(t time.Time, values interface{}) error {
	return nil
}

//func (sh *BoxUntypedStub[S]) Check() *time.Time {
//	return nil
//}
