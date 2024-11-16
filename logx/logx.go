package logx

import (
	"os"
	"strings"

	logrus "github.com/sirupsen/logrus"
)

var logEntry *logrus.Entry
var appLogger *logrus.Logger
var logAppLevel logrus.Level

func Init(className string, level string) {

	logger := &logrus.Logger{
		Formatter: &logrus.JSONFormatter{
			FieldMap: logrus.FieldMap{

				logrus.FieldKeyTime:  "@timestamp",
				logrus.FieldKeyLevel: "@level",
				logrus.FieldKeyMsg:   "@message",
				logrus.FieldKeyFunc:  "@caller",
			},
			TimestampFormat:  "2006-01-02 15:04:05.000Z07:00",
			CallerPrettyfier: CallerPrettyfier,
		},
		Out:   os.Stdout,
		Level: GetLogLevel(level),
	}
	logger.SetOutput(os.Stdout)

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		logger.SetLevel(logrus.InfoLevel)
	} else {
		logger.SetLevel(lvl)
	}

	entry := logger.WithFields(logrus.Fields{
		"@class": className,
	})
	logAppLevel = GetLogLevel(level)
	logEntry = entry
	appLogger = logger
}

func GetLog() *logrus.Entry {
	return logEntry

}
func GetLogger() *logrus.Logger {
	return appLogger
}
func GetLogAppLevel() logrus.Level {
	return logAppLevel
}

func GetLogLevel(logLevel string) logrus.Level {

	switch strings.ToLower(logLevel) {
	case "error":
		return logrus.ErrorLevel
	case "debug":
		return logrus.DebugLevel
	case "warn":
		return logrus.WarnLevel
	case "fatal":
		return logrus.FatalLevel
	case "trace":
		return logrus.TraceLevel
	case "info":
		return logrus.InfoLevel
	}

	return logrus.InfoLevel

}
func WithFields(fields logrus.Fields) *logrus.Entry {
	return logEntry.WithFields(fields)
}
