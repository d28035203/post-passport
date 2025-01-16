// post.go — post-passport.
// Author: d28035203

package handlers

import (
	"net/http"
	"strconv"

	"github.com/d28035203/post-passport/internal/models"
	"github.com/d28035203/post-passport/internal/services"
	"github.com/gin-gonic/gin"
)

// PostHandler handles post-related operations
type PostHandler struct {
	PostService *services.PostService
}

// NewPostHandler creates a new PostHandler
func NewPostHandler(postSvc *services.PostService) *PostHandler {
	return &PostHandler{
		PostService: postSvc,
	}
}

// List returns a list of all posts
func (h *PostHandler) List(c *gin.Context) {
	posts, err := h.PostService.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

// Get returns a single post by ID
func (h *PostHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	post, err := h.PostService.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// Increment views
	h.PostService.IncrementViews(uint(id))

	c.JSON(http.StatusOK, post)
}

// Create creates a new post (protected route)
func (h *PostHandler) Create(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userIDStr := c.GetString("auth_user_id")
	if userIDStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing user ID in context"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	post.UserID = uint(userID)

	if err := h.PostService.Create(&post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, post)
}

// Update updates an existing post (protected route)
func (h *PostHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	post.ID = uint(id)

	if err := h.PostService.Update(&post); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, post)
}

// Delete deletes a post (protected route)
func (h *PostHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid post ID"})
		return
	}

	if err := h.PostService.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Post deleted successfully"})
}
