package errors

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/pavlo67/common/common"
)

type Key = common.ErrorKey

type Error interface {
	error
	Cause() error
	Key() Key
	Data() common.Map
	Append(interface{}) Error
}

func CommonError(any ...interface{}) Error {
	var err *commonError
	for _, anything := range any {
		err = err.append(anything)
	}

	if reflect.ValueOf(err).IsNil() {
		return nil
	}

	return err
}

func AddErrors(err, errToAdd error) error {
	if err == nil {
		return errToAdd
	} else if v := reflect.ValueOf(err); v.Kind() == reflect.Interface && v.IsNil() {
		return errToAdd
	} else if e, _ := err.(Error); e != nil {
		return e.Append(errToAdd)
	}

	return CommonError(err, errToAdd)
}

//func CommonError(key Key, data common.Map) Error {
//	return &commonError{
//		errs: nil,
//		key:  key,
//		data: data,
//	}
//}

func Keyed(err error) Key {
	if err == nil {
		return ""
	}

	if errs, _ := err.(Error); errs != nil {
		return errs.Key()
	}
	return ""
}

func Cause(err error) error {
	if err == nil {
		return nil
	}

	if errs, _ := err.(Error); errs != nil {
		return errs.Cause()
	}
	return nil
}

func Data(err error) common.Map {
	if err == nil {
		return nil
	}

	if errs, _ := err.(Error); errs != nil {
		return errs.Data()
	}
	return nil
}

// commonError -------------------------------------------------------------------------------------------------------

var _ Error = &commonError{}

type commonError struct {
	errs multipleErrors
	key  Key
	data common.Map
}

func (ce *commonError) Cause() error {
	if ce != nil {
		if len(ce.errs) > 0 {
			return ce.errs[0]
		} else if ce.key != "" {
			// ???
			New(string(ce.key))
		}
	}

	return nil
}

func (ce *commonError) Error() string {
	if ce == nil {
		return ""
	}
	errStr := strings.TrimSpace(string(ce.key))
	if errStr == "" {
		// errStr = "<no key> "
	}

	if len(ce.data) > 0 {
		errStr += fmt.Sprintf(" (%v) ", ce.data)
	}

	return errStr + ce.errs.String()
}

func (ce *commonError) Key() Key {
	if ce == nil {
		return ""
	}

	return ce.key
}

func (ce *commonError) Data() common.Map {
	if ce == nil {
		return nil
	}

	return ce.data
}

func (ce *commonError) append(anything interface{}) *commonError {
	if anything == nil {
		return ce
	}

	switch v := anything.(type) {
	case Key:
		if ce == nil {
			return &commonError{key: v}
		} else if ce.key == "" {
			ce.key = v
			return ce
		} else {
			anything = commonError{key: v}
		}
	case common.Map:
		if ce == nil {
			return &commonError{data: v}
		} else if len(ce.data) < 1 {
			ce.data = v
			return ce
		} else {
			anything = commonError{data: v}
		}
	case string:
		anything = errors.New(v)
	}

	if ce == nil {
		switch v := anything.(type) {
		case commonError:
			v1 := v //  to prevent recursion in the case: ke1 := CommonError(...); ke2 := CommonError(ke1, ke1)
			return &v1
		case *commonError:
			v1 := *v // to prevent recursion in the case: ke1 := CommonError(...); ke2 := CommonError(ke1, ke1)
			return &v1
		case Error:
			return &commonError{
				errs: multipleErrors{errors.New(v.Error())},
				key:  v.Key(),
				data: v.Data(),
			}
		case error:
			return &commonError{errs: multipleErrors{v}}
		case string:
			return &commonError{errs: multipleErrors{errors.New(v)}}
		}
		ce = &commonError{}
	}

	var err error
	switch v := anything.(type) {
	case commonError:
		err = &v
	case *commonError:
		err = v
	case Error:
		err = v
	case error:
		err = v
	case string:
		err = errors.New(v)
	default:
	}

	ce.errs = append(ce.errs, err)
	return ce
}

func (ce *commonError) Append(anything interface{}) Error {
	return ce.append(anything)

}
