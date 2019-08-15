package logger

import (
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func GetLogger(fields ...zap.Option) *zap.SugaredLogger {
	verbose, ok := viper.Get("verbose").(bool)
	if !ok {
		verbose = false
	}

	logConfig := zap.NewProductionConfig()
	logConfig.Development = false
	logConfig.Encoding = "console"
	logConfig.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	if verbose {
		logConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}
	logConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logger, _ := logConfig.Build(fields...)

	return logger.Sugar()
}
