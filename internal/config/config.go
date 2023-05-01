package config

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Config describes the configuration for Server
type Config struct {
	LogLevel log.Level `mapstructure:"log_level"`
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
}
