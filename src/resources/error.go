package resources

import (
	"github.com/gin-gonic/gin"
)

func RespondWithError(ctx *gin.Context, statusCode int, err error) {
	ctx.JSON(statusCode, gin.H{
		"error": err.Error(),
	})
}
