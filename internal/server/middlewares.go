package server

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	h "github.com/ugcompsoc/apid/internal/helpers"
)

/*
 * Does any random things to the context we want done before reaching the endpoint func
 */
func (s *Server) ContextMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		val := ctx.Writer.Header().Get("context-id")
		if len(val) == 0 {
			ctx.Writer.Header().Set("context-id", uuid.New().String())
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

		logFields := log.Fields{
			"method": ctx.Request.Method,
			"path":   ctx.Request.URL.Path,
			"ip":     ctx.RemoteIP(),
		}

		ctx.Next()

		logFields["latency_ns"] = time.Since(start).Nanoseconds()
		logFields["status"] = ctx.Writer.Status()
		log.WithFields(logFields).Info("request")
	}
}

/*
 * This middleware prints a panic to the log and responds to the user with an error
 */
func RecoveryMiddlware(ctx *gin.Context, recovered interface{}) {
	message, ok := recovered.(string)
	logFields := log.Fields{}
	if ok && len(message) != 0 {
		logFields["error"] = message
	}
	log.WithFields(logFields).Error("server error recovery")
	h.RespondWithError(ctx, errors.New("a server error was encountered"), http.StatusInternalServerError)
	return
}
