package logger

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// logger strcut
type logger struct {
	atom  zap.AtomicLevel
	sugar *zap.SugaredLogger
}

// custom encoder
func encoder() zapcore.Encoder {
	return zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		MessageKey:  "msg",
		LevelKey:    "level",
		EncodeLevel: zapcore.CapitalLevelEncoder,
		TimeKey:     "ts",
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		},
		CallerKey:    "file",
		EncodeCaller: zapcore.ShortCallerEncoder,
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendInt64(int64(d) / 1000000)
		},
	})
}

// New create logger
func New() (*logger, error) {
	l := &logger{}

	l.atom = zap.NewAtomicLevel()

	//encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg := zap.NewDevelopmentEncoderConfig()

	zaplogger := zap.New(zapcore.NewCore(
		//zapcore.NewJSONEncoder(encoderCfg),
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		l.atom,
	), zap.AddCaller())

	l.atom.SetLevel(zap.InfoLevel)
	l.sugar = zaplogger.Sugar()

	return l, nil
}

func (l *logger) SetLevel(level string) {

	var zaplevel zapcore.Level

	switch level {
	case "debug":
		zaplevel = zap.DebugLevel
	case "info":
		zaplevel = zap.InfoLevel
	case "warn":
		zaplevel = zap.WarnLevel
	case "error":
		zaplevel = zap.ErrorLevel
	default:
		zaplevel = zap.InfoLevel
	}

	l.atom.SetLevel(zaplevel)
}

func (l *logger) Info(args interface{}) {
	l.sugar.Info(args)
}

func (l *logger) Infof(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}

func (l *logger) Warn(args interface{}) {
	l.sugar.Warn(args)
}

func (l *logger) Warnf(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}

func (l *logger) Debug(args interface{}) {
	l.sugar.Debug(args)
}

func (l *logger) Debugf(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}

func (l *logger) Error(args interface{}) {
	l.sugar.Error(args)
}

func (l *logger) Errorf(template string, args ...interface{}) {
	l.sugar.Errorf(template, args...)
}

func (l *logger) Panic(args interface{}) {
	l.sugar.Panic(args)
}

func (l *logger) Panicf(template string, args ...interface{}) {
	l.sugar.Panicf(template, args...)
}

func (l *logger) Fatal(args interface{}) {
	l.sugar.Fatal(args)
}

func (l *logger) Fatalf(template string, args ...interface{}) {
	l.sugar.Fatalf(template, args...)
}

func (l *logger) Close() error {
	return l.sugar.Sync()
}
