package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server  ServerConfig
	MongoDB MongoDBConfig
	Redis   RedisConfig
	JWT     JWTConfig
	API     APIConfig
	Swagger SwaggerConfig
}

type SwaggerConfig struct {
	AllowWithoutAuth bool
	Username         string
	Password         string
}

type ServerConfig struct {
	Port string
}

type MongoDBConfig struct {
	URI      string
	Database string
}

type RedisConfig struct {
	URI string
}

type JWTConfig struct {
	Secret     string
	Expiration time.Duration
}

type APIConfig struct {
	Timeout time.Duration
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
		},
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGODB_URI", "mongodb://localhost:27017"),
			Database: getEnv("MONGODB_DATABASE", "trading_chats"),
		},
		Redis: RedisConfig{
			URI: getEnv("REDIS_URI", "redis://localhost:6379/0"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "change_me_in_production"),
			Expiration: getDurationEnv("JWT_EXPIRATION", 24*time.Hour),
		},
		API: APIConfig{
			Timeout: getDurationEnv("API_TIMEOUT", 60*time.Second),
		},
		Swagger: SwaggerConfig{
			AllowWithoutAuth: getEnv("SWAGGER_ALLOW_WITHOUT_AUTH", "true") == "true",
			Username:         getEnv("SWAGGER_USERNAME", "admin"),
			Password:         getEnv("SWAGGER_PASSWORD", "swagger123"),
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
