package helpers

import (
	"errors"

	"github.com/gin-gonic/gin"
)

func RespondWithError(c *gin.Context, err error, statusCode int) {
	if err == nil || len(err.Error()) == 0 {
		err = errors.New("unknown error")
	}
	c.JSON(statusCode, gin.H{"error": err.Error()})
}

func RespondWithString(c *gin.Context, message string, statusCode int) {
	c.JSON(statusCode, gin.H{"message": message})
}
