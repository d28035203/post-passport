// test_api_test.go — fuzzy-adventure.
// Author: d28035203

package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/d28035203/fuzzy-adventure/internal/handlers"
	"github.com/d28035203/fuzzy-adventure/internal/middleware"
	"github.com/d28035203/fuzzy-adventure/internal/models"
	"github.com/d28035203/fuzzy-adventure/internal/repo"
	"github.com/d28035203/fuzzy-adventure/internal/services"
	"github.com/d28035203/fuzzy-adventure/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TestRouter represents the test router with all middleware
var TestRouter *gin.Engine
var testDB *gorm.DB

// SetupTest initializes the test environment
func SetupTest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Load .env for JWT secret
	os.Getenv("JWT_SECRET")

	// Setup SQLite for testing (fast, no external deps)
	var err error
	testDB, err = gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Auto-migrate models
	if err := testDB.AutoMigrate(&models.User{}, &models.Post{}); err != nil {
		t.Fatalf("Failed to migrate test models: %v", err)
	}

	// Setup in-memory logger
	logger.Init("DEBUG")

	// Setup repositories
	userRepo := repo.NewUserRepository(testDB)
	postRepo := repo.NewPostRepository(testDB)

	// Setup services
	userSvc := services.NewUserService(userRepo)
	postSvc := services.NewPostService(postRepo, testDB)

	// Setup handlers
	authHandler := handlers.NewAuthHandler(userSvc)
	postHandler := handlers.NewPostHandler(postSvc)

	// Setup router with middleware
	TestRouter = gin.Default()
	TestRouter.Use(middleware.CORS())
	TestRouter.Use(middleware.Recover())
	TestRouter.Use(middleware.Logger())
	TestRouter.Use(middleware.Metrics())

	// Initialize JWT auth
	jwtAuth := middleware.NewJWTAuth("super_secret_key")

	// Auth routes
	TestRouter.POST("/register", authHandler.Register)
	TestRouter.POST("/login", authHandler.Login)
	TestRouter.GET("/profile", jwtAuth.Auth(), func(c *gin.Context) {
		user, err := authHandler.Profile(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": user})
	})

	// Post routes
	TestRouter.GET("/posts", jwtAuth.Auth(), postHandler.List)
	TestRouter.GET("/posts/:id", jwtAuth.Auth(), postHandler.Get)
	TestRouter.POST("/posts", jwtAuth.Auth(), postHandler.Create)
	TestRouter.PUT("/posts/:id", jwtAuth.Auth(), postHandler.Update)
	TestRouter.DELETE("/posts/:id", jwtAuth.Auth(), postHandler.Delete)

	// Health check
	TestRouter.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "fuzzy-adventure-api"})
	})
}

// TestHealthEndpoint tests the health check endpoint
func TestHealthEndpoint(t *testing.T) {
	SetupTest(t)
	ResetDB()

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	TestRouter.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for health check, got %d", rr.Code)
	}
	var resp map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp["status"] != "ok" {
		t.Fatalf("Expected status 'ok', got '%v'", resp["status"])
	}
	if resp["service"] != "fuzzy-adventure-api" {
		t.Fatalf("Expected service 'fuzzy-adventure-api', got '%v'", resp["service"])
	}
}

// TestUserRegistration tests user registration flow
func TestUserRegistration(t *testing.T) {
	SetupTest(t)
	ResetDB()

	// Test 1: Register a user
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(`{"username":"testuser","email":"test@example.com","password":"password123"}`)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	TestRouter.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for new user, got %d", rr.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if _, ok := resp["message"]; !ok {
		t.Fatalf("Response should contain 'message'")
	}
	if _, ok := resp["user"]; !ok {
		t.Fatalf("Response should contain 'user'")
	}
}

// TestUserLogin tests user login flow
func TestUserLogin(t *testing.T) {
	SetupTest(t)
	ResetDB()

	// First, register a user
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(`{"username":"loginuser","email":"login@example.com","password":"password123"}`)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	TestRouter.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for registration, got %d", rr.Code)
	}

	// Now login
	req, err = http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(`{"username":"loginuser","password":"password123"}`)))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	TestRouter.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for login, got %d", rr.Code)
	}

	var resp map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	// Check response contains JWT token
	if _, ok := resp["token"]; !ok {
		t.Fatalf("Response should contain 'token'")
	}
	if _, ok := resp["token_type"]; !ok {
		t.Fatalf("Response should contain 'token_type'")
	}

	// Verify JWT token is valid
	tokenString, ok := resp["token"].(string)
	if !ok {
		t.Fatal("Token should be a string")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("super_secret_key"), nil
	})
	if err != nil {
		t.Fatalf("Failed to parse JWT token: %v", err)
	}
	if !token.Valid {
		t.Fatal("JWT token should be valid")
	}
}

