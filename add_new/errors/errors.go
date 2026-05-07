package errors_new

//import (
//	"fmt"
//
//	"github.com/pavlo67/common/common/errors"
//)
//
//func Errorf(format string, args ...interface{}) error {
//	return fmt.Errorf(format, args...)
//}
//
//func New(text string) error {
//	return errors.New(text)
//}
//
//func Wrapf(err error, format string, args ...interface{}) error {
//	return errors.CommonError(err, fmt.Sprintf(format, args...))
//}
//
//func Wrap(err error, msg string) error {
//	if err == nil {
//		return nil
//	}
//	return errors.CommonError(err, msg)
//}

import (
	stderrors "errors"
	"fmt"
)

type Key string

func (k Key) Error() string {
	return string(k)
}

type Error struct {
	key Key
	msg string
	err error
}

func New(key Key, msg string) error {
	return &Error{key: key, msg: msg}
}

func Newf(key Key, format string, args ...any) error {
	return &Error{
		key: key,
		msg: fmt.Sprintf(format, args...),
	}
}

func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return &Error{msg: msg, err: err}
}

func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil
	}
	return &Error{
		msg: fmt.Sprintf(format, args...),
		err: err,
	}
}

func WithKey(err error, key Key) error {
	if err == nil {
		return nil
	}
	return &Error{key: key, err: err}
}

func Append(err error, next error) error {
	if next == nil {
		return err
	}
	if err == nil {
		return next
	}
	return stderrors.Join(err, next)
}

func Is(err error, key Key) bool {
	return stderrors.Is(err, key)
}

func (e *Error) Error() string {
	if e.err == nil {
		return e.msg
	}
	if e.msg == "" {
		return e.err.Error()
	}
	return e.msg + ": " + e.err.Error()
}

func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Is(target error) bool {
	k, ok := target.(Key)
	return ok && e.key != "" && e.key == k
}

//Приклад:
//
//const (
//	ErrConfig errors.Key = "config"
//	ErrInput  errors.Key = "input"
//)
//
//func load() error {
//	err := readConfig()
//	if err != nil {
//		return errors.WithKey(
//			errors.Wrap(err, "failed to read config"),
//			ErrConfig,
//		)
//	}
//
//	return nil
//}
//
//func validate() error {
//	var err error
//
//	err = errors.Append(err, errors.WithKey(checkName(), ErrInput))
//	err = errors.Append(err, errors.WithKey(checkConfig(), ErrConfig))
//
//	return err
//}
