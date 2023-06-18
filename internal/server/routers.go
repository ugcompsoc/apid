package server

import "github.com/gin-gonic/gin"

// Returns the routes associated with /v2
func (s *Server) v2Router(r *gin.RouterGroup) {
	r.Use(s.ContextMiddleware())
	r.Use(s.LoggingMiddleware())

	r.GET("/", s.RootV2Get)
	r.GET("/healthcheck", s.MiscV2HealthcheckGet)
	r.GET("/brew", s.MiscV2BrewGet)
	r.GET("/ping", s.MiscV2PingGet)
}

func SetupRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.CustomRecovery(RecoveryMiddlware))
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		}
		c.Next()
	})
	return r
}
