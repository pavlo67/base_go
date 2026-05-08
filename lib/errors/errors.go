package errors

import (
	"errors"
	"fmt"
	"slices"
)

var _ error = &Error{}

type Key string

type Error struct {
	Key
	error
	stack []string
}

func (e *Error) Error() string {
	if e == nil || e.error == nil {
		return ""
	} else if len(e.stack) == 0 {
		return e.error.Error()
	}

	var err string
	for _, msg := range slices.Backward(e.stack) {
		err += msg + " / "
	}
	return err + e.error.Error()
}

func New(key Key, err any) error {
	switch v := err.(type) {
	case nil:
		return nil

	case string:
		if v == "" {
			return nil
		}
		return &Error{Key: key, error: errors.New(v)}

	case error:
		if v == nil {
			return nil
		}
		return &Error{Key: key, error: v}

	default:
		return &Error{
			Key:   key,
			error: fmt.Errorf("%v", err),
		}
	}
}

func Newf(key Key, format string, args ...any) error {
	err := fmt.Sprintf(format, args...)
	if err == "" {
		return nil
	}
	return &Error{Key: key, error: errors.New(err)}
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.error
}

func (e *Error) Is(target error) bool {
	if e == nil {
		return target == nil
	}
	return errors.Is(e.error, target)
}

func (e *Error) As(target any) bool {
	if e == nil {
		return target == nil
	} else if target == nil {
		return false
	}
	return errors.As(e.error, target)
}

func (e *Error) Join(next error) error {
	if e == nil {
		return New("", next)
	} else if next == nil {
		return e
	}
	return errors.Join(e, next)
}

func ErrorKey(err error) Key {
	if e, ok := err.(*Error); ok && e != nil {
		return e.Key
	}

	return ""
}

func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	} else if e, ok := err.(*Error); ok && e != nil {
		e.stack = append(e.stack, msg)
		return e
	}
	return &Error{stack: []string{msg}, error: err}
}

func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	msg := fmt.Sprintf(format, args...)
	if e, ok := err.(*Error); ok && e != nil {
		e.stack = append(e.stack, msg)
		return e
	}
	return &Error{stack: []string{msg}, error: err}
}
