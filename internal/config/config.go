package config

import (
	"time"

	"github.com/rs/zerolog"
)

// Config describes the configuration for Server
type Config struct {
	LogLevel string `mapstructure:"log_level"`
	Timeouts struct {
		Startup  time.Duration
		Shutdown time.Duration
	}

	HTTP struct {
		ListenAddress string `mapstructure:"listen_address"`

		CORS struct {
			AllowedOrigins []string `mapstructure:"allowed_origins"`
		}
	}

	Database struct {
		Host             string `mapstructure:"host"`
		Name             string `mapstructure:"name"`
		Username         string `mapstructure:"username"`
		Password         string `mapstructure:"password"`
		EventsCollection string `mapstructure:"events_collection"`
	}
}

func (c *Config) GetZeroLogLevel() zerolog.Level {
	switch c.LogLevel {
	case "trace":
		return zerolog.TraceLevel
	case "disabled":
		return zerolog.Disabled
	case "panic":
		return zerolog.PanicLevel
	case "fatal":
		return zerolog.FatalLevel
	case "error":
		return zerolog.ErrorLevel
	case "warn":
		return zerolog.WarnLevel
	case "info":
		return zerolog.InfoLevel
	case "debug":
		return zerolog.DebugLevel
	default:
		return zerolog.NoLevel
	}
}
