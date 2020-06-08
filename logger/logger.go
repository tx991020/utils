package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger, _ = zap.NewDevelopment(zap.AddCallerSkip(1))
var sugar = logger.Sugar()

// Debug outputs a message at debug level.
// This call is a wrapper around [Logger.Debug](https://godoc.org/go.uber.org/zap#Logger.Debug)
func Debug(msg string, fields ...zapcore.Field) {
	logger.Debug(msg, fields...)
}

// Debuga uses fmt.Sprint to construct and log a message at debug level.
// This call is a wrapper around [Sugaredlogger.Debug](https://godoc.org/go.uber.org/zap#Sugaredlogger.Debug)
func Debuga(args ...interface{}) {
	sugar.Debug(args...)
}

// Debugf uses fmt.Sprintf to construct and log a message at debug level.
// This call is a wrapper around [Sugaredlogger.Debugf](https://godoc.org/go.uber.org/zap#Sugaredlogger.Debugf)
func Debugf(template string, args ...interface{}) {
	sugar.Debugf(template, args...)
}

// Debugw logs a message at debug level with some additional context.
// This call is a wrapper around [Sugaredlogger.Debugw](https://godoc.org/go.uber.org/zap#Sugaredlogger.Debugw)
func Debugw(msg string, keysAndValues ...interface{}) {
	sugar.Debugw(msg, keysAndValues...)
}

// DebugEnabled returns whether output of messages at the debug level is currently enabled.
func DebugEnabled() bool {
	return logger.Core().Enabled(zap.DebugLevel)
}

// Error outputs a message at error level.
// This call is a wrapper around [logger.Error](https://godoc.org/go.uber.org/zap#logger.Error)
func Error(msg string, fields ...zapcore.Field) {
	logger.Error(msg, fields...)
}

// Errora uses fmt.Sprint to construct and log a message at error level.
// This call is a wrapper around [Sugaredlogger.Error](https://godoc.org/go.uber.org/zap#Sugaredlogger.Error)
func Errora(args ...interface{}) {
	sugar.Error(args...)
}

// Errorf uses fmt.Sprintf to construct and log a message at error level.
// This call is a wrapper around [Sugaredlogger.Errorf](https://godoc.org/go.uber.org/zap#Sugaredlogger.Errorf)
func Errorf(template string, args ...interface{}) {
	sugar.Errorf(template, args...)
}

// Errorw logs a message at error level with some additional context.
// This call is a wrapper around [Sugaredlogger.Errorw](https://godoc.org/go.uber.org/zap#Sugaredlogger.Errorw)
func Errorw(msg string, keysAndValues ...interface{}) {
	sugar.Errorw(msg, keysAndValues...)
}

// ErrorEnabled returns whether output of messages at the error level is currently enabled.
func ErrorEnabled() bool {
	return logger.Core().Enabled(zap.ErrorLevel)
}

// Warn outputs a message at warn level.
// This call is a wrapper around [logger.Warn](https://godoc.org/go.uber.org/zap#logger.Warn)
func Warn(msg string, fields ...zapcore.Field) {
	logger.Warn(msg, fields...)
}

// Warna uses fmt.Sprint to construct and log a message at warn level.
// This call is a wrapper around [Sugaredlogger.Warn](https://godoc.org/go.uber.org/zap#Sugaredlogger.Warn)
func Warna(args ...interface{}) {
	sugar.Warn(args...)
}

// Warnf uses fmt.Sprintf to construct and log a message at warn level.
// This call is a wrapper around [Sugaredlogger.Warnf](https://godoc.org/go.uber.org/zap#Sugaredlogger.Warnf)
func Warnf(template string, args ...interface{}) {
	sugar.Warnf(template, args...)
}

// Warnw logs a message at warn level with some additional context.
// This call is a wrapper around [Sugaredlogger.Warnw](https://godoc.org/go.uber.org/zap#Sugaredlogger.Warnw)
func Warnw(msg string, keysAndValues ...interface{}) {
	sugar.Warnw(msg, keysAndValues...)
}

// WarnEnabled returns whether output of messages at the warn level is currently enabled.
func WarnEnabled() bool {
	return logger.Core().Enabled(zap.WarnLevel)
}

// Info outputs a message at information level.
// This call is a wrapper around [logger.Info](https://godoc.org/go.uber.org/zap#logger.Info)
func Info(msg string, fields ...zapcore.Field) {
	logger.Info(msg, fields...)
}

// Infoa uses fmt.Sprint to construct and log a message at info level.
// This call is a wrapper around [Sugaredlogger.Info](https://godoc.org/go.uber.org/zap#Sugaredlogger.Info)
func Infoa(args ...interface{}) {
	sugar.Info(args...)
}

// Infof uses fmt.Sprintf to construct and log a message at info level.
// This call is a wrapper around [Sugaredlogger.Infof](https://godoc.org/go.uber.org/zap#Sugaredlogger.Infof)
func Infof(template string, args ...interface{}) {
	sugar.Infof(template, args...)
}

// Infow logs a message at info level with some additional context.
// This call is a wrapper around [Sugaredlogger.Infow](https://godoc.org/go.uber.org/zap#Sugaredlogger.Infow)
func Infow(msg string, keysAndValues ...interface{}) {
	sugar.Infow(msg, keysAndValues...)
}

// InfoEnabled returns whether output of messages at the info level is currently enabled.
func InfoEnabled() bool {
	return logger.Core().Enabled(zap.InfoLevel)
}

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
// This call is a wrapper around [logger.With](https://godoc.org/go.uber.org/zap#logger.With)
func With(fields ...zapcore.Field) *zap.Logger {
	return logger.With(fields...)
}

// Sync flushes any buffered log entries.
// Processes should normally take care to call Sync before exiting.
// This call is a wrapper around [logger.Sync](https://godoc.org/go.uber.org/zap#logger.Sync)
func Sync() error {
	return logger.Sync()
}
