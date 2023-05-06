package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	h "github.com/ugcompsoc/apid/internal/helpers"
)

/*
 * This middlware adds a context ID to the request so we can track all requests from this user through the logs
 */
func (s *Server) ContextMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Set the context-id field if the user hasn't given us one (or its not valid)
		_, err := uuid.Parse(ctx.Writer.Header().Get("context-id"))
		if err != nil {
			ctx.Writer.Header().Set("context-id", uuid.New().String())
		}
		// If this container is behind traefik we need to rely on the X-Real-Ip header to get the IP
		if ctx.Request.Header.Get("X-Real-Ip") == "" {
			ctx.Request.Header.Set("X-Real-Ip", ctx.RemoteIP())
		}
		ctx.Next()
	}
}

/*
 * This middleware logs primarily the request path, method, response status and completion latency
 */
func (s *Server) LoggingMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		requestLog := log.With().Str("ip", ctx.Request.Header.Get("X-Real-Ip")).Logger().
			With().Str("method", ctx.Request.Method).Logger().
			With().Str("path", ctx.Request.URL.Path).Logger().
			With().Str("context_id", ctx.Writer.Header().Get("context-id")).Logger()
		ctx.Next()
		requestLog = requestLog.With().Int64("latency_ns", time.Since(start).Nanoseconds()).Logger().
			With().Int("status", ctx.Writer.Status()).Logger()
		requestLog.Info().Msg("request_info")
	}
}

/*
 * This middleware prints a panic to the log and responds to the user with an error
 */
func RecoveryMiddlware(ctx *gin.Context, recovered interface{}) {
	log.Error().Any("error", recovered).Msg("recovery middleware")
	h.RespondWithError(ctx, errors.New("a server error was encountered"), http.StatusInternalServerError)
	return
}
