package filelib

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pavlo67/common/common/errors"
)

const onAppendFile = "on filelib.AppendFile()"

func AppendFile(filename string, data []byte) error {
	_, err := Dir(filepath.Dir(filename))
	if err != nil {
		return errors.Wrap(err, onAppendFile)
	}

	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return errors.Wrap(err, onAppendFile)
	}
	defer f.Close()

	n, err := f.Write(data)
	if err != nil {
		return errors.Wrap(err, onAppendFile)
	} else if n != len(data) {
		return fmt.Errorf("wrote %d bytes of %d required / "+onAppendFile, n, len(data))
	}

	return nil
}
