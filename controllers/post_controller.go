package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ponyjackal/eventual-consistency-blog/models"
	"github.com/ponyjackal/eventual-consistency-blog/repository"
)

type PostController struct {
    postRepo *repository.PostRepository
}

func NewPostController() *PostController {
    return &PostController{
        postRepo: &repository.PostRepository{},
    }
}

func (pc *PostController) SavePost(ctx *gin.Context)  {
	var post models.Post

    if err := ctx.BindJSON(&post); err != nil {
        ctx.AbortWithError(http.StatusBadRequest, err)
        return
    }

	err := pc.postRepo.SavePost(&post);
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save post"})
		return
	}

	ctx.JSON(http.StatusCreated, post)
}


func (pc *PostController) GetPosts(ctx *gin.Context)  {
	posts, err := pc.postRepo.FindAllPosts();
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get posts"})
		return
	}

	ctx.JSON(http.StatusOK, posts)
}
