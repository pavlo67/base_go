package logger_zap

import (
	"fmt"
	"testing"
	"time"

	"github.com/pavlo67/common/common/logger"
	"github.com/stretchr/testify/require"
)

func TestLoggerJ(t *testing.T) {

	testKey := "testKey"
	cfg := logger.Config{
		Key:         testKey,
		LogLevel:    logger.InfoLevel,
		OutputPaths: []string{fmt.Sprintf("%d.log", time.Now().Unix())},
	}

	lj, err := New(cfg)
	require.NoError(t, err)
	require.NotNil(t, lj)

	logger.TestJ(t, lj)

}

func TestLoggerFile(t *testing.T) {

	testKey := "testKey"
	cfg := logger.Config{
		Key:         testKey,
		LogLevel:    logger.InfoLevel,
		SaveFiles:   true,
		OutputPaths: []string{fmt.Sprintf("%d.log", time.Now().Unix())},
	}

	lj, err := New(cfg)
	require.NoError(t, err)
	require.NotNil(t, lj)

	lj.File(fmt.Sprintf("tracking/%04d.txt", 1), true, []byte("aaa"))

}

//func TestLoggerZapOutputPaths(t *testing.T) {
//
//	testKey := "testKey"
//
//	cfg := logger.Config{
//		Key:         testKey,
//		LogLevel:    logger.InfoLevel,
//		OutputPaths: []string{"test.log", "stdout"},
//		//ErrorOutputPaths: nil,
//		//Encoding:         "",
//	}
//	l, err := New(cfg)
//	require.NoError(t, err)
//	require.NotNil(t, l)
//
//	require.Equal(t, testKey, l.Key())
//
//	l.Comment("INITIAL BLOCK")
//
//	l.Info("IT'S TEST FOR LOGGING INFO. OK!")
//
//	l.Comment("NEXT BLOCK")
//
//	l.Error("IT'S TEST FOR LOGGING ERROR. OK!")
//}
