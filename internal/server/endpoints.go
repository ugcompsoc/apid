package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ugcompsoc/apid/internal/helpers"
)

// RootGet					godoc
// @Summary					Redirect to swagger docs
// @Description				Redirect to swagger docs
// @Tags					Root
// @Success					307
// @Header					307	{string}	Location	"docs/index.html"
// @Router					/ [get]
func (s *Server) RootGet(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, "docs/index.html")
}

// MiscV2HealthcheckGet		godoc
// @Summary					Get health of API
// @Description				Responds with 'Root V2' message
// @Tags					V2
// @Produce					json
// @Success					200	{object}	helpers.Message
// @Router					/v2	[get]
func (s *Server) RootV2Get(c *gin.Context) {
	c.JSON(http.StatusOK, helpers.Message{Message: "Root V2"})
}

/***************************
 *
 * === MISC V2 ENDPOINTS ===
 *
 ***************************/

// MiscV2HealthcheckGet		godoc
// @Summary					Get health of API
// @Description				Responds with any service errors
// @Tags					V2
// @Produce					json
// @Success					200	{object}	helpers.ErrorsArray
// @Failure					500	{object}	helpers.Error
// @Router					/v2/healthcheck [get]
func (s *Server) MiscV2HealthcheckGet(c *gin.Context) {
	var errs []string = []string{}
	err := s.Datastore.Client.Ping(context.TODO(), nil)
	if err != nil {
		errs = append(errs, "cannot ping database")
	}
	c.JSON(200, helpers.ErrorsArray{Errors: errs})
	return
}

// MiscV2BrewGet			godoc
// @Summary					Brew coffee
// @Description				Responds with refusal to brew coffee
// @Tags					V2
// @Produce					json
// @Success					200	{object}	helpers.Error
// @Router					/v2/brew [get]
func (s *Server) MiscV2BrewGet(c *gin.Context) {
	helpers.RespondWithError(c, errors.New("I refuse to brew coffee because I am, permanently, a teapot."), http.StatusTeapot)
	return
}

// MiscV2PingGet			godoc
// @Summary					Ping pong
// @Description				Responds with a pong
// @Tags					V2
// @Produce					json
// @Success					200	{object}	helpers.Message
// @Router					/v2/ping [get]
func (s *Server) MiscV2PingGet(c *gin.Context) {
	helpers.RespondWithString(c, "Pong!", http.StatusOK)
	return
}
