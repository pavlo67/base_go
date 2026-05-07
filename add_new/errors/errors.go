package errors_new

import (
	"fmt"

	"github.com/pavlo67/common/common/errors"
)

func Errorf(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

func New(text string) error {
	return errors.New(text)
}

func Wrapf(err error, format string, args ...interface{}) error {
	return errors.CommonError(err, fmt.Sprintf(format, args...))
}

func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return errors.CommonError(err, msg)
}
