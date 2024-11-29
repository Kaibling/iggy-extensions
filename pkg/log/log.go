package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func loglevel(level string) zapcore.Level {
	var logLevel zapcore.Level

	switch level {
	case "debug":
		logLevel = zapcore.DebugLevel
	case "info":
		logLevel = zapcore.InfoLevel
	case "warn":
		logLevel = zapcore.WarnLevel
	case "error":
		logLevel = zapcore.ErrorLevel
	default:
		logLevel = zapcore.InfoLevel
	}

	return logLevel
}

func New(logLevel string, json bool) *Logger {
	var encoding string
	if json {
		encoding = "json"
	} else {
		encoding = "console"
	}

	cfg := zap.Config{
		Encoding:      encoding,
		Level:         zap.NewAtomicLevelAt(loglevel(logLevel)),
		OutputPaths:   []string{"stdout"},
		EncoderConfig: zap.NewProductionEncoderConfig(),
	}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := cfg.Build()
	if err != nil {
		// todo remove panic
		panic(err)
	}

	return &Logger{
		l:        logger,
		Fields:   map[string]zapcore.Field{},
		logLevel: logLevel,
		json:     json,
	}
}

func logCopy(log *Logger) *Logger {
	newLog := New(log.logLevel, log.json)

	for k, v := range log.Fields {
		newLog.Fields[k] = v
	}

	return newLog
}

type Logger struct {
	l        *zap.Logger
	logLevel string
	Fields   map[string]zapcore.Field
	json     bool
}

func (l *Logger) fields() []zapcore.Field {
	values := []zapcore.Field{}
	for _, value := range l.Fields {
		values = append(values, value)
	}

	return values
}

func (l *Logger) AddStringField(key string, value string) {
	l.Fields[key] = zap.String(key, value)
}

func (l *Logger) AddIntField(key string, value int) {
	l.Fields[key] = zap.Int(key, value)
}

func (l *Logger) AddAnyField(key string, value any) {
	l.Fields[key] = zap.Any(key, value)
}

func (l *Logger) ErrorMsg(msg string) {
	l.l.Error(msg, l.fields()...)
}

func (l *Logger) Error(err error) {
	l.l.Error(err.Error(), l.fields()...)
}

func (l *Logger) Debug(msg string) {
	l.l.Debug(msg, l.fields()...)
}

func (l *Logger) Warn(msg string) {
	l.l.Warn(msg, l.fields()...)
}

func (l *Logger) Info(msg string) {
	l.l.Info(msg, l.fields()...)
}

func (l *Logger) NewScope(value string) *Logger {
	newLogger := logCopy(l)
	newLogger.Fields["scope"] = zap.Any("scope", value)

	return newLogger
}
