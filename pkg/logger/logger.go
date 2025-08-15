package logger

import (
	"member-link-lite/config"
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

// Init 初始化日志系统
func Init() {
	log = logrus.New()

	// 设置输出
	log.SetOutput(os.Stdout)

	// 设置日志级别
	level := config.GetString("log.level")
	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	// 设置日志格式
	format := config.GetString("log.format")
	if format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}
}

// GetLogger 获取日志实例
func GetLogger() *logrus.Logger {
	return log
}

// Debug 调试日志
func Debug(args ...interface{}) {
	log.Debug(args...)
}

// Info 信息日志
func Info(args ...interface{}) {
	log.Info(args...)
}

// Warn 警告日志
func Warn(args ...interface{}) {
	log.Warn(args...)
}

// Error 错误日志
func Error(args ...interface{}) {
	log.Error(args...)
}

// Fatal 致命错误日志
func Fatal(args ...interface{}) {
	log.Fatal(args...)
}

// WithField 添加字段
func WithField(key string, value interface{}) *logrus.Entry {
	return log.WithField(key, value)
}

// WithFields 添加多个字段
func WithFields(fields logrus.Fields) *logrus.Entry {
	return log.WithFields(fields)
}
