package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/cors"

	"github.com/ugcompsoc/apid/internal/config"
)

type Server struct {
	Config config.Config
	HTTP   *http.Server
}

// NewServer returns an initialized Server
func NewServer(config config.Config) *Server {
	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: config.HTTP.CORS.AllowedOrigins,
		AllowedMethods: []string{
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	r := SetupRouter()
	httpSrv := &http.Server{
		Addr:    config.HTTP.ListenAddress,
		Handler: corsMiddleware.Handler(r),
	}

	s := &Server{
		Config: config,
		HTTP:   httpSrv,
	}

	// v2 route
	v2 := r.Group("v2")
	s.v2Router(v2)

	return s
}

// Start begins listening
func (s *Server) Start(ctx context.Context) error {
	if err := s.HTTP.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to start HTTP server: %w", err)
	}
	return nil
}

// Stop shuts down the server and listener
func (s *Server) Stop(ctx context.Context) error {
	if err := s.HTTP.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to stop HTTP server: %w", err)
	}
	return nil
}
