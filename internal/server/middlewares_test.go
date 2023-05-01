package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

var logHook *test.Hook

func init() {
	gin.SetMode(gin.TestMode)
	logHook = test.NewGlobal()
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
}

func TestLoggingMiddleware(t *testing.T) {
	t.Run("check all of the parameters were logged", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, engine := gin.CreateTestContext(w)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
		assert.NoError(t, err, "could not create http request")

		s := &Server{}
		engine.Use(s.LoggingMiddleware())
		req.RemoteAddr = "127.0.0.1:80"
		engine.ServeHTTP(w, req)
		lastEntry := logHook.LastEntry().Data
		assert.Equal(t, http.MethodGet, lastEntry["method"], "expected HTTP method to be logged")
		assert.Equal(t, "/", lastEntry["path"], "expected path to be logged")
		assert.Equal(t, "127.0.0.1", lastEntry["ip"], "expected ip to be logged")
		assert.LessOrEqual(t, int64(0), lastEntry["latency_ns"], "expected latency_ns to be logged")
		assert.Equal(t, 404, lastEntry["status"], "expected status to be logged")
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
		engine.ServeHTTP(w, req)
		lastEntry := logHook.LastEntry().Data
		assert.Equal(t, "Oh no :(", lastEntry["error"], "expected panic message to be logged")
	})
}
