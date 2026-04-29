package repository

import (
	"context"
	"errors"
	"trading-chats-backend/internal/models"
	"trading-chats-backend/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ModelAPIRepository struct {
	collection *mongo.Collection
}

func NewModelAPIRepository(db *mongo.Database) *ModelAPIRepository {
	return &ModelAPIRepository{collection: db.Collection("model_api_configs")}
}

func (r *ModelAPIRepository) Create(ctx context.Context, config *models.ModelAPIConfig) error {
	authCtx := models.GetAuthContext(ctx)
	tenantID := models.ResolveTenantID(authCtx, config.TenantID)
	if tenantID == "" {
		return errors.New("tenant_id is required")
	}
	config.TenantID = tenantID
	config.CreatedAt = utils.NowString()
	config.UpdatedAt = utils.NowString()
	result, err := r.collection.InsertOne(ctx, config)
	if err != nil {
		return err
	}
	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		config.ID = id
	}
	return nil
}

func (r *ModelAPIRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.ModelAPIConfig, error) {
	filter := bson.M{"_id": id}
	applyTenantFilter(ctx, filter)
	var config models.ModelAPIConfig
	if err := r.collection.FindOne(ctx, filter).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *ModelAPIRepository) GetAll(ctx context.Context) ([]models.ModelAPIConfig, error) {
	filter := bson.M{}
	applyTenantFilter(ctx, filter)
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var configs []models.ModelAPIConfig
	if err := cursor.All(ctx, &configs); err != nil {
		return nil, err
	}
	return configs, nil
}

func (r *ModelAPIRepository) GetByProvider(ctx context.Context, provider string) ([]models.ModelAPIConfig, error) {
	filter := bson.M{"provider": provider}
	applyTenantFilter(ctx, filter)
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var configs []models.ModelAPIConfig
	if err := cursor.All(ctx, &configs); err != nil {
		return nil, err
	}
	return configs, nil
}

func (r *ModelAPIRepository) GetEnabledByTabTag(ctx context.Context, tabTag string) ([]models.ModelAPIConfig, error) {
	normalizedTab := models.NormalizeTabTag(tabTag)
	filter := bson.M{
		"tab_settings": bson.M{
			"$elemMatch": bson.M{
				"tab_tag": normalizedTab,
				"enabled": true,
			},
		},
	}
	applyTenantFilter(ctx, filter)
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var configs []models.ModelAPIConfig
	if err := cursor.All(ctx, &configs); err != nil {
		return nil, err
	}
	return configs, nil
}

func (r *ModelAPIRepository) Update(ctx context.Context, config *models.ModelAPIConfig) error {
	config.UpdatedAt = utils.NowString()
	filter := bson.M{"_id": config.ID}
	applyTenantFilter(ctx, filter)
	_, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": config})
	return err
}

func (r *ModelAPIRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	applyTenantFilter(ctx, filter)
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}
