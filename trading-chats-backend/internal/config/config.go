package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	MongoDB  MongoDBConfig
	Redis    RedisConfig
	JWT      JWTConfig
	API      APIConfig
}

type ServerConfig struct {
	Port string
	Env  string
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
	// Load .env file if it exists
	_ = godotenv.Load()

	jwtExpiration, err := time.ParseDuration(getEnv("JWT_EXPIRATION", "24h"))
	if err != nil {
		jwtExpiration = 24 * time.Hour
	}

	apiTimeout, err := time.ParseDuration(getEnv("API_TIMEOUT", "30s"))
	if err != nil {
		apiTimeout = 30 * time.Second
	}

	return &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		MongoDB: MongoDBConfig{
			URI:      getEnv("MONGODB_URI", "mongodb://admin:mongo123@122.51.4.80:27017"),
			Database: getEnv("MONGODB_DATABASE", "trading_chats"),
		},
		Redis: RedisConfig{
			URI: getEnv("REDIS_URI", "redis://:redis123@122.51.4.80:6379/0"),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "your_jwt_secret_key"),
			Expiration: jwtExpiration,
		},
		API: APIConfig{
			Timeout: apiTimeout,
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
