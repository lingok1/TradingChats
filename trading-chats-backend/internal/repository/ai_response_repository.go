package repository

import (
	"context"
	"trading-chats-backend/internal/models"
	"trading-chats-backend/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AIResponseRepository struct {
	collection *mongo.Collection
}

func NewAIResponseRepository(db *mongo.Database) *AIResponseRepository {
	return &AIResponseRepository{
		collection: db.Collection("ai_responses"),
	}
}

// Create 创建AI响应信息
func (r *AIResponseRepository) Create(ctx context.Context, response *models.AIResponse) error {
	response.CreatedAt = utils.NowString()
	response.UpdatedAt = utils.NowString()
	result, err := r.collection.InsertOne(ctx, response)
	if err != nil {
		return err
	}
	// 更新ID字段为MongoDB生成的ID
	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		response.ID = id
	}
	return nil
}

// GetByID 根据ID获取AI响应信息
func (r *AIResponseRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.AIResponse, error) {
	var response models.AIResponse
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&response)
	if err != nil {
		return nil, err
	}
	return &response, nil
}

// GetByBatchID 根据批次ID获取AI响应信息
func (r *AIResponseRepository) GetByBatchID(ctx context.Context, batchID string) ([]models.AIResponse, error) {
	var responses []models.AIResponse
	cursor, err := r.collection.Find(ctx, bson.M{"batch_id": batchID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &responses); err != nil {
		return nil, err
	}

	return responses, nil
}

// GetAll 获取所有AI响应信息
func (r *AIResponseRepository) GetAll(ctx context.Context) ([]models.AIResponse, error) {
	var responses []models.AIResponse
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &responses); err != nil {
		return nil, err
	}

	return responses, nil
}

// Update 更新AI响应信息
func (r *AIResponseRepository) Update(ctx context.Context, response *models.AIResponse) error {
	response.UpdatedAt = utils.NowString()
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"_id": response.ID},
		bson.M{"$set": response},
	)
	return err
}

// GetLatestSuccessfulBatchID 获取最近一个成功的批次ID
func (r *AIResponseRepository) GetLatestSuccessfulBatchID(ctx context.Context) (string, error) {
	var response models.AIResponse
	// 按 CreatedAt 降序排列，找到第一个状态为 completed 的记录
	opts := options.FindOne().SetSort(bson.M{"created_at": -1})
	err := r.collection.FindOne(ctx, bson.M{"status": "completed"}, opts).Decode(&response)
	if err != nil {
		return "", err
	}
	return response.BatchID, nil
}

// Delete 删除AI响应信息
func (r *AIResponseRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}
