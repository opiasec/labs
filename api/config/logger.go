package config

import (
	"io"
	"os"

	"github.com/labstack/gommon/log"
	"go.uber.org/zap"
)

type ZapLogger struct {
	Sugar *zap.SugaredLogger
}

func NewZapLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	logger, err := config.Build()
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	return logger
}

func (l *ZapLogger) Output() io.Writer {
	return os.Stdout
}

func (l *ZapLogger) SetOutput(w io.Writer) {}

func (l *ZapLogger) Prefix() string {
	return ""
}

func (l *ZapLogger) SetPrefix(p string) {}

func (l *ZapLogger) Level() log.Lvl {
	return log.INFO
}

func (l *ZapLogger) SetLevel(v log.Lvl) {}

func (l *ZapLogger) Printj(j log.JSON) {
	l.Sugar.Infow("Printj", "data", j)
}

func (l *ZapLogger) Debugj(j log.JSON) {
	l.Sugar.Debugw("Debugj", "data", j)
}

func (l *ZapLogger) Infoj(j log.JSON) {
	l.Sugar.Infow("Infoj", "data", j)
}

func (l *ZapLogger) Warnj(j log.JSON) {
	l.Sugar.Warnw("Warnj", "data", j)
}

func (l *ZapLogger) Errorj(j log.JSON) {
	l.Sugar.Errorw("Errorj", "data", j)
}

func (l *ZapLogger) Fatalj(j log.JSON) {
	l.Sugar.Fatalw("Fatalj", "data", j)
}

func (l *ZapLogger) Panicj(j log.JSON) {
	l.Sugar.Panicw("Panicj", "data", j)
}

func (l *ZapLogger) Print(i ...interface{}) {
	l.Sugar.Info(i...)
}

func (l *ZapLogger) Debug(i ...interface{}) {
	l.Sugar.Debug(i...)
}

func (l *ZapLogger) Info(i ...interface{}) {
	l.Sugar.Info(i...)
}

func (l *ZapLogger) Warn(i ...interface{}) {
	l.Sugar.Warn(i...)
}

func (l *ZapLogger) Error(i ...interface{}) {
	l.Sugar.Error(i...)
}

func (l *ZapLogger) Fatal(i ...interface{}) {
	l.Sugar.Fatal(i...)
}

func (l *ZapLogger) Panic(i ...interface{}) {
	l.Sugar.Panic(i...)
}

func (l *ZapLogger) Printf(format string, v ...interface{}) {
	l.Sugar.Infof(format, v...)
}

func (l *ZapLogger) Debugf(format string, v ...interface{}) {
	l.Sugar.Debugf(format, v...)
}

func (l *ZapLogger) Infof(format string, v ...interface{}) {
	l.Sugar.Infof(format, v...)
}

func (l *ZapLogger) Warnf(format string, v ...interface{}) {
	l.Sugar.Warnf(format, v...)
}

func (l *ZapLogger) Errorf(format string, v ...interface{}) {
	l.Sugar.Errorf(format, v...)
}

func (l *ZapLogger) Fatalf(format string, v ...interface{}) {
	l.Sugar.Fatalf(format, v...)
}

func (l *ZapLogger) Panicf(format string, v ...interface{}) {
	l.Sugar.Panicf(format, v...)
}

func (l *ZapLogger) SetHeader(h string) {}
