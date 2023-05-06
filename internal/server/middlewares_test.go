package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestContextMiddleware(t *testing.T) {
	t.Run("context ID should be added to the context", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, engine := gin.CreateTestContext(w)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
		assert.NoError(t, err, "could not create http request")

		s := &Server{}
		engine.Use(s.ContextMiddleware())
		engine.ServeHTTP(w, req)
		contextId := ctx.Writer.Header().Get("context-id")
		_, err = uuid.Parse(contextId)
		assert.NoError(t, err, "expected context id to exist/get parsed correctly")
	})

	t.Run("context ID should not be added if it already exists", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, engine := gin.CreateTestContext(w)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
		assert.NoError(t, err, "could not create http request")

		s := &Server{}
		expectedContextId := uuid.New().String()
		ctx.Writer.Header().Set("context-id", expectedContextId)
		engine.Use(s.ContextMiddleware())
		engine.ServeHTTP(w, req)
		contextId := ctx.Writer.Header().Get("context-id")
		actualContextId, err := uuid.Parse(contextId)
		assert.NoError(t, err, "expected context id to exist/get parsed correctly")
		assert.Equal(t, expectedContextId, actualContextId.String(), "expected context id to equal the one previously set")
	})

	t.Run("context ID should be added if an invalid uuid is supplied", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, engine := gin.CreateTestContext(w)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
		assert.NoError(t, err, "could not create http request")

		s := &Server{}
		ctx.Writer.Header().Set("context-id", "i am an invalid uuid")
		engine.Use(s.ContextMiddleware())
		engine.ServeHTTP(w, req)
		contextId := ctx.Writer.Header().Get("context-id")
		_, err = uuid.Parse(contextId)
		assert.NoError(t, err, "expected context id to exist/get parsed correctly")
	})
}

func TestLoggingMiddleware(t *testing.T) {
	t.Run("check all of the parameters were logged", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, engine := gin.CreateTestContext(w)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/v2", nil)
		assert.NoError(t, err, "could not create http request")

		var buf bytes.Buffer
		writer := io.Writer(&buf)
		log.Logger = log.Output(writer)

		s := &Server{}
		engine.Use(s.ContextMiddleware())
		engine.Use(s.LoggingMiddleware())
		req.RemoteAddr = "127.0.0.1:80"

		engine.ServeHTTP(w, req)
		var logResult map[string]interface{}
		err = json.Unmarshal(buf.Bytes(), &logResult)
		assert.NoError(t, err, "could not unmarshal logging result to interface")
		assert.Equal(t, http.MethodGet, logResult["method"], "expected HTTP method to be logged")
		assert.Equal(t, "/v2", logResult["path"], "expected path to be logged")
		assert.Equal(t, "127.0.0.1", logResult["ip"], "expected ip to be logged")
		assert.LessOrEqual(t, float64(0), logResult["latency_ns"], "expected latency_ns to be logged")
		assert.Equal(t, float64(404), logResult["status"], "expected status to be logged")
	})

	t.Run("check X-Real-Ip header is used", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, engine := gin.CreateTestContext(w)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/v2", nil)
		assert.NoError(t, err, "could not create http request")

		var buf bytes.Buffer
		writer := io.Writer(&buf)
		log.Logger = log.Output(writer)

		s := &Server{}
		engine.Use(s.ContextMiddleware())
		engine.Use(s.LoggingMiddleware())
		req.RemoteAddr = "127.0.0.2:80"
		req.Header.Set("X-Real-Ip", "127.0.0.1")

		engine.ServeHTTP(w, req)
		var logResult map[string]interface{}
		err = json.Unmarshal(buf.Bytes(), &logResult)
		assert.NoError(t, err, "could not unmarshal logging result to interface")
		assert.Equal(t, "127.0.0.1", logResult["ip"], "expected ip to be logged")
	})
}

func TestRecoveryMiddlware(t *testing.T) {
	t.Run("recovery with valid recovered interface", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, engine := gin.CreateTestContext(w)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/recovery", nil)
		assert.NoError(t, err, "could not create http request")

		engine.Use(gin.CustomRecovery(RecoveryMiddlware))
		engine.GET("/recovery", func(c *gin.Context) {
			panic("Oh no :(")
		})

		var buf bytes.Buffer
		writer := io.Writer(&buf)
		log.Logger = log.Output(writer)

		engine.ServeHTTP(w, req)
		var logResult map[string]string
		err = json.Unmarshal(buf.Bytes(), &logResult)
		assert.NoError(t, err, "could not unmarshal logging result to interface")
		assert.Equal(t, "Oh no :(", logResult["error"], "expected panic message to be logged")
	})
}
