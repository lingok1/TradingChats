package db

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"trading-chats-backend/internal/config"
)

var (
	MongoClient *mongo.Client
	RedisClient *redis.Client
)

// Connect establishes connections to MongoDB and Redis
func Connect(cfg *config.Config) error {
	// Connect to MongoDB
	ctx := context.Background()
	mongoOpts := options.Client().ApplyURI(cfg.MongoDB.URI)
	mongoClient, err := mongo.Connect(ctx, mongoOpts)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Ping MongoDB to ensure connection
	if err := mongoClient.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	MongoClient = mongoClient
	log.Println("Connected to MongoDB successfully")

	// Connect to Redis
	opts, err := redis.ParseURL(cfg.Redis.URI)
	if err != nil {
		return fmt.Errorf("failed to parse Redis URI: %w", err)
	}

	redisClient := redis.NewClient(opts)

	// Ping Redis to ensure connection
	if err := redisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	RedisClient = redisClient
	log.Println("Connected to Redis successfully")

	return nil
}

// Disconnect closes connections to MongoDB and Redis
func Disconnect() error {
	ctx := context.Background()

	// Disconnect from MongoDB
	if MongoClient != nil {
		if err := MongoClient.Disconnect(ctx); err != nil {
			return fmt.Errorf("failed to disconnect from MongoDB: %w", err)
		}
		log.Println("Disconnected from MongoDB")
	}

	// Disconnect from Redis
	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			return fmt.Errorf("failed to disconnect from Redis: %w", err)
		}
		log.Println("Disconnected from Redis")
	}

	return nil
}
