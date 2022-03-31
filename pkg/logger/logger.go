// Package logger create new logger configuration
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New logger instance
func New(appName, env, release string, buildDate, buildNumber string, gitHash string) (*zap.Logger, error) {
	var logger *zap.Logger
	var err error

	var config zap.Config

	switch env {
	case "prod":
		config = zap.NewProductionConfig()
		config.OutputPaths = []string{"stderr"}
		// config.OutputPaths = []string{logPath + appName + "-" + strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day) + ".json"}
		// config.ErrorOutputPaths = []string{logPath + strconv.Itoa(year) + "-" + strconv.Itoa(int(month)) + "-" + strconv.Itoa(day) + ".error.json"}
	case "dev":
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.OutputPaths = []string{"stderr"}
	case "debug":
		config = zap.NewProductionConfig()
		config.Encoding = "json"
		config.OutputPaths = []string{"stderr"}
	default:
		config = zap.NewDevelopmentConfig()
		config.OutputPaths = []string{"stderr"}
	}

	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.MessageKey = "message"

	logger, err = config.Build()
	if err != nil {
		return nil, err
	}

	logger = logger.With(
		zap.String("release", release),
		zap.String("appName", appName),
		zap.String("buildDate", buildDate),
		zap.String("gitHash", gitHash),
	) // .WithOptions(zap.AddCallerSkip(1))

	logger.Info("start logger",
		zap.String("buildNumber", buildNumber),
	)

	return logger, nil
}
