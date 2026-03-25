package repository

import (
	"context"
	"time"
	"trading-chats-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SystemConfigRepository interface {
	GetConfig(ctx context.Context) (*models.SystemConfig, error)
	SaveConfig(ctx context.Context, config *models.SystemConfig) error
}

type systemConfigRepository struct {
	collection *mongo.Collection
}

func NewSystemConfigRepository(db *mongo.Database) SystemConfigRepository {
	return &systemConfigRepository{
		collection: db.Collection("system_configs"),
	}
}

func (r *systemConfigRepository) GetConfig(ctx context.Context) (*models.SystemConfig, error) {
	var config models.SystemConfig
	err := r.collection.FindOne(ctx, bson.M{"_id": models.GlobalSystemConfigID}).Decode(&config)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return &models.SystemConfig{
				ID:          models.GlobalSystemConfigID,
				SystemTitle: "Trading Chats",
				SystemLogo:  "",
				Parameters:  map[string]string{},
				UpdatedAt:   time.Now(),
			}, nil
		}
		return nil, err
	}

	if config.ID == "" {
		config.ID = models.GlobalSystemConfigID
	}
	if config.Parameters == nil {
		config.Parameters = map[string]string{}
	}

	return &config, nil
}

func (r *systemConfigRepository) SaveConfig(ctx context.Context, config *models.SystemConfig) error {
	if config == nil {
		config = &models.SystemConfig{}
	}

	doc := &models.SystemConfig{
		ID:          config.ID,
		SystemTitle: config.SystemTitle,
		SystemLogo:  config.SystemLogo,
		Parameters:  cloneParameters(config.Parameters),
		UpdatedAt:   config.UpdatedAt,
	}
	if doc.ID == "" {
		doc.ID = models.GlobalSystemConfigID
	}
	if doc.Parameters == nil {
		doc.Parameters = map[string]string{}
	}
	if doc.UpdatedAt.IsZero() {
		doc.UpdatedAt = time.Now()
	}

	opts := options.Replace().SetUpsert(true)
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": models.GlobalSystemConfigID}, doc, opts)
	return err
}

func cloneParameters(source map[string]string) map[string]string {
	if source == nil {
		return map[string]string{}
	}

	cloned := make(map[string]string, len(source))
	for key, value := range source {
		cloned[key] = value
	}

	return cloned
}
