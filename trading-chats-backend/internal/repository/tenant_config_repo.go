package repository

import (
	"context"
	"time"
	"trading-chats-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TenantConfigRepository struct {
	collection *mongo.Collection
}

func NewTenantConfigRepository(db *mongo.Database) *TenantConfigRepository {
	return &TenantConfigRepository{collection: db.Collection("tenant_configs")}
}

func (r *TenantConfigRepository) GetByTenantID(ctx context.Context, tenantID string) (*models.TenantConfig, error) {
	var cfg models.TenantConfig
	err := r.collection.FindOne(ctx, bson.M{"_id": tenantID}).Decode(&cfg)
	if err == mongo.ErrNoDocuments {
		return &models.TenantConfig{
			ID:         tenantID,
			Parameters: map[string]string{},
			MenuConfig: models.TenantMenuConfig{
				VisibleTabs:     models.DefaultVisibleTabs,
				VisibleSettings: models.DefaultVisibleSettings,
			},
		}, nil
	}
	if err != nil {
		return nil, err
	}
	if cfg.Parameters == nil {
		cfg.Parameters = map[string]string{}
	}
	if len(cfg.MenuConfig.VisibleTabs) == 0 {
		cfg.MenuConfig.VisibleTabs = models.DefaultVisibleTabs
	}
	if len(cfg.MenuConfig.VisibleSettings) == 0 {
		cfg.MenuConfig.VisibleSettings = models.DefaultVisibleSettings
	}
	return &cfg, nil
}

func (r *TenantConfigRepository) Save(ctx context.Context, cfg *models.TenantConfig) error {
	cfg.UpdatedAt = time.Now()
	opts := options.Replace().SetUpsert(true)
	_, err := r.collection.ReplaceOne(ctx, bson.M{"_id": cfg.ID}, cfg, opts)
	return err
}

func (r *TenantConfigRepository) GetAll(ctx context.Context) ([]models.TenantConfig, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var configs []models.TenantConfig
	if err := cursor.All(ctx, &configs); err != nil {
		return nil, err
	}
	return configs, nil
}
