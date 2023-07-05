package routers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(route *gin.Engine) {
	route.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"status": http.StatusNotFound, "message": "Route Not Found"})
	})
	route.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"live": "ok"})
	})
}