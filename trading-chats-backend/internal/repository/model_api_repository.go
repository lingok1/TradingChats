package repository

import (
	"context"
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
	return &ModelAPIRepository{
		collection: db.Collection("model_api_configs"),
	}
}

// Create 创建模型与API配置
func (r *ModelAPIRepository) Create(ctx context.Context, config *models.ModelAPIConfig) error {
	config.CreatedAt = utils.NowString()
	config.UpdatedAt = utils.NowString()
	result, err := r.collection.InsertOne(ctx, config)
	if err != nil {
		return err
	}
	// 更新ID字段为MongoDB生成的ID
	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		config.ID = id
	}
	return nil
}

// GetByID 根据ID获取模型与API配置
func (r *ModelAPIRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.ModelAPIConfig, error) {
	var config models.ModelAPIConfig
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// GetAll 获取所有模型与API配置
func (r *ModelAPIRepository) GetAll(ctx context.Context) ([]models.ModelAPIConfig, error) {
	var configs []models.ModelAPIConfig
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &configs); err != nil {
		return nil, err
	}

	return configs, nil
}

// GetByProvider 根据提供商获取模型与API配置
func (r *ModelAPIRepository) GetByProvider(ctx context.Context, provider string) ([]models.ModelAPIConfig, error) {
	var configs []models.ModelAPIConfig
	cursor, err := r.collection.Find(ctx, bson.M{"provider": provider})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &configs); err != nil {
		return nil, err
	}

	return configs, nil
}

// Update 更新模型与API配置
func (r *ModelAPIRepository) Update(ctx context.Context, config *models.ModelAPIConfig) error {
	config.UpdatedAt = utils.NowString()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": config.ID},
		bson.M{"$set": config},
	)
	return err
}

// Delete 删除模型与API配置
func (r *ModelAPIRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
