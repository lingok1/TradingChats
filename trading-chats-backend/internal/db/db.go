package db

import (
	"context"
	"fmt"
	"time"
	"trading-chats-backend/internal/config"

	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	MongoClient *mongo.Client
	MongoDB     *mongo.Database
	RedisClient *redis.Client
)

func Connect(cfg *config.Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoOptions := options.Client().ApplyURI(cfg.MongoDB.URI)
	client, err := mongo.Connect(ctx, mongoOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("failed to ping MongoDB: %w", err)
	}

	MongoClient = client
	MongoDB = client.Database(cfg.MongoDB.Database)

	redisOptions, err := redis.ParseURL(cfg.Redis.URI)
	if err != nil {
		return fmt.Errorf("failed to parse Redis URI: %w", err)
	}
	redisClient := redis.NewClient(redisOptions)

	if err := redisClient.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("failed to ping Redis: %w", err)
	}

	RedisClient = redisClient
	return nil
}

func Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if RedisClient != nil {
		if err := RedisClient.Close(); err != nil {
			return fmt.Errorf("failed to close Redis connection: %w", err)
		}
	}

	if MongoClient != nil {
		if err := MongoClient.Disconnect(ctx); err != nil {
			return fmt.Errorf("failed to disconnect MongoDB: %w", err)
		}
	}

	return nil
}
