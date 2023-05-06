package config

import (
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestGetZeroLogLevel(t *testing.T) {
	t.Run("check all levels can be converted", func(t *testing.T) {
		c := &Config{}

		c.LogLevel = "trace"
		assert.Equal(t, zerolog.TraceLevel, c.GetZeroLogLevel(), "expected to return trace zerolog ENUM")

		c.LogLevel = "disabled"
		assert.Equal(t, zerolog.Disabled, c.GetZeroLogLevel(), "expected to return disabled zerolog ENUM")

		c.LogLevel = "panic"
		assert.Equal(t, zerolog.PanicLevel, c.GetZeroLogLevel(), "expected to return panic zerolog ENUM")

		c.LogLevel = "fatal"
		assert.Equal(t, zerolog.FatalLevel, c.GetZeroLogLevel(), "expected to return fatal zerolog ENUM")

		c.LogLevel = "error"
		assert.Equal(t, zerolog.ErrorLevel, c.GetZeroLogLevel(), "expected to return error zerolog ENUM")

		c.LogLevel = "warn"
		assert.Equal(t, zerolog.WarnLevel, c.GetZeroLogLevel(), "expected to return warn zerolog ENUM")

		c.LogLevel = "info"
		assert.Equal(t, zerolog.InfoLevel, c.GetZeroLogLevel(), "expected to return info zerolog ENUM")

		c.LogLevel = "debug"
		assert.Equal(t, zerolog.DebugLevel, c.GetZeroLogLevel(), "expected to return debug zerolog ENUM")

		c.LogLevel = "dummy"
		assert.Equal(t, zerolog.NoLevel, c.GetZeroLogLevel(), "expected to return nolevel zerolog ENUM")
	})
}
