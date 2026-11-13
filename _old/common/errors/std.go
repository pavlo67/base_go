package errors

import (
	"errors"
	"fmt"
)

func Errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

func New(text string) error {
	return errors.New(text)
}

// DEPRECATED
func Wrapf(err error, format string, args ...interface{}) error {
	return CommonError(err, fmt.Sprintf(format, args...))
}

// DEPRECATED
func Wrap(err error, msg string) error {
	return CommonError(err, msg)
}
