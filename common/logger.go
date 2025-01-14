package common

import (
	"go.uber.org/zap"
)

func NewLogger(env string) *zap.SugaredLogger {
	var logger *zap.Logger
	var err error

	development := env == "development"

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
