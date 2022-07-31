package utils

import (
	"github.com/cosmotek/loguago"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
)

func NewFileLogger(logName string) *loguago.Logger {
	zlogger := zerolog.New(&lumberjack.Logger{
		Filename:   logName,
		MaxSize:    50, // megabytes
		MaxBackups: 30,
		MaxAge:     28, // days
	}).With().Timestamp().Logger()
	return loguago.NewLogger(zlogger)
}
