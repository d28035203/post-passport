// user.go — post-passport.
// Author: d28035203

package services

import (
	"github.com/d28035203/post-passport/internal/models"
	"github.com/d28035203/post-passport/internal/repo"
)

// UserService handles business logic for User operations
type UserService struct {
	UserRepo *repo.UserRepository
}

// NewUserService creates a new UserService
func NewUserService(userRepo *repo.UserRepository) *UserService {
	return &UserService{UserRepo: userRepo}
}

// Create creates a new user
func (s *UserService) Create(user *models.User) error {
	return s.UserRepo.Create(user)
}

// FindByID retrieves a user by their ID
func (s *UserService) FindByID(id uint) (*models.User, error) {
	return s.UserRepo.FindByID(id)
}

// FindByUsername retrieves a user by their username
func (s *UserService) FindByUsername(username string) (*models.User, error) {
	return s.UserRepo.FindByUsername(username)
}

// FindByEmail retrieves a user by their email
func (s *UserService) FindByEmail(email string) (*models.User, error) {
	return s.UserRepo.FindByEmail(email)
}

// FindOrCreate checks if user exists, creates if not
func (s *UserService) FindOrCreate(user *models.User) (*models.User, error) {
	return s.UserRepo.FindOrCreate(user)
}

// FindAll retrieves all users from the database
func (s *UserService) FindAll() ([]models.User, error) {
	return s.UserRepo.FindAll()
}

// Delete removes a user from the database and their posts
func (s *UserService) Delete(id uint) error {
	return s.UserRepo.Delete(id)
}
