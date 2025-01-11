package common

import "go.uber.org/zap"

var Logger *zap.SugaredLogger

func init() {
	Logger = newLogger()
}

func newLogger() *zap.SugaredLogger {
	logger := zap.Must(zap.NewProduction()).Sugar()
	// defer logger.Sync()

	return logger
}
