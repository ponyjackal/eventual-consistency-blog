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

func (pr *PostRepository) GetAll() ([]models.Post, error) {
	var posts []models.Post

	err := database.DB.Find(&posts).Error
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func (pr *PostRepository) GetByID(id string) (*models.Post, error) {
	var product models.Post
	err := database.DB.First(&product, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (pr *PostRepository) Update(product *models.Post) error {
	err := database.DB.Save(product).Error
	return err
}