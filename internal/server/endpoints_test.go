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

		assert.Equal(t, http.StatusOK, ctx.Writer.Status(), "expected status 200 from endpoint")
		assert.Equal(t, "{\"message\":\"Root V2\"}", w.Body.String(), "unexpected response")
	})
}
