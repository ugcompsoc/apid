package helpers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var url string

func init() {
	gin.SetMode(gin.TestMode)
	url = "https://apid.testbox.compsoc.ie"
}

func TestRespondWithError(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, engine := gin.CreateTestContext(w)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", new(bytes.Buffer))
		engine.GET("/", func(c *gin.Context) {
			RespondWithError(ctx, errors.New("testing"), http.StatusInternalServerError)
		})
		assert.NoError(t, err, "could not create http request")
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode, "expected status code 500 was not received")
		expectedRespond := gin.H{"error": "testing"}
		b, err := json.Marshal(&expectedRespond)
		assert.NoError(t, err, "expected there to be no error marshalling response")
		assert.Equal(t, string(b), w.Body.String(), "expected error message not in response")
	})

	t.Run("error message is empty string", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, engine := gin.CreateTestContext(w)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", new(bytes.Buffer))
		engine.GET("/", func(c *gin.Context) {
			RespondWithError(ctx, errors.New(""), http.StatusInternalServerError)
		})
		assert.NoError(t, err, "could not create http request")
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode, "expected status code 500 was not received")
		expectedRespond := gin.H{"error": "unknown error"}
		b, err := json.Marshal(&expectedRespond)
		assert.NoError(t, err, "expected there to be no error marshalling response")
		assert.Equal(t, string(b), w.Body.String(), "expected error message not in response")
	})

	t.Run("error is nil", func(t *testing.T) {
		w := httptest.NewRecorder()
		ctx, engine := gin.CreateTestContext(w)
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "/", new(bytes.Buffer))
		engine.GET("/", func(c *gin.Context) {
			RespondWithError(ctx, nil, http.StatusInternalServerError)
		})
		assert.NoError(t, err, "could not create http request")
		engine.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode, "expected status code 500 was not received")
		expectedRespond := gin.H{"error": "unknown error"}
		b, err := json.Marshal(&expectedRespond)
		assert.NoError(t, err, "expected there to be no error marshalling response")
		assert.Equal(t, string(b), w.Body.String(), "expected error message not in response")
	})
}