// TestPostCRUD tests Create, Read, Update, Delete operations
func TestPostCRUD(t *testing.T) {
	SetupTest(t)
	ResetDB()

	// First, login to get JWT token
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(`{"username":"postuser","email":"post@example.com","password":"password123"}`)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	TestRouter.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created, got %d", rr.Code)
	}

	// Login to get token
	req, err = http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(`{"username":"postuser","password":"password123"}`)))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	TestRouter.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", rr.Code)
	}

	var loginResp map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&loginResp); err != nil {
		t.Fatal(err)
	}
	tokenString, ok := loginResp["token"].(string)
	if !ok {
		t.Fatal("Token should be a string")
	}

	// Create a post
	req, err = http.NewRequest("POST", "/posts", bytes.NewBuffer([]byte(`{"title":"Test Post","content":"This is a test post content"}`)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+tokenString)

	rr = httptest.NewRecorder()
	TestRouter.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for new post, got %d", rr.Code)
	}

	var postResp map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&postResp); err != nil {
		t.Fatal(err)
	}
	postID, ok := postResp["id"].(float64)
	if !ok {
		t.Fatal("Post ID should be a number")
	}

	// Test 2: Get the post we just created
	req, err = http.NewRequest("GET", fmt.Sprintf("/posts/%d", int(postID)), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+tokenString)

	rr = httptest.NewRecorder()
	TestRouter.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for GET post, got %d", rr.Code)
	}

	var getResp map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&getResp); err != nil {
		t.Fatal(err)
	}
	if getResp["title"] != "Test Post" {
		t.Fatalf("Post title should match, got '%v'", getResp["title"])
	}

	// Test 3: List all posts
	req, err = http.NewRequest("GET", "/posts", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+tokenString)

	rr = httptest.NewRecorder()
	TestRouter.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for listing posts, got %d", rr.Code)
	}

	var listResp map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&listResp); err != nil {
		t.Fatal(err)
	}
	posts, ok := listResp["posts"].([]interface{})
	if !ok {
		t.Fatal("Posts should be a slice")
	}
	if len(posts) != 1 {
		t.Fatalf("Should have 1 post, got %d", len(posts))
	}

	// Test 4: Update the post
	updatePost := models.Post{
		Title:   "Updated Test Post",
		Content: "This is an updated test post content",
	}
	updateJSON, _ := json.Marshal(updatePost)

	req, err = http.NewRequest("PUT", fmt.Sprintf("/posts/%d", int(postID)), bytes.NewBuffer(updateJSON))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+tokenString)

	rr = httptest.NewRecorder()
	TestRouter.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for updated post, got %d", rr.Code)
	}

	var updateResp map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&updateResp); err != nil {
		t.Fatal(err)
	}
	if updateResp["title"] != "Updated Test Post" {
		t.Fatalf("Post title should be updated, got '%v'", updateResp["title"])
	}

	// Test 5: Delete the post
	req, err = http.NewRequest("DELETE", fmt.Sprintf("/posts/%d", int(postID)), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+tokenString)

	rr = httptest.NewRecorder()
	TestRouter.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for deleted post, got %d", rr.Code)
	}

	// Verify post is deleted
	req, err = http.NewRequest("GET", fmt.Sprintf("/posts/%d", int(postID)), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+tokenString)

	rr = httptest.NewRecorder()
	TestRouter.ServeHTTP(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("Expected 404 Not Found for deleted post, got %d", rr.Code)
	}
}

