package server

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRootV2Get(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, engine := gin.CreateTestContext(w)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/v2", new(bytes.Buffer))
		assert.NoError(t, err, "could not create http request")
		s := &Server{}
		engine.GET("/v2", s.RootV2Get)
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "expected status 200 from endpoint")
		assert.Equal(t, "{\"message\":\"Root V2\"}", w.Body.String(), "unexpected response")
	})
}

func TestBrewV2Get(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, engine := gin.CreateTestContext(w)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/v2/brew", new(bytes.Buffer))
		assert.NoError(t, err, "could not create http request")
		s := &Server{}
		engine.GET("/v2/brew", s.MiscV1BrewGet)
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusTeapot, w.Code, "expected status 418 from endpoint")
		assert.Equal(t, "{\"error\":\"I refuse to brew coffee because I am, permanently, a teapot.\"}", w.Body.String(), "unexpected response")
	})
}

func TestMiscV2Get(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, engine := gin.CreateTestContext(w)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/v2/ping", new(bytes.Buffer))
		assert.NoError(t, err, "could not create http request")
		s := &Server{}
		engine.GET("/v2/ping", s.MiscV1PingGet)
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code, "expected status 200 from endpoint")
		assert.Equal(t, "{\"message\":\"Pong!\"}", w.Body.String(), "unexpected response")
	})
}
