package log

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.SugaredLogger

	loggerInit sync.Once
)

func NewLogger(level zapcore.Level, out zapcore.WriteSyncer) {
	var atomicLevel = zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		out,
		atomicLevel)
	lg := zap.New(core, zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller(), zap.AddCallerSkip(1))
	logger = lg.Sugar()
}

func getLogger() *zap.SugaredLogger {
	loggerInit.Do(func() {
		if logger == nil {
			// logger is not initialized, for example, running `go test`
			NewLogger(zapcore.InfoLevel, os.Stdout)
		}
	})
	return logger
}

func Infof(template string, args ...interface{}) {
	getLogger().Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	getLogger().Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	getLogger().Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	getLogger().Fatalf(template, args...)
}
