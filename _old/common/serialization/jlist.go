package serialization

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/filelib"
)

const onJSONList = "on serialization.JSONList()"

func JSONList[T any](list []T, prefix, indent string) ([]byte, error) {
	prefixBytes, indentBytes := []byte(prefix), []byte(indent)

	jsonBytes := []byte{'[', '\n'}

	for i, item := range list {
		itemBytes, err := json.Marshal(item)
		if err != nil {
			return nil, errors.Wrap(err, onJSONList)
		}

		if i > 0 {
			jsonBytes = append(append(jsonBytes, ','), prefixBytes...)
		}
		jsonBytes = append(append(jsonBytes, indentBytes...), itemBytes...)
	}

	return append(jsonBytes, '\n', ']'), nil
}

const onSavePart = "on serialization.SavePart()"

func SavePart(data interface{}, marshaler Marshaler, filename string) error {

	logData, err := marshaler.Marshal(data)
	if err != nil {
		return errors.Wrap(err, onSavePart)
	} else if err = filelib.AppendFile(filename, append(logData, '\n')); err != nil {
		return errors.Wrap(err, onSavePart)
	}

	return nil
}

const onReadPart = "on serialization.ReadPart()"

func ReadPart(filename string, n int, marshaler Marshaler, data interface{}) error {
	dataBytes, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("reading %s got: %s / "+onReadPart, filename, err)
	}
	lines := bytes.Split(dataBytes, []byte{'\n'})
	if n < 0 || n >= len(lines) {
		return fmt.Errorf("wrong n (%d) to get from %d lines / "+onReadPart, n, len(lines))
	}

	if err := marshaler.Unmarshal(lines[n], data); err != nil {
		return fmt.Errorf("unmarshaling %s got: %s / "+onReadPart, dataBytes, err)
	}

	return nil
}

var reEmptyLine = regexp.MustCompile(`^\s*$`)

const onReadAllPartsJSON = "on serialization.ReadAllPartsJSON()"

func ReadAllPartsJSON(filename string, data interface{}) error {
	exists, isDir := filelib.FileExistsAny(filename)
	if !exists {
		return nil
	} else if isDir {
		return fmt.Errorf("%s is a directory / "+onReadAllPartsJSON, filename)
	}

	dataBytesRaw, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("reading %s got: %s / "+onReadAllPartsJSON, filename, err)
	}
	var lines [][]byte
	for _, line := range bytes.Split(dataBytesRaw, []byte{'\n'}) {
		if !reEmptyLine.Match(line) {
			lines = append(lines, line)
		}
	}

	dataBytesLines := bytes.Join(lines, []byte{','})
	dataBytes := append([]byte{'['}, append(dataBytesLines, ']')...)

	if err := json.Unmarshal(dataBytes, data); err != nil {
		return fmt.Errorf("unmarshaling %s got: %s / "+onReadAllPartsJSON, dataBytes, err)
	}

	return nil
}

const onSaveAllPartsJSON = "on serialization.SaveAllPartsJSON()"

func SaveAllPartsJSON[T any](data []T, filename string) error {
	var dataBytes []byte

	for _, item := range data {
		itemBytes, err := json.Marshal(item)
		if err != nil {
			return fmt.Errorf("marshaling %v got: %s / "+onSaveAllPartsJSON, item, err)
		}

		dataBytes = append(dataBytes, itemBytes...)
		dataBytes = append(dataBytes, '\n')
	}

	err := os.WriteFile(filename, dataBytes, 0644)
	if err != nil {
		return errors.Wrap(err, onSaveAllPartsJSON)
	}

	return nil
}
