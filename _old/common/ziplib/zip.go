package ziplib

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/pkg/errors"
)

const onUnzipFile = "on ziplib.UnzipFile()"

func UnzipFile(sourceFilename, filename string) ([]byte, error) {
	r, err := zip.OpenReader(sourceFilename)
	if err != nil {
		return nil, fmt.Errorf(onUnzipFile+": can't zip.OpenReader(%s): %s", sourceFilename, err)
	}
	defer r.Close()

	for _, f := range r.File {

		if f.Name == filename || filename == "" {
			if f.FileInfo().IsDir() {
				return nil, fmt.Errorf(onUnzipFile+": requested file (%s) is a directory", filename)
			}
		}

		rc, err := f.Open()
		if err != nil {
			return nil, fmt.Errorf(onUnzipFile+": on fOpen() got %s", err)
		}

		buffer := bytes.NewBuffer(nil)
		_, err = io.Copy(buffer, rc)
		rc.Close()

		return buffer.Bytes(), nil
	}

	return nil, fmt.Errorf(onUnzipFile+": requested file (%s) isn't found", filename)
}

type ToZip struct {
	Data     []byte
	Filename string
}

const onZip = "on ziplib.Zip()"

func ZipFiles(targetFilename string, content []ToZip, perm os.FileMode) ([]byte, error) {
	buffer := bytes.NewBuffer(nil)
	zipWriter := zip.NewWriter(buffer)

	for _, c := range content {
		w1, err := zipWriter.Create(c.Filename)
		if err != nil {
			return nil, errors.Wrap(err, onZip)
		}
		if _, err := io.Copy(w1, bytes.NewBuffer(c.Data)); err != nil {
			return nil, errors.Wrap(err, onZip)
		}
	}
	zipWriter.Close()

	zippedBytes := buffer.Bytes()

	if err := os.WriteFile(targetFilename, zippedBytes, perm); err != nil {
		return nil, errors.Wrap(err, onZip)
	}

	return zippedBytes, nil
}
