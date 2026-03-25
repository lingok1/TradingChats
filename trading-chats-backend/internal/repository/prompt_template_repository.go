package repository

import (
	"context"
	"trading-chats-backend/internal/models"
	"trading-chats-backend/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PromptTemplateRepository struct {
	collection *mongo.Collection
}

func NewPromptTemplateRepository(db *mongo.Database) *PromptTemplateRepository {
	return &PromptTemplateRepository{
		collection: db.Collection("prompt_templates"),
	}
}

// Create 创建提示词模版
func (r *PromptTemplateRepository) Create(ctx context.Context, template *models.PromptTemplate) error {
	template.CreatedAt = utils.NowString()
	template.UpdatedAt = utils.NowString()
	result, err := r.collection.InsertOne(ctx, template)
	if err != nil {
		return err
	}
	// 更新ID字段为MongoDB生成的ID
	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		template.ID = id
	}
	return nil
}

// GetByID 根据ID获取提示词模版
func (r *PromptTemplateRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.PromptTemplate, error) {
	var template models.PromptTemplate
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&template)
	if err != nil {
		return nil, err
	}
	return &template, nil
}

// GetAll 获取所有提示词模版
func (r *PromptTemplateRepository) GetAll(ctx context.Context) ([]models.PromptTemplate, error) {
	var templates []models.PromptTemplate
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &templates); err != nil {
		return nil, err
	}

	return templates, nil
}

// GetByTag 根据标签获取提示词模版
func (r *PromptTemplateRepository) GetByTag(ctx context.Context, tag string) ([]models.PromptTemplate, error) {
	var templates []models.PromptTemplate
	cursor, err := r.collection.Find(ctx, bson.M{"tags": tag})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &templates); err != nil {
		return nil, err
	}

	return templates, nil
}

// Update 更新提示词模版
func (r *PromptTemplateRepository) Update(ctx context.Context, template *models.PromptTemplate) error {
	template.UpdatedAt = utils.NowString()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": template.ID},
		bson.M{"$set": template},
	)
	return err
}

// Delete 删除提示词模版
func (r *PromptTemplateRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
