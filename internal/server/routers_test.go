package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	t.Run("should include two middlewares", func(t *testing.T) {
		r := SetupRouter()
		assert.Len(t, r.Handlers, 2, "should include 2 middlewares from engine")
	})

	t.Run("base path should be /", func(t *testing.T) {
		r := SetupRouter()
		assert.Equal(t, r.BasePath(), "/", "base path should be /")
	})
}

func TestV2Router(t *testing.T) {
	t.Run("should include four middlewares", func(t *testing.T) {
		s := &Server{}
		r := SetupRouter()
		v2 := r.Group("v2")
		s.v2Router(v2)
		assert.Len(t, v2.Handlers, 4, "should include 2 middlewares from engine, 2 from router")
	})

	t.Run("base path should be v2", func(t *testing.T) {
		s := &Server{}
		r := SetupRouter()
		v2 := r.Group("v2")
		s.v2Router(v2)
		assert.Equal(t, v2.BasePath(), "/v2", "base path should be v2")
	})

	t.Run("should include one route", func(t *testing.T) {
		s := &Server{}
		r := SetupRouter()
		v2 := r.Group("v2")
		s.v2Router(v2)
		assert.Len(t, r.Routes(), 1, "v2 router should have added one route to the API")
	})
}
