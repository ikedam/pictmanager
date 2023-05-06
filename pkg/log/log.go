package log

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logLevel   zap.AtomicLevel
	rootLogger *zap.Logger
)

func init() {
	var err error
	logLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	logConfig := zap.NewProductionConfig()
	logConfig.Level = logLevel
	logConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	rootLogger, err = logConfig.Build()
	if err != nil {
		panic(err)
	}
}

type loggerContextKeyType string

var loggerContextKey = loggerContextKeyType("logger")

func NewLogger(fields ...zap.Field) *zap.Logger {
	return rootLogger.With(fields...)
}

func CtxWithLogger(ctx context.Context, fields ...zap.Field) context.Context {
	logger := rootLogger.With(fields...)
	return context.WithValue(ctx, loggerContextKey, logger)
}

func FromCtx(ctx context.Context) *zap.Logger {
	logger := ctx.Value(loggerContextKey)
	if logger == nil {
		return rootLogger
	}
	return logger.(*zap.Logger)
}

func SetLevelByName(level string) error {
	parsedLevel, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	logLevel.SetLevel(parsedLevel.Level())
	return nil
}

func Debug(ctx context.Context, msg string, fields ...zap.Field) {
	FromCtx(ctx).Debug(msg, fields...)
}

func Debugf(ctx context.Context, format string, args ...any) {
	Debug(ctx, fmt.Sprintf(format, args...))
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	FromCtx(ctx).Info(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	FromCtx(ctx).Error(msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	FromCtx(ctx).Fatal(msg, fields...)
}

func Fatalf(ctx context.Context, format string, args ...any) {
	Fatal(ctx, fmt.Sprintf(format, args...))
}
