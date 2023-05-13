package routes

import (
	"github.com/ahmedmohamed24/azan/config"
	"github.com/ahmedmohamed24/azan/controllers"
	"github.com/gin-gonic/gin"
)

func API_V1(group *gin.RouterGroup, app *config.APP) {
	group.POST("/azan", func(ctx *gin.Context) { controllers.FetchAzan(ctx, app) })
}
