package logger_zap

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/pavlo67/common/common"
	"github.com/pavlo67/common/common/filelib"
	"github.com/pavlo67/common/common/imagelib"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/mathlib/sets"
)

var _ logger.OperatorJ = &loggerZap{}

type loggerZap struct {
	zap.SugaredLogger
	logger.Config
	j logger.Operator
}

func New(cfg logger.Config) (logger.OperatorJ, error) {
	consoleOut, consoleErr := "stdout", "stderr"

	var err error

	cfg.OutputPaths, err = logger.ModifiedPaths(cfg.OutputPaths, cfg.BasePath, consoleOut)
	if err != nil {
		return nil, err
	}

	cfg.ErrorPaths, err = logger.ModifiedPaths(cfg.ErrorPaths, cfg.BasePath, consoleErr)
	if err != nil {
		return nil, err
	}

	if cfg.Encoding == "" {
		cfg.Encoding = "console"
	}

	if cfg.Key == "" {
		cfg.Key = time.Now().Format(time.RFC3339)[:19]
	}

	c := zap.NewProductionConfig()
	c.DisableStacktrace = true
	c.Level.SetLevel(zapLevel(cfg.LogLevel))
	c.OutputPaths = cfg.OutputPaths
	c.ErrorOutputPaths = cfg.ErrorPaths
	c.Encoding = cfg.Encoding
	c.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	c.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	l, err := c.Build()
	if err != nil {
		return nil, fmt.Errorf("can't create logger (%#v --> %#v), got %s", cfg, c, err)
	}

	fmt.Printf("LOGS WILL BE STORED INTO %v, ERROR LOGS: %v\n", c.OutputPaths, c.ErrorOutputPaths)

	return &loggerZap{SugaredLogger: *l.Sugar(), Config: cfg}, nil

}

func (op *loggerZap) J() (_ logger.Operator, outputPaths []string, _ error) {
	if op == nil {
		return nil, nil, fmt.Errorf("op == nil / on logger_zap.OperatorJ()")
	} else if op.j == nil {
		var err error
		cfg := op.Config

		var journallingPaths []string

		for _, path := range append(cfg.OutputPaths, cfg.ErrorPaths...) {
			if path == "stdin" || path == "stdout" || path == "stderr" || sets.In(journallingPaths, path) {
				continue
			}
			journallingPaths = append(journallingPaths, path)
		}
		cfg.OutputPaths, cfg.ErrorPaths = journallingPaths, nil

		c := zap.NewProductionConfig()
		c.DisableStacktrace = true
		c.Level.SetLevel(zapLevel(cfg.LogLevel))
		c.OutputPaths = cfg.OutputPaths
		c.ErrorOutputPaths = cfg.ErrorPaths
		c.Encoding = cfg.Encoding
		c.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		c.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

		l, err := c.Build()
		if err != nil {
			return nil, nil, fmt.Errorf("can't create logger (%#v --> %#v), got %s", cfg, c, err)
		}

		op.j, outputPaths = &loggerZap{SugaredLogger: *l.Sugar(), Config: cfg}, c.OutputPaths
	}

	return op.j, outputPaths, nil
}

func zapLevel(level logger.Level) zapcore.Level {
	switch level {
	case logger.TraceLevel:
		return zapcore.DebugLevel
	case logger.DebugLevel:
		return zapcore.DebugLevel
	case logger.InfoLevel:
		return zapcore.InfoLevel
	case logger.WarnLevel:
		return zapcore.WarnLevel
	case logger.ErrorLevel:
		return zapcore.ErrorLevel
	case logger.FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func (op loggerZap) Comment(text string) {
	text = "\n---------- " + text + "\n\n"

	paths := op.Config.OutputPaths
	for _, commentPath := range paths {
		if commentPath == "stdout" {
			fmt.Fprint(os.Stdout, text)
		} else if commentPath == "stderr" {
			fmt.Fprint(os.Stderr, text)
		} else if err := filelib.AppendFile(commentPath, []byte(text)); err != nil {
			op.Errorf("CAN'T WRITE COMMENT (%s) TO %s: %s", text, commentPath, err)
			continue
		}
	}
}

func (op *loggerZap) Path() string {
	if op == nil {
		return ""
	}
	return op.Config.BasePath
}

func (op *loggerZap) SetPath(basePath string) {
	if op == nil {
		return
	}
	if basePath = strings.TrimSpace(basePath); basePath == "" {
		op.Config.BasePath = ""
		return
	}

	var err error
	if basePath, err = filelib.Dir(basePath); err != nil {
		op.Errorf("can't create basePath (%s): %s / on logger.SetPath()", basePath, err)
		return
	}

	op.Config.BasePath = basePath
}

func (op loggerZap) File(path string, appending bool, data []byte) {
	if op.Config.SaveFiles {
		basedPaths, err := logger.ModifiedPaths([]string{path}, op.Config.BasePath, "")
		if err != nil {
			op.Error(err)
		} else {
			if appending {
				err = filelib.AppendFile(basedPaths[0], data)
			} else {
				filename := basedPaths[0]
				basedPath := filepath.Dir(filename)
				if _, err = filelib.Dir(basedPath); err != nil {
					op.Error(err)
				} else {
					err = os.WriteFile(filename, data, 0644)
				}
			}

			if err != nil {
				op.Errorf("CAN'T WRITE TO FILE %s: %s", path, err)
			} else if !appending {
				op.Infof("FILE IS WRITTEN  %s", basedPaths[0])
			}

		}
	}
}

func (op loggerZap) Image(path string, getImage logger.GetImage, opts common.Map) {
	if op.Config.SaveFiles {
		img, info, err := getImage.Image(opts)
		if info != "" {
			op.File(path+".txt", false, []byte(info))

			//_, filename, line, _ := runtime.Caller(1)
			//op.Infof("from %s:%d: "+info, filename, line)
		}
		if img != nil {
			basedPaths, err := logger.ModifiedPaths([]string{path}, op.Config.BasePath, "")
			if err != nil {
				op.Error(err)
			} else if err = imagelib.SavePNG(img, basedPaths[0]); err != nil {
				op.Error(err)
			}
		}
		if err != nil {
			op.Error(err)
		}
	}
}

func (op *loggerZap) SetKey(key string) {
	if op == nil {
		return
	}
	op.Config.Key = key
}

func (op loggerZap) Key() string {
	return op.Config.Key
}
