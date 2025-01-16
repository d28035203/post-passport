// auth.go — post-passport.
// Author: d28035203

package handlers

import (
	"net/http"
	"time"

	"github.com/d28035203/post-passport/internal/models"
	"github.com/d28035203/post-passport/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// AuthHandler handles user authentication operations
type AuthHandler struct {
	UserService *services.UserService
}

// NewAuthHandler creates a new AuthHandler
func NewAuthHandler(userSvc *services.UserService) *AuthHandler {
	return &AuthHandler{
		UserService: userSvc,
	}
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists
	if exists, err := h.UserService.FindByEmail(user.Email); err == nil && exists != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	// Create user
	if err := h.UserService.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered", "user": user})
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by username
	dbUser, err := h.UserService.FindByUsername(user.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// In production, hash password and compare
	// For now, we'll use plain text password
	if dbUser.Password != user.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	claims := jwt.MapClaims{
		"sub":  dbUser.Username,
		"uid":  dbUser.ID,
		"exp":  time.Now().Add(10 * time.Minute).Unix(),
		"iat":  time.Now().Unix(),
	}
	// JWT v5 API: Create and sign token in one step
	tokenString, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte("super_secret_key"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString, "token_type": "Bearer"})
}

// Profile returns the current user's profile
func (h *AuthHandler) Profile(c *gin.Context) (
	user *models.User,
	err error,
) {
	token := c.GetString("auth_token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing auth token"})
		return user, err
	}

	// Parse token to get username
	tokenParsed, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("super_secret_key"), nil
	})
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return user, err
	}

	// JWT v5: claims field is of type jwt.Claims (interface), cast to MapClaims
	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims type"})
		return user, err
	}

	username, ok := claims["sub"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid sub claim"})
		return user, err
	}

	// Fetch the actual user data to get email and proper ID
	user, err = h.UserService.FindByUsername(username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return user, err
	}

	return user, nil
}
