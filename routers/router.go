package routers

import (
	"errors"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/ponyjackal/eventual-consistency-blog/models"
	"github.com/ponyjackal/eventual-consistency-blog/routers/middleware"
	"github.com/ponyjackal/eventual-consistency-blog/services"
)

func SetupRoute() *gin.Engine {
	api, closeAPI := services.NewAPI()
	
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
	
	// Middleware to close the API service after request handlers
	router.Use(func(ctx *gin.Context) {
		defer closeAPI()
		ctx.Next()
	})

	// RegisterRoutes(router) // routes register

	router.POST("/post", func(ctx *gin.Context) {
		var post models.Post

		if err := ctx.BindJSON(&post); err != nil {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}

		_, err := api.NewMessage(post.Title, post.Content)
		if err != nil {
		   ctx.AbortWithError(http.StatusInternalServerError, err)
		   return
		}

		ctx.JSON(http.StatusCreated, gin.H{"message": "We have received your post and it will be published sooner or later."})
	})
	router.GET("/post/:slug", func(ctx *gin.Context) {
		post, err := api.GetPost(ctx.Param("slug"))
		if err != nil {
		   	if errors.Is(err, redis.Nil) {
				ctx.AbortWithError(http.StatusNotFound, err)
			  	return
		   	}
		   	ctx.AbortWithError(http.StatusInternalServerError, err)
		   	return
		}
		ctx.JSON(http.StatusOK, post)
	})

	return router
}