package logger

import (
	"github.com/khaled4vokalz/gourl_shortener/internal/config"
	"go.uber.org/zap"
)

type Logger interface {
	Info(msg string)
	Infow(msg string, keysAndValues ...interface{})
	Debug(msg string)
	Debugw(msg string, keysAndValues ...interface{})
	Error(msg string)
	Errorw(msg string, keysAndValues ...interface{})
	Warn(msg string)
	Warnw(msg string, keysAndValues ...interface{})
	Fatal(msg string)
	Fatalw(msg string, keysAndValues ...interface{})
	Sync()
}

type loggerImpl struct {
	log *zap.Logger
}

var singleton *loggerImpl

func GetLogger() Logger {
	if singleton != nil {
		return singleton
	}
	var zapLog *zap.Logger
	var err error

	if config.GetConfig().Environment == "dev" {
		zapLog, err = zap.NewDevelopment(
			zap.AddCallerSkip(1), // skip our Logger interface
		)
	} else {
		zapLog, err = zap.NewProduction(
			zap.AddCallerSkip(1), // skip our Logger interface
		)
	}

	if err != nil {
		panic("failed to initialize logger: " + err.Error())
	}

	singleton = &loggerImpl{log: zapLog}
	return singleton
}

func (l *loggerImpl) Info(msg string) {
	l.log.Sugar().Info(msg)
	l.log.Sugar().Level()
}
func (l *loggerImpl) Infow(msg string, keysAndValues ...interface{}) {
	l.log.Sugar().Infow(msg, keysAndValues...)
}

func (l *loggerImpl) Debug(msg string) {
	l.log.Sugar().Debug(msg)
}
func (l *loggerImpl) Debugw(msg string, keysAndValues ...interface{}) {
	l.log.Sugar().Debugw(msg, keysAndValues...)
}

func (l *loggerImpl) Error(msg string) {
	l.log.Sugar().Error(msg)
}
func (l *loggerImpl) Errorw(msg string, keysAndValues ...interface{}) {
	l.log.Sugar().Errorw(msg, keysAndValues...)
}

func (l *loggerImpl) Warn(msg string) {
	l.log.Sugar().Warn(msg)
}
func (l *loggerImpl) Warnw(msg string, keysAndValues ...interface{}) {
	l.log.Sugar().Warnw(msg, keysAndValues...)
}

func (l *loggerImpl) Fatal(msg string) {
	l.log.Sugar().Fatal(msg)
}
func (l *loggerImpl) Fatalw(msg string, keysAndValues ...interface{}) {
	l.log.Sugar().Fatalw(msg, keysAndValues...)
}

func (l *loggerImpl) Sync() {
	_ = l.log.Sync()
}
