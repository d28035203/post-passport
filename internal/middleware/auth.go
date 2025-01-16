// auth.go — post-passport.
// Author: d28035203

package middleware

import (
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
)

// JWTAuth is a JWT authentication middleware
type JWTAuth struct {
	JWTSecret string
}

// NewJWTAuth creates a new JWTAuth middleware
func NewJWTAuth(secret string) *JWTAuth {
	return &JWTAuth{JWTSecret: secret}
}

// Auth validates JWT token and sets user context
func (m *JWTAuth) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing Authorization header"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Authorization header format. Expected: Bearer <token>"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := m.parseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("auth_token", tokenString)
		c.Set("auth_username", claims["sub"])

		// Extract user ID from JWT claims
		if uid, ok := claims["uid"].(string); ok && uid != "" {
			c.Set("auth_user_id", uid)
		} else if uid, ok := claims["uid"].(float64); ok {
			c.Set("auth_user_id", strconv.Itoa(int(uid)))
		}

		c.Next()
	}
}

// parseToken parses a JWT token and returns the claims
func (m *JWTAuth) parseToken(tokenString string) (jwt.MapClaims, error) {
	claims, err := jwt.ParseWithClaims(tokenString, jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !claims.Valid {
		return nil, jwt.ErrTokenExpired
	}

	return claims.Claims.(jwt.MapClaims), nil
}

// CORS handles Cross-Origin Resource Sharing
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization, X-Request-ID")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// Recover is a panic recovery middleware
func Recover() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// Logger is a request logger middleware
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.RequestURI

		c.Next()

		// Log after response is sent
		latency := time.Since(startTime)
		status := c.Writer.Status()

		log.Printf("REQUEST: %s %s %d %s",
			c.Request.Method,
			path,
			status,
			latency,
		)
	}
}

// Metrics is a Prometheus metrics middleware
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		path := c.Request.RequestURI

		c.Next()

		// Log metrics
		latency := time.Since(startTime)
		status := c.Writer.Status()

		// In production, you would send these to Prometheus
		// metrics.Increment(path, status)
		// metrics.RecordTimer(path, latency)

		// For now, just log
		log.Printf("METRICS: %s %s %d %s",
			c.Request.Method,
			path,
			status,
			latency,
		)
	}
}
