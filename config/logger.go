package config

import "github.com/gofiber/fiber/v3/log"

type LoggerConfig struct {
	Level string
}

func NewLoggerConfig() *LoggerConfig {
	return &LoggerConfig{
		Level: GetEnv("LOG_LEVEL", "info"),
	}
}

func (l *LoggerConfig) Setup() {
	switch l.Level {
	case "debug":
		log.SetLevel(log.LevelDebug)
	case "warn":
		log.SetLevel(log.LevelWarn)
	case "error":
		log.SetLevel(log.LevelError)
	default:
		log.SetLevel(log.LevelInfo)
	}
}
