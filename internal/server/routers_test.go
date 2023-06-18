package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	t.Run("check setup router", func(t *testing.T) {
		r := SetupRouter()
		assert.Len(t, r.Handlers, 2, "should include 2 middlewares from engine")
		assert.Equal(t, r.BasePath(), "/", "base path should be /")
	})
}

func TestV2Router(t *testing.T) {
	t.Run("check V2 router defaults", func(t *testing.T) {
		s := &Server{}
		r := SetupRouter()
		v2 := r.Group("v2")
		assert.Len(t, v2.Handlers, 2, "should include 2 from router")
		assert.Equal(t, v2.BasePath(), "/v2", "base path should be v2")
		s.v2Router(v2)
		assert.Len(t, r.Routes(), 4, "v2 router should have added 4 routes to the API")
	})
}
