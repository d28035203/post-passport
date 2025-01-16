// main.go — post-passport.
// Author: d28035203

package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/d28035203/post-passport/internal/config"
	"github.com/d28035203/post-passport/internal/handlers"
	"github.com/d28035203/post-passport/internal/middleware"
	"github.com/d28035203/post-passport/internal/models"
	"github.com/d28035203/post-passport/internal/repo"
	"github.com/d28035203/post-passport/internal/services"
	"github.com/d28035203/post-passport/pkg/logger"
	"github.com/gin-gonic/gin"
)

func main() {
	// Setup config (env vars already loaded by caller or from .env in cwd)
	cfg := config.LoadDefaultsWithEnv()

	// Setup logger
	logger.Init(cfg.LogLevel)

	// Setup database
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Migrate
	db.AutoMigrate(&models.User{}, &models.Post{})
	log.Println("Database migrated successfully")

	// Setup repositories
	userRepo := repo.NewUserRepository(db)
	postRepo := repo.NewPostRepository(db)

	// Setup services
	userSvc := services.NewUserService(userRepo)
	postSvc := services.NewPostService(postRepo, db)

	// Setup handlers
	authHandler := handlers.NewAuthHandler(userSvc)
	postHandler := handlers.NewPostHandler(postSvc)

	// Setup router
	r := gin.Default()

	// Global middleware
	r.Use(middleware.CORS())
	r.Use(middleware.Recover())
	r.Use(middleware.Logger())
	r.Use(middleware.Metrics())

	// Initialize JWT auth middleware
	jwtAuth := middleware.NewJWTAuth(JWTSecret)

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "post-passport-api"})
	})

	// Auth routes
	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)
	r.GET("/profile", jwtAuth.Auth(), func(c *gin.Context) {
		user, err := authHandler.Profile(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user": user})
	})

	// Post routes (protected)
	r.GET("/posts", jwtAuth.Auth(), postHandler.List)
	r.GET("/posts/:id", jwtAuth.Auth(), postHandler.Get)
	r.POST("/posts", jwtAuth.Auth(), postHandler.Create)
	r.PUT("/posts/:id", jwtAuth.Auth(), postHandler.Update)
	r.DELETE("/posts/:id", jwtAuth.Auth(), postHandler.Delete)

	// Create server
	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	// Remove leading colon if present
	if strings.HasPrefix(port, ":") {
		port = strings.TrimPrefix(port, ":")
	}

	log.Printf("Starting server on port %s", port)
	log.Printf("Port type: %T, Port value: %q", port, port)

	err = r.Run(":" + port)
	if err != nil {
		log.Fatal(err)
	}
	err = r.Run(port)
	if err != nil {
		log.Fatal(err)
	}
}

// Graceful shutdown
func init() {
	go func() {
		sig := <-make(chan os.Signal, 2)
		log.Printf("Shutdown signal received: %v", sig)
	}()
}

// Global JWT secret
const JWTSecret = "super_secret_key_change_this_in_production"
