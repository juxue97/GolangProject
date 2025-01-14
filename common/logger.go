package common

import (
	"os"

	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func init() {
	Logger = newLogger()
	Logger.Info("Logger initialized. ", "Environment: ", GetString("ENV", "development"))
}

func newLogger() *zap.SugaredLogger {
	// Check environment to determine logger type
	development := os.Getenv("ENV") == "development"

	var logger *zap.Logger
	var err error

	if development {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		panic(err)
	}

	return logger.Sugar()
}
