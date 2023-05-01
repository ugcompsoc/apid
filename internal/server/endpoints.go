package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RootV2Get(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "Root V2"})
}
