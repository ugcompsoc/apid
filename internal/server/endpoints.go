package server

import (
	"context"
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
 * === MISC V2 ENDPOINTS ===
 *
 ***************************/

func (s *Server) MiscV2HealthcheckGet(c *gin.Context) {
	var errs []string = []string{}
	err := s.Datastore.Client.Ping(context.TODO(), nil)
	if err != nil {
		errs = append(errs, "cannot ping database")
	}
	c.JSON(200, gin.H{"errors": errs})
	return
}

func (s *Server) MiscV2BrewGet(c *gin.Context) {
	h.RespondWithError(c, errors.New("I refuse to brew coffee because I am, permanently, a teapot."), http.StatusTeapot)
	return
}

func (s *Server) MiscV2PingGet(c *gin.Context) {
	h.RespondWithString(c, "Pong!", http.StatusOK)
	return
}
