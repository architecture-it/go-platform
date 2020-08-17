package log

import (
	"os"
	"path/filepath"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Logger *zap.Logger
	// SugarLogger es un wrapper del Logger común.
	// Es más flexible ya que permite imprimir con formato, como un Printf().
	// Debido a esto, es un poco más lento que el Logger.
	SugarLogger *zap.SugaredLogger
)

func init() {
	Logger, SugarLogger = configureLogger()
}

func configureLogger() (*zap.Logger, *zap.SugaredLogger) {
	logEncoder := os.Getenv("LOG_ENCODING")
	if logEncoder == "" {
		logEncoder = "console"
	}
	config := zap.NewDevelopmentConfig()
	config.Encoding = logEncoder
	config.EncoderConfig.EncodeTime = ELKLogTimeEncoder
	if logEncoder == "json" {
		config.EncoderConfig.CallerKey = "context"
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.MessageKey = "message"
		config.EncoderConfig.LevelKey = "severity"
		config.EncoderConfig.StacktraceKey = "" //Oculto el stacktrace
	} else {
		config.EncoderConfig.EncodeLevel = ConsoleLevelEncoder
		config.EncoderConfig.EncodeCaller = ConsoleCallerEncoder
		//TODO: Descomentar con el próximo release de la librería go.uber.org/zap
		//config.EncoderConfig.ConsoleSeparator = "|"
	}
	Logger, _ = config.Build()
	if logEncoder == "json" {
		Logger = Logger.With(zap.Int("threadId", 0), zap.String("applicationName", filepath.Base(os.Args[0])))
	}
	return Logger, Logger.Sugar()
}

func ELKLogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func ConsoleLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(" | 0 | " + level.CapitalString() + " | " + filepath.Base(os.Args[0]))
}

func ConsoleCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	// TODO: consider using a byte-oriented API to save an allocation.
	enc.AppendString("| " + caller.TrimmedPath() + " |")
}

//2020-08-11 23:55:42,522 | 89 | INFO  | Andreani.Tracking.CollectorManager.exe | TrackingWorker-7 | Processing Message Worker #7
