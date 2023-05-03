package server

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/ugcompsoc/apid/internal/config"
)

func TestStop(t *testing.T) {
	t.Run("server should stop", func(t *testing.T) {
		s := &Server{
			HTTP: &http.Server{
				Addr: ":8081",
			},
		}
		go func() {
			time.Sleep(1 * time.Second)
			err := s.Stop(context.Background())
			assert.NoError(t, err, "could not stop server")
			l, err := net.Listen("tcp", ":8081")
			assert.NoError(t, err, "expected server to not be running on port")
			if err == nil {
				l.Close()
			}
		}()
		err := s.Start(context.Background())
		assert.NoError(t, err, "expected server to start without without error")
	})
}

func TestStart(t *testing.T) {
	t.Run("server should start", func(t *testing.T) {
		s := &Server{
			HTTP: &http.Server{
				Addr: ":8082",
			},
		}
		go func() {
			time.Sleep(1 * time.Second)
			l, err := net.Listen("tcp", ":8082")
			assert.Error(t, err, "expected server to be running on port")
			if err == nil {
				l.Close()
			}
			s.Stop(context.Background())
		}()
		err := s.Start(context.Background())
		assert.NoError(t, err, "expected server to start without error")
	})

	t.Run("server should not start", func(t *testing.T) {
		s := &Server{
			HTTP: &http.Server{
				Addr: ":8083",
			},
		}
		l, err := net.Listen("tcp", ":8083")
		assert.NoError(t, err, "expected to be able to listen on port")
		defer l.Close()
		err = s.Start(context.Background())
		assert.Error(t, err, "expected server to start without error")
	})
}

// TODO lets make this better
func TestNewServer(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		var config config.Config
		config.HTTP.ListenAddress = ":8090"
		config.HTTP.CORS.AllowedOrigins = []string{"*"}
		s := NewServer(config)
		assert.IsType(t, &Server{}, s, "expected NewServer to return Server type")
	})
}
