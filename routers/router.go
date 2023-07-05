package router

import (
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ponyjackal/eventual-consistency-blog/routers/middleware"
)

func SetupRoute() *gin.Engine {
	environment, _ := strconv.ParseBool(os.Getenv("DEBUG"))
	if environment {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	allowedHosts := os.Getenv("ALLOWED_HOSTS")
	router := gin.New()
	router.SetTrustedProxies([]string{allowedHosts})
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())

	RegisterRoutes(router) // routes register

	return router
}