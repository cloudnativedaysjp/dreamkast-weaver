package logger

import (
	"context"
	"strings"

	"dreamkast-weaver/internal/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	logLevelEnv = "LOG_LEVEL"
)

type loggerKey struct{}

var (
	logLevel    zapcore.Level
	localLogger *zap.SugaredLogger
)

type config struct {
	LogLevel   zapcore.Level
	AppVersion string
	Service    string
}

func NewLogger() *zap.SugaredLogger {
	return localLogger
}

func init() {
	logLevel := utils.GetEnv(logLevelEnv, "info")
	appVersion := utils.GetEnv(utils.AppVersionEnv, "unknown")
	serviceName := utils.GetEnv(utils.ServiceNameEnv, "unknown")

	configure(config{
		LogLevel:   getZapLogLevelFromEnv(logLevel),
		AppVersion: appVersion,
		Service:    serviceName,
	})
}

func configure(config config) {
	zapConfig := defaultZapConfig()

	logger, _ := zapConfig.Build()
	fields := zap.Fields([]zap.Field{
		zap.String(utils.AppVersionKey, config.AppVersion),
		zap.String(utils.ServiceNameKey, config.Service),
	}...)

	localLogger = logger.WithOptions(fields).Sugar()
}

func defaultZapConfig() zap.Config {
	return zap.Config{
		Level:    zap.NewAtomicLevelAt(logLevel),
		Encoding: "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "severity",
			TimeKey:        "time",
			NameKey:        "logger",
			CallerKey:      "caller",
			StacktraceKey:  "stacktrace",
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
}

func FromCtx(ctx context.Context) *zap.SugaredLogger {
	logger, ok := ctx.Value(loggerKey{}).(*zap.SugaredLogger)
	if ok {
		return logger
	}

	return NewLogger()
}

func ToCtx(ctx context.Context, logger *zap.SugaredLogger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

const (
	debugStr = "debug"
	infoStr  = "info"
	warnStr  = "warn"
	errorStr = "error"
)

func getZapLogLevelFromEnv(logLevel string) zapcore.Level {
	switch strings.ToLower(logLevel) {
	case debugStr:
		return zapcore.DebugLevel
	case infoStr:
		return zapcore.InfoLevel
	case warnStr:
		return zapcore.WarnLevel
	case errorStr:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
