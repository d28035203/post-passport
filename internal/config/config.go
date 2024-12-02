// config.go — fuzzy-adventure.
// Author: d28035203

package config

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Config holds application configuration
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPass     string
	DBName     string
	RedisHost  string
	RedisPort  string
	RedisPass  string
	JWTSecret  string
	JWTExp     int
	Port       string
	LogLevel   string
}

// Load returns the configuration from environment variables
func Load() *Config {
	return &Config{
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
		DBUser:     os.Getenv("DB_USER"),
		DBPass:     os.Getenv("DB_PASS"),
		DBName:     os.Getenv("DB_NAME"),
		RedisHost:  os.Getenv("REDIS_HOST"),
		RedisPort:  os.Getenv("REDIS_PORT"),
		RedisPass:  os.Getenv("REDIS_PASS"),
		JWTSecret:  os.Getenv("JWT_SECRET"),
		JWTExp:     3600,
		Port:       os.Getenv("PORT"),
		LogLevel:   os.Getenv("LOG_LEVEL"),
	}
}

// LoadDefaultsWithEnv overrides defaults with env vars if set (for backwards compatibility)
func LoadDefaultsWithEnv() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPass:     getEnv("DB_PASS", "postgres"),
		DBName:     getEnv("DB_NAME", "fuzzy_adventure"),
		RedisHost:  getEnv("REDIS_HOST", "localhost"),
		RedisPort:  getEnv("REDIS_PORT", "6379"),
		RedisPass:  getEnv("REDIS_PASS", ""),
		JWTSecret:  getEnv("JWT_SECRET", "super_secret_key"),
		JWTExp:     3600,
		Port:       getEnv("PORT", "8080"),
		LogLevel:   getEnv("LOG_LEVEL", "INFO"),
	}
}

// getEnv returns the environment variable value or a default
func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

// ConnectDB creates a new database connection using GORM
func ConnectDB() (*gorm.DB, error) {
	cfg := Load()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort)
	db, err := gorm.Open(postgres.New(postgres.Config{DSN: dsn}), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
