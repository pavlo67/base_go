package serialization

import (
	"fmt"
	"os"

	"github.com/pavlo67/common/common/filelib"
)

const onSave = "on serialization.Save()"

func Save(data interface{}, marshaler Marshaler, filename string) error {
	dataBytes, err := marshaler.Marshal(data)
	if err != nil {
		return fmt.Errorf("saving %#v got: %s / "+onSave, data, err)
	}

	if err := os.WriteFile(filename, dataBytes, 0644); err != nil {
		return fmt.Errorf("saving %s got: %s / "+onSave, filename, err)
	}

	return nil
}

const onRead = "on serialization.Read()"

func Read(filename string, marshaler Marshaler, data interface{}) error {
	dataBytes, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("reading %s got: %s / "+onRead, filename, err)
	}

	if err := marshaler.Unmarshal(dataBytes, data); err != nil {
		return fmt.Errorf("unmarshaling %s got: %s / "+onRead, dataBytes, err)
	}

	return nil
}

const onReadIfExists = "on serialization.ReadIfExists()"

func ReadIfExists(path string, marshaler Marshaler, data interface{}) (ok bool, err error) {
	exists, isDir := filelib.FileExistsAny(path)
	if !exists {
		return false, nil
	} else if isDir {
		return false, fmt.Errorf("%s is a directory / "+onReadIfExists, path)
	}

	dataBytes, err := os.ReadFile(path)
	if err != nil {
		return false, fmt.Errorf("reading %s got: %s / "+onReadIfExists, path, err)
	}

	if err := marshaler.Unmarshal(dataBytes, data); err != nil {
		return false, fmt.Errorf("unmarshaling %s got: %s / "+onReadIfExists, dataBytes, err)
	}

	return true, nil
}
