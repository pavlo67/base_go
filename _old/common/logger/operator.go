package logger

import (
	"image"
	"path/filepath"
	"strings"

	"github.com/pavlo67/common/common/mathlib/sets"

	"github.com/pavlo67/common/common"

	"github.com/pavlo67/common/common/errors"
	"github.com/pavlo67/common/common/filelib"

	"github.com/pavlo67/common/common/joiner"
)

const InterfaceKey joiner.InterfaceKey = "logger"

type Level int

type Config struct {
	Key         string
	LogLevel    Level
	BasePath    string
	OutputPaths []string
	ErrorPaths  []string
	Encoding    string
	SaveFiles   bool
}

const TraceLevel Level = -2
const DebugLevel Level = -1
const InfoLevel Level = 0
const WarnLevel Level = 1
const ErrorLevel Level = 2
const FatalLevel Level = 4

type GetImage interface {
	Image(opts common.Map) (image.Image, string, error)
	Bounds() image.Rectangle
}

type Operator interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})

	Info(args ...interface{})
	Infof(template string, args ...interface{})

	Warn(args ...interface{})
	Warnf(template string, args ...interface{})

	Error(args ...interface{})
	Errorf(template string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(template string, args ...interface{})

	Comment(text string)

	SetKey(key string)
	Key() string

	SetPath(basePath string)
	Path() string

	File(path string, append bool, data []byte)
	Image(path string, getImage GetImage, opts common.Map)
}

type OperatorJ interface {
	Operator

	J() (_ Operator, outputPaths []string, _ error)
}

func ModifiedPaths(paths []string, basePath, systemStream string) ([]string, error) {
	basePath = strings.TrimSpace(basePath)

	if basePath != "" {
		var err error
		if basePath, err = filelib.Dir(basePath); err != nil {
			return nil, errors.Wrapf(err, "on logger.ModifiedPaths()")
		}
	}

	var modifiedPaths []string

	for _, path := range paths {
		if path == "stdin" || path == "stdout" || path == "stderr" {
			continue
		} else if filepath.IsAbs(path) || basePath == "" {
			modifiedPaths = append(modifiedPaths, path)
		} else {
			modifiedPaths = append(modifiedPaths, filepath.Join(basePath, path))
		}
	}

	if systemStream != "" && !sets.In(modifiedPaths, systemStream) {
		return append(modifiedPaths, systemStream), nil
	}

	return modifiedPaths, nil
}
