package logger

import (
	"os"
	"sync"

	"github.com/rs/zerolog"
)

var (
	once sync.Once

	logger zerolog.Logger

	DebugFlag bool = false
)

func Get() *zerolog.Logger {
	once.Do(func() {
		zerolog.TimeFieldFormat = "2006-01-02 15:04:05"
		zerolog.LevelFieldName = "level"
		zerolog.MessageFieldName = "msg"

		logLevel := zerolog.PanicLevel
		if DebugFlag {
			logLevel = zerolog.DebugLevel
		}
		zerolog.SetGlobalLevel(logLevel)

		logger = zerolog.New(os.Stdout).With().
			Timestamp().
			Caller().
			Logger()
	})
	return &logger
}
