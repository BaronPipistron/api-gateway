package logging

import (
	"github.com/BaronPipistron/api-gateway/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var (
	Logger *zap.SugaredLogger
)

func Init(cfg *config.Config) {
	if cfg.Stage.IsDev {
		initDevelopment()
		return
	}
	initProduction(cfg.Stage.LogFilePath)
}

func initDevelopment() {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.Lock(os.Stdout),
		zapcore.DebugLevel,
	)
	l := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.WarnLevel),
	)
	Logger = l.Sugar()
}

func initProduction(logFilePath string) {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    100,
		MaxBackups: 3,
		MaxAge:     14,
		Compress:   true,
	})

	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), zapcore.InfoLevel)
	fileCore := zapcore.NewCore(consoleEncoder, fileWriter, zapcore.InfoLevel)

	core := zapcore.NewTee(consoleCore, fileCore)
	l := zap.New(
		core,
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)

	Logger = l.Sugar()
}

func Info(args ...interface{}) {
	Logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	Logger.Infof(template, args...)
}

func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

func Error(args ...interface{}) {
	Logger.Error(args...)
	_ = Logger.Sync()
	os.Exit(1)
}

func Errorf(template string, args ...interface{}) {
	Logger.Errorf(template, args...)
	_ = Logger.Sync()
	os.Exit(1)
}

func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	Logger.Fatalf(template, args...)
}
