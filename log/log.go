package log

import (
	"encoding/json"
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
	jsonLogConfig := os.Getenv("LOG_CONFIG")
	if jsonLogConfig != "" {
		var cfg zap.Config
		if err := json.Unmarshal([]byte(jsonLogConfig), &cfg); err != nil {
			panic(err)
		}
		if cfg.Encoding == "console" { //Agrego los campos adicionales para la consola.
			cfg.EncoderConfig.EncodeLevel = ConsoleLevelEncoder
			cfg.EncoderConfig.EncodeCaller = ConsoleCallerEncoder
		} else {
			camposAdicionales := make(map[string]interface{})
			camposAdicionales["threadId"] = 0
			camposAdicionales["applicationName"] = filepath.Base(os.Args[0])
			cfg.InitialFields = camposAdicionales
			cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
			cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
		}
		cfg.EncoderConfig.EncodeTime = ELKLogTimeEncoder
		cfg.EncoderConfig.LineEnding = zapcore.DefaultLineEnding
		cfg.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder

		logger, err := cfg.Build()
		if err != nil {
			panic("Ocurió un error al crear el Logger a partir de la configuración. Revise la variable de entorno LOG_CONFIG.")
		}
		return logger, logger.Sugar()
	}

	//Creo una configuración por defecto
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeTime = ELKLogTimeEncoder
	config.EncoderConfig.EncodeLevel = ConsoleLevelEncoder
	config.EncoderConfig.EncodeCaller = ConsoleCallerEncoder
	//TODO: Descomentar con el próximo release de la librería go.uber.org/zap
	//config.EncoderConfig.ConsoleSeparator = "|"
	logger, _ := config.Build()

	return logger, logger.Sugar()
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
