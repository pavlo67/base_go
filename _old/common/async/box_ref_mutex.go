package async

import (
	"sync"
	"time"

	"github.com/pavlo67/common/common/errors"
)

var _ Box[int] = NewBoxRefMutex[int]()

func NewBoxRefMutex[S any]() *BoxRefMutex[S] {
	return &BoxRefMutex[S]{m: &sync.Mutex{}}
}

type BoxRefMutex[S any] struct {
	v *S
	t *time.Time
	m *sync.Mutex
}

func (brm *BoxRefMutex[S]) Set(v S, t time.Time) error {
	if brm == nil {
		return errors.New("brm == nil / on BoxRefMutex.Set()")
	}

	brm.m.Lock()
	defer brm.m.Unlock()
	brm.t, brm.v = &t, &v

	return nil
}

func (brm *BoxRefMutex[S]) Get() *ValueTimed[S] {
	if brm == nil {
		return nil
	}

	brm.m.Lock()
	defer brm.m.Unlock()

	if brm.v == nil {
		return nil
	}

	return &ValueTimed[S]{*brm.v, *brm.t}

}
