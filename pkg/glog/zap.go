package glog

import (
	"GoFocusMicroService/conf"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var Log *zap.Logger

// TimeEncoder 时间编码参数
func TimeEncoder(time time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendString(time.Format("2006-01-02 15:04:05"))
}

// NewEncoderConfig 生成编码参数
func NewEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		FunctionKey:    "file",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// SetUp 实例化日志配置
func SetUp() {
	encoder := zapcore.NewConsoleEncoder(NewEncoderConfig())
	priority := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zapcore.DebugLevel
	})

	// 打印到控制台cores
	console := zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), priority)
	// 打印的日志cores
	lumberJackLogger := &lumberjack.Logger{
		Filename: conf.Conf.LogConf.FilePath,
		//MaxSize:    1,
		MaxBackups: 10,
		MaxAge:     1,
		Compress:   false,
	}
	file := zapcore.NewCore(encoder, zapcore.AddSync(lumberJackLogger), priority)
	// 2个cores添加进去
	core := zapcore.NewTee(console, file)

	// 实例化日志实例
	// https://stackoverflow.com/questions/53250323/uber-zap-logger-not-printing-caller-information-in-the-log-statement
	log := zap.New(core, zap.AddCaller())
	defer log.Sync()
	Log = log
}
