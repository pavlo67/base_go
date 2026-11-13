package config

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"regexp"
	"runtime/debug"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pavlo67/common/common/cli"
	"github.com/pavlo67/common/common/filelib"
	"github.com/pavlo67/common/common/logger"
	"github.com/pavlo67/common/common/logger/logger_zap"
	"github.com/pavlo67/common/common/serialization"
)

func ShowVCSInfo() {
	if bi, ok := debug.ReadBuildInfo(); ok && bi != nil {
		for _, s := range bi.Settings {
			if ok, _ = regexp.MatchString(`^vcs\.`, s.Key); ok {
				fmt.Printf("%s\t%s\n", s.Key, s.Value)
			}
		}
		fmt.Print("\n")
	}
}

func PrepareApp(envPath, logPath, logFile string) (Envs, logger.OperatorJ) {

	rand.Seed(time.Now().UnixNano())

	// get config --------------------------------------------------------------------

	configEnv, ok := os.LookupEnv("ENV")
	if !ok {
		configEnv = "local"
	}

	envs, err := Get(filepath.Join(envPath, configEnv+".yaml"), serialization.MarshalerYAML)
	if err != nil || envs == nil {
		log.Fatalf("on PrepareApp(%s, %s) got %#v / %s", envPath, configEnv+".yaml", envs, err)
	}

	// get logger --------------------------------------------------------------------

	var saveFiles bool
	if err = envs.Value("logger_save_files", &saveFiles); err != nil {
		fmt.Fprintf(os.Stderr, cli.Red+"on PrepareApp(%s, %s), reading of 'logger_save_files' key produces the error: %s\n"+cli.Reset, envPath, configEnv+".yaml", err)
	}

	if logPath == "" {
		if err = envs.Value("logger_path", &logPath); err != nil {
			fmt.Fprintf(os.Stderr, cli.Red+"on PrepareApp(%s, %s), reading of 'logger_path' key produces the error: %s\n"+cli.Reset, envPath, configEnv+".yaml", err)
		}
	}

	logPath, err = filelib.Dir(logPath)
	if err != nil {
		log.Fatalf("on PrepareApp(%s, %s): can't create log path (%s): %s", envPath, configEnv+".yaml", logPath, err)
	}

	cfg := logger.Config{
		Key:       strconv.FormatInt(time.Now().Unix(), 10),
		BasePath:  logPath,
		SaveFiles: saveFiles,
	}

	if logFile != "" {
		cfg.OutputPaths = []string{logFile}
		cfg.ErrorPaths = []string{logFile}
	}

	l, err := logger_zap.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	return *envs, l
}

func PrepareTests(t *testing.T, envPath, logFile string) (Envs, logger.OperatorJ) {

	configEnv := "test"
	os.Setenv("ENV", configEnv)

	envsPtr, err := Get(filepath.Join(envPath, configEnv+".yaml"), serialization.MarshalerYAML)
	require.NoError(t, err)
	require.NotNil(t, envsPtr)

	cfg := logger.Config{
		Key:      strings.ReplaceAll(time.Now().Format(time.RFC3339)[:19], ":", "_"),
		LogLevel: logger.TraceLevel,
	}

	if logFile != "" {
		cfg.OutputPaths = []string{logFile}
		cfg.ErrorPaths = []string{logFile}
	}

	if err = envsPtr.Value("logger_save_files", &cfg.SaveFiles); err != nil {
		fmt.Fprintf(os.Stderr, cli.Red+"on PrepareApp(%s, %s), reading of 'logger_save_files' key produces the error: %s\n"+cli.Reset, envPath, configEnv+".yaml", err)
	}

	if logFile != "" || cfg.SaveFiles {
		var loggerPath string
		if err = envsPtr.Value("logger_path", &loggerPath); err != nil {
			fmt.Fprintf(os.Stderr, cli.Red+"on PrepareApp(%s, %s), reading of 'logger_path' key produces the error: %s\n"+cli.Reset, envPath, configEnv+".yaml", err)
		}
		cfg.BasePath, err = filelib.Dir(filepath.Join(loggerPath, cfg.Key))
		require.NoError(t, err)
	}

	// log.Fatalf("%+v", cfg)

	l, err := logger_zap.New(cfg)
	require.NoError(t, err)
	require.NotNil(t, l)

	return *envsPtr, l

}