// TestErrorHandling tests error cases
func TestErrorHandling(t *testing.T) {
	SetupTest(t)
	ResetDB()

	// Test 1: Register the user first
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(`{"username":"loginuser","email":"login@example.com","password":"password123"}`)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	TestRouter.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created for registration, got %d", rr.Code)
	}

	// Test 2: Login with wrong password
	req, err = http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(`{"username":"loginuser","password":"wrongpassword"}`)))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	TestRouter.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401 Unauthorized for wrong password, got %d", rr.Code)
	}

	// Test 3: Login with non-existent user
	req, err = http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(`{"username":"nonexistent","password":"password123"}`)))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	TestRouter.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401 Unauthorized for non-existent user, got %d", rr.Code)
	}

	// Test 4: Register duplicate user
	req, err = http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(`{"username":"loginuser","email":"login@example.com","password":"password123"}`)))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	TestRouter.ServeHTTP(rr, req)
	if rr.Code != http.StatusConflict {
		t.Fatalf("Expected 409 Conflict for duplicate user, got %d", rr.Code)
	}
}

// TestJWTExpiry tests JWT token expiry
func TestJWTExpiry(t *testing.T) {
	SetupTest(t)
	ResetDB()

	// Register and login
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(`{"username":"expiryuser","email":"expiry@example.com","password":"password123"}`)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	TestRouter.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created, got %d", rr.Code)
	}

	req, err = http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(`{"username":"expiryuser","password":"password123"}`)))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	TestRouter.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", rr.Code)
	}

	var loginResp map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&loginResp); err != nil {
		t.Fatal(err)
	}
	tokenString, ok := loginResp["token"].(string)
	if !ok {
		t.Fatal("Token should be a string")
	}

	req, err = http.NewRequest("POST", "/posts", bytes.NewBuffer([]byte(`{"title":"Expiry Test","content":"Testing JWT expiry"}`)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+tokenString)

	rr = httptest.NewRecorder()
	TestRouter.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created, got %d", rr.Code)
	}

	var postResp map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&postResp); err != nil {
		t.Fatal(err)
	}
	postID, ok := postResp["id"].(float64)
	if !ok {
		t.Fatal("Post ID should be a number")
	}

	// Wait for token to expire (10 min + buffer)
	time.Sleep(12 * time.Minute)

	// Try to access the post with expired token
	req, err = http.NewRequest("GET", fmt.Sprintf("/posts/%d", int(postID)), nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+tokenString)

	rr = httptest.NewRecorder()
	TestRouter.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("Expected 401 Unauthorized for expired token, got %d", rr.Code)
	}
}

// TestProfileEndpoint tests user profile retrieval
func TestProfileEndpoint(t *testing.T) {
	SetupTest(t)
	ResetDB()

	// Register and login
	req, err := http.NewRequest("POST", "/register", bytes.NewBuffer([]byte(`{"username":"profileuser","email":"profile@example.com","password":"password123"}`)))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	TestRouter.ServeHTTP(rr, req)
	if rr.Code != http.StatusCreated {
		t.Fatalf("Expected 201 Created, got %d", rr.Code)
	}

	req, err = http.NewRequest("POST", "/login", bytes.NewBuffer([]byte(`{"username":"profileuser","password":"password123"}`)))
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	TestRouter.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK, got %d", rr.Code)
	}

	var loginResp map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&loginResp); err != nil {
		t.Fatal(err)
	}
	tokenString, ok := loginResp["token"].(string)
	if !ok {
		t.Fatal("Token should be a string")
	}

	// Access profile endpoint
	req, err = http.NewRequest("GET", "/profile", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+tokenString)

	rr = httptest.NewRecorder()
	TestRouter.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("Expected 200 OK for profile endpoint, got %d", rr.Code)
	}

	var profileResp map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&profileResp); err != nil {
		t.Fatal(err)
	}
	user, ok := profileResp["user"].(map[string]interface{})
	if !ok {
		t.Fatal("User should be a map")
	}
	if user["username"] != "profileuser" {
		t.Fatalf("Username should match, got '%v'", user["username"])
	}
	if user["email"] != "profile@example.com" {
		t.Fatalf("Email should match, got '%v'", user["email"])
	}
}

// ResetDB clears the in-memory database (for testing)
func ResetDB() {
	if testDB != nil {
		// Auto-migrate models to create fresh tables
		if err := testDB.AutoMigrate(&models.User{}, &models.Post{}); err != nil {
			log.Printf("Failed to migrate test models: %v", err)
		}
	}
}

// Cleanup closes the test database after all tests
func Cleanup() {
	if testDB != nil {
		if db, err := testDB.DB(); err == nil {
			db.Close()
		}
	}
}

// Run cleanup before exiting
var _ = Cleanup

