package server

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	h "github.com/ugcompsoc/apid/internal/helpers"
)

func (s *Server) RootV2Get(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Root V2"})
}

/***************************
 *
 * === MISC V1 ENDPOINTS ===
 *
 ***************************/

func (s *Server) MiscV1BrewGet(c *gin.Context) {
	h.RespondWithError(c, errors.New("I refuse to brew coffee because I am, permanently, a teapot."), http.StatusTeapot)
	return
}

func (s *Server) MiscV1PingGet(c *gin.Context) {
	h.RespondWithString(c, "Pong!", 200)
	return
}
