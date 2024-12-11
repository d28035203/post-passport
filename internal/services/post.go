// post.go — fuzzy-adventure.
// Author: d28035203

package services

import (
	"github.com/d28035203/fuzzy-adventure/internal/models"
	"github.com/d28035203/fuzzy-adventure/internal/repo"
	"gorm.io/gorm"
)

// PostService handles business logic for Post operations
type PostService struct {
	PostRepo *repo.PostRepository
	DB       *gorm.DB
}

// NewPostService creates a new PostService
func NewPostService(postRepo *repo.PostRepository, db *gorm.DB) *PostService {
	return &PostService{
		PostRepo: postRepo,
		DB:       db,
	}
}

// Create creates a new post
func (s *PostService) Create(post *models.Post) error {
	return s.PostRepo.Create(post)
}

// FindByID retrieves a post by its ID
func (s *PostService) FindByID(id uint) (*models.Post, error) {
	return s.PostRepo.FindByID(id)
}

// FindAll retrieves all posts from the database
func (s *PostService) FindAll() ([]models.Post, error) {
	return s.PostRepo.FindAll()
}

// FindByUserID retrieves all posts by a specific user
func (s *PostService) FindByUserID(userID uint) ([]models.Post, error) {
	return s.PostRepo.FindByUserID(userID)
}

// FindOrCreate checks if post exists, creates if not
func (s *PostService) FindOrCreate(post *models.Post) (*models.Post, error) {
	return s.PostRepo.FindOrCreate(post)
}

// Update updates an existing post
func (s *PostService) Update(post *models.Post) error {
	return s.PostRepo.Update(post)
}

// Delete removes a post from the database
func (s *PostService) Delete(id uint) error {
	return s.PostRepo.Delete(id)
}

// IncrementViews increments the view count for a post
func (s *PostService) IncrementViews(id uint) error {
	return s.PostRepo.IncrementViews(id)
}
