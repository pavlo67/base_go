package errors

import (
	"encoding/json"
	"errors"
	"strings"
)

// multipleErrors ------------------------------------------------------------------------------------------------

type multipleErrors []error

func (errs multipleErrors) String() string {
	var errstrings []string
	for _, err := range errs {
		if err != nil {
			errstring := err.Error()
			if errstring == "" {
				errstrings = append(errstrings, "???")
			} else {
				errstrings = append(errstrings, errstring)
			}
		}
	}
	return strings.Join(errstrings, " / ")
}

func (errs multipleErrors) Append(err error) multipleErrors {
	if err != nil {
		return append(errs, err)
	}

	return errs
}

func (errs multipleErrors) AppendErrs(errsToAppend multipleErrors) multipleErrors {
	if len(errs) == 0 {
		return errsToAppend
	}

	for _, err := range errsToAppend {
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func (errs multipleErrors) Err() error {

	// TODO!!! errs.Error() must keep Keyable interface

	errstring := errs.String()
	if errstring != "" {
		return errors.New(errstring)
	}

	return nil
}

func (errs multipleErrors) MarshalJSON() ([]byte, error) {
	messages := []string{}

	for _, err := range errs {
		if err != nil {
			messages = append(messages, err.Error())
		}
	}

	return json.Marshal(messages)
}
