package requests

import (
	"github.com/gin-gonic/gin"
)

type AzanRequest struct {
	Lat      float64 `json:"lat" binding:"required,number,lte=180"`
	Lng      float64 `json:"lng" binding:"required,number"`
	TimeZone float64 `json:"timezone" binding:"omitempty"`
}

func (request *AzanRequest) Validate(ctx *gin.Context) bool {
	return Validate(ctx, request)
}
