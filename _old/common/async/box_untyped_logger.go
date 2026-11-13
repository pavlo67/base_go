package async

import (
	"fmt"
	"time"

	"github.com/pavlo67/common/common/errors"
)

// BoxUntypedLogger ------------------------------------------------------------------------------------------

var _ BoxUntyped = NewBoxUntypedLogger[int](1, 1, nil)

func NewBoxUntypedLogger[S any](c0, c1 Code, fmt func(time.Time, *S)) *BoxUntypedLogger[S] {
	return &BoxUntypedLogger[S]{c0: c0, c1: c1, fmt: fmt}
}

type BoxUntypedLogger[S any] struct {
	c0  Code
	c1  Code
	fmt func(time.Time, *S)
}

func (bul *BoxUntypedLogger[S]) Code() (_, _ Code) {
	if bul == nil {
		return 0, 0
	}

	return bul.c0, bul.c1
}

const onLoggerSend = "on BoxUntypedLogger.Set()"

func (bul *BoxUntypedLogger[S]) Set(t time.Time, values interface{}) error {
	if bul == nil {
		return errors.New("bul == nil / " + onLoggerSend)
	}

	var sensorValues *S

	switch v := values.(type) {
	case S:
		sensorValues = &v
	case *S:
		sensorValues = v
	default:
		return fmt.Errorf("wrong values: %#v / "+onLoggerSend, values)
	}

	bul.fmt(t, sensorValues)

	return nil
}

//func (ls *BoxUntypedLogger[S]) Check() *time.Time {
//	return nil
//}
