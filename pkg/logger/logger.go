package logger

import (
	"go.uber.org/zap"
)

var log *zap.Logger

func Init() error {
	var err error
	log, err = zap.NewDevelopment()
	if err != nil {
		return err
	}
	return nil
}

func Logger() *zap.Logger {
	return log
}

func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}
