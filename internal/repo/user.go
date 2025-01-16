// user.go — post-passport.
// Author: d28035203

package repo

import (
	"github.com/d28035203/post-passport/internal/models"
	"gorm.io/gorm"
)

// UserRepository handles data access for User entities
type UserRepository struct {
	DB *gorm.DB
}

// NewUserRepository creates a new UserRepository instance
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

// Create inserts a new user into the database
func (r *UserRepository) Create(user *models.User) error {
	return r.DB.Create(user).Error
}

// FindByID retrieves a user by their ID
func (r *UserRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.DB.Where("id = ?", id).First(&user).Error
	return &user, err
}

// FindByUsername retrieves a user by their username
func (r *UserRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

// FindByEmail retrieves a user by their email
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("email = ?", email).First(&user).Error
	return &user, err
}

// FindAll retrieves all users from the database
func (r *UserRepository) FindAll() ([]models.User, error) {
	var users []models.User
	err := r.DB.Find(&users).Error
	return users, err
}

// FindOrCreate checks if user exists, creates if not
func (r *UserRepository) FindOrCreate(user *models.User) (*models.User, error) {
	existingUser, err := r.FindByEmail(user.Email)
	if err == nil {
		return existingUser, nil
	}

	if err != nil {
		return nil, err
	}

	return existingUser, nil
}

// Delete removes a user from the database and their posts
func (r *UserRepository) Delete(id uint) error {
	var user models.User
	err := r.DB.Where("id = ?", id).First(&user).Error
	if err != nil {
		return err
	}

	return r.DB.Delete(&user).Error
}
