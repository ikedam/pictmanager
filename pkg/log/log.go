package log

import (
	"context"
	"fmt"

	"go.uber.org/zap"
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
	rootLogger, err = logConfig.Build()
	if err != nil {
		panic(err)
	}
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
	rootLogger.Debug(msg, fields...)
}

func Debugf(ctx context.Context, format string, args ...any) {
	Debug(ctx, fmt.Sprintf(format, args...))
}

func Info(ctx context.Context, msg string, fields ...zap.Field) {
	rootLogger.Info(msg, fields...)
}

func Error(ctx context.Context, msg string, fields ...zap.Field) {
	rootLogger.Error(msg, fields...)
}

func Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	rootLogger.Fatal(msg, fields...)
}

func Fatalf(ctx context.Context, format string, args ...any) {
	Fatal(ctx, fmt.Sprintf(format, args...))
}
