package component

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
	"github.com/stellar-payment/sp-payment/internal/config"
	"github.com/stellar-payment/sp-payment/internal/util/timeutil"
)

type NewLoggerParams struct {
	PrettyPrint bool
	ServiceName string
}

func CallerNameHook() zerolog.HookFunc {
	return func(e *zerolog.Event, l zerolog.Level, msg string) {
		pc, file, line, ok := runtime.Caller(4)
		if !ok {
			return
		}

		funcname := runtime.FuncForPC(pc).Name()
		fn := funcname[strings.LastIndex(funcname, "/")+1:]
		e.Str("caller", fn)

		if l == zerolog.ErrorLevel {
			filename := file[strings.LastIndex(file, "/")+1:]
			e.Str("file", fmt.Sprintf("%s:%d", filename, line))
		}
	}
}

func NewLogger(params NewLoggerParams) zerolog.Logger {
	var output zerolog.LevelWriter
	conf := config.Get()

	if env := os.Getenv("ENVIRONMENT"); env == string(config.EnvironmentLocal) && conf.FFJsonLogger != "1" {
		output = zerolog.MultiLevelWriter(zerolog.ConsoleWriter{
			Out: os.Stdout,
		})
	} else {
		runtimeLog, err := os.OpenFile(
			filepath.Join(conf.FilePath, "logs", fmt.Sprintf("%s.log", timeutil.FormatVerboseLogTime(conf.RunSince))),
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0664)
		if err != nil {
			panic(fmt.Errorf("failed to open logfile err: %+w", err))
		}

		output = zerolog.MultiLevelWriter(os.Stdout, runtimeLog)
	}

	return zerolog.New(output).With().Timestamp().Str("service", params.ServiceName).Logger().Hook(CallerNameHook())
}
