package server

import (
	"context"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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
		assert.NoError(t, err, "expected server to start without without error")
	})
}

func TestNewServer(t *testing.T) {
	t.Run("///TODO", func(t *testing.T) {

	})
}
