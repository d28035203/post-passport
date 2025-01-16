// post.go — post-passport.
// Author: d28035203

package repo

import (
	"github.com/d28035203/post-passport/internal/models"
	"gorm.io/gorm"
)

// PostRepository handles data access for Post entities
type PostRepository struct {
	DB *gorm.DB
}

// NewPostRepository creates a new PostRepository instance
func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{DB: db}
}

// Create inserts a new post into the database
func (r *PostRepository) Create(post *models.Post) error {
	return r.DB.Create(post).Error
}

// FindByID retrieves a post by its ID
func (r *PostRepository) FindByID(id uint) (*models.Post, error) {
	var post models.Post
	err := r.DB.Where("id = ?", id).First(&post).Error
	return &post, err
}

// FindAll retrieves all posts from the database
func (r *PostRepository) FindAll() ([]models.Post, error) {
	var posts []models.Post
	err := r.DB.Find(&posts).Error
	return posts, err
}

// FindByUserID retrieves all posts by a specific user
func (r *PostRepository) FindByUserID(userID uint) ([]models.Post, error) {
	var posts []models.Post
	err := r.DB.Where("user_id = ?", userID).Find(&posts).Error
	return posts, err
}

// FindOrCreate checks if post exists, creates if not
func (r *PostRepository) FindOrCreate(post *models.Post) (*models.Post, error) {
	existingPost, err := r.FindByID(post.ID)
	if err == nil {
		return existingPost, nil
	}

	if err != nil {
		return nil, err
	}

	return existingPost, nil
}

// Update updates an existing post
func (r *PostRepository) Update(post *models.Post) error {
	return r.DB.Save(post).Error
}

// Delete removes a post from the database
func (r *PostRepository) Delete(id uint) error {
	var post models.Post
	err := r.DB.Where("id = ?", id).First(&post).Error
	if err != nil {
		return err
	}

	return r.DB.Delete(&post).Error
}

// IncrementViews increments the view count for a post
func (r *PostRepository) IncrementViews(id uint) error {
	var post models.Post
	err := r.DB.Where("id = ?", id).First(&post).Error
	if err != nil {
		return err
	}

	post.Views++
	return r.DB.Save(&post).Error
}
