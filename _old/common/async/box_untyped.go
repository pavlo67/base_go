package async

import (
	"fmt"
	"sync"
	"time"

	"github.com/pavlo67/common/common/errors"
)

type Code = uint32

// implementations of BoxUntyped must be thread-safe
type BoxUntyped interface {
	Code() (_, _ Code)
	Set(_ time.Time, data interface{}) error
	// Check() *time.Time
}

// BoxUntypedGeneric ------------------------------------------------------------------------------------------

var _ BoxUntyped = NewBoxUntypedGeneric[int](1, 1)

func NewBoxUntypedGeneric[S any](c0, c1 Code) *BoxUntypedGeneric[S] {
	return &BoxUntypedGeneric[S]{c0: c0, c1: c1, m: &sync.Mutex{}}
}

type BoxUntypedGeneric[S any] struct {
	c0 Code
	c1 Code
	t  time.Time
	v  *S
	m  *sync.Mutex
}

func (bug *BoxUntypedGeneric[S]) Code() (_, _ Code) {
	if bug == nil {
		return 0, 0
	}

	return bug.c0, bug.c1
}

const onSet = "on BoxUntypedGeneric.Set()"

func (bug *BoxUntypedGeneric[S]) Set(t time.Time, values interface{}) error {
	if bug == nil {
		return errors.New("bug == nil / " + onSet)
	}

	var sensorValues *S

	switch v := values.(type) {
	case S:
		sensorValues = &v
	case *S:
		sensorValues = v
	default:
		return fmt.Errorf("wrong values: %#v / "+onSet, values)
	}

	bug.m.Lock()
	defer bug.m.Unlock()
	bug.t, bug.v = t, sensorValues

	return nil
}

func (bug *BoxUntypedGeneric[S]) SetTyped(t time.Time, values S) {
	if bug == nil {
		// TODO???
		return
	}

	bug.m.Lock()
	defer bug.m.Unlock()
	bug.t, bug.v = t, &values
}

func (bug *BoxUntypedGeneric[S]) GetTyped() (time.Time, *S) {
	if bug == nil {
		return time.Time{}, nil
	}

	bug.m.Lock()
	defer bug.m.Unlock()

	return bug.t, bug.v
}

//func (gh *BoxUntypedGeneric[S]) Check() *time.Time {
//	if gh == nil {
//		return nil
//	}
//
//	gh.m.Lock()
//	defer gh.m.Unlock()
//
//	if gh.v == nil {
//		return nil
//	}
//	return &gh.t
//}
