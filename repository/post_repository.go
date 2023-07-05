package repository

import (
	"github.com/ponyjackal/eventual-consistency-blog/infra/database"
	"github.com/ponyjackal/eventual-consistency-blog/infra/logger"
	"github.com/ponyjackal/eventual-consistency-blog/models"
)

type PostRepository struct{}

func (pr *PostRepository) SavePost(post *models.Post) interface{} {
	err := database.DB.Create(post).Error
	if err != nil {
		logger.Errorf("error, not save data %v", err)
	}

	return err
}

func (pr *PostRepository) FindAllPosts() ([]models.Post, error) {
	var posts []models.Post

	err := database.DB.Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (pr *PostRepository) FindPostByID(id string) (*models.Post, error) {
	var post models.Post
	err := database.DB.First(&post, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (pr *PostRepository) UpdatePost(post *models.Post) error {
	err := database.DB.Save(post).Error
	return err
}