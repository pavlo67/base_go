package async

import (
	"fmt"
	"time"

	"github.com/pavlo67/common/common/errors"
)

// BoxChan ------------------------------------------------------------------------------------------

var _ Box[int] = NewBoxChan[int](1)

func NewBoxChan[S any](length int) Box[S] {
	return &BoxChan[S]{ch: make(chan ValueTimed[S], max(length, 1))}
}

type BoxChan[S any] struct {
	ch chan ValueTimed[S]
}

const onThrow = "on BoxChan.Set()"

func (bc *BoxChan[S]) Set(s S, t time.Time) error {
	if bc == nil {
		return errors.New("bc == nil / " + onThrow)
	} else if bc.ch == nil {
		return errors.New("bc.ch == nil / " + onThrow)
	} else if len(bc.ch) >= cap(bc.ch) {
		return fmt.Errorf("len(bc.ch): %d >= cap(bc.ch): %d / "+onThrow, len(bc.ch), cap(bc.ch))
	}

	bc.ch <- ValueTimed[S]{s, t}

	return nil
}

func (bc *BoxChan[S]) Get() *ValueTimed[S] {
	if bc == nil {
		return nil
	}

	select {
	case st := <-bc.ch:
		return &st
	default:
		return nil
	}
}
