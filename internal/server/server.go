package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/ugcompsoc/apid/internal/config"
	"github.com/ugcompsoc/apid/internal/services/database"
)

type Server struct {
	Config    config.Config
	HTTP      *http.Server
	Datastore *database.Datastore
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

	var err error
	s.Datastore, err = database.NewDatastore(&s.Config)
	if err != nil {
		log.Fatal().Err(err).Msg("database")
	}

	// root route
	r.GET("", s.RootGet)

	// v2 route
	v2 := r.Group("v2")
	s.v2Router(v2)

	// docs route
	r.GET("docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
