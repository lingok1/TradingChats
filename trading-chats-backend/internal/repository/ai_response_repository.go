package repository

import (
	"context"
	"errors"
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
	return &AIResponseRepository{collection: db.Collection("ai_responses")}
}

func (r *AIResponseRepository) Create(ctx context.Context, response *models.AIResponse) error {
	authCtx := models.GetAuthContext(ctx)
	tenantID := models.ResolveTenantID(authCtx, response.TenantID)
	if tenantID == "" {
		return errors.New("tenant_id is required")
	}
	response.TenantID = tenantID
	response.CreatedAt = utils.NowString()
	response.UpdatedAt = utils.NowString()
	result, err := r.collection.InsertOne(ctx, response)
	if err != nil {
		return err
	}
	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		response.ID = id
	}
	return nil
}

func (r *AIResponseRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.AIResponse, error) {
	filter := bson.M{"_id": id}
	applyTenantFilter(ctx, filter)
	var response models.AIResponse
	if err := r.collection.FindOne(ctx, filter).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (r *AIResponseRepository) GetByBatchID(ctx context.Context, batchID string) ([]models.AIResponse, error) {
	filter := bson.M{"batch_id": batchID}
	applyTenantFilter(ctx, filter)
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var responses []models.AIResponse
	if err := cursor.All(ctx, &responses); err != nil {
		return nil, err
	}
	return responses, nil
}

func (r *AIResponseRepository) GetAll(ctx context.Context) ([]models.AIResponse, error) {
	filter := bson.M{}
	applyTenantFilter(ctx, filter)
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var responses []models.AIResponse
	if err := cursor.All(ctx, &responses); err != nil {
		return nil, err
	}
	return responses, nil
}

func (r *AIResponseRepository) Update(ctx context.Context, response *models.AIResponse) error {
	response.UpdatedAt = utils.NowString()
	filter := bson.M{"_id": response.ID}
	applyTenantFilter(ctx, filter)
	_, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": response})
	return err
}

func (r *AIResponseRepository) GetLatestSuccessfulBatchID(ctx context.Context) (string, error) {
	filter := bson.M{"status": "completed"}
	applyTenantFilter(ctx, filter)
	var response models.AIResponse
	opts := options.FindOne().SetSort(bson.M{"created_at": -1})
	if err := r.collection.FindOne(ctx, filter, opts).Decode(&response); err != nil {
		return "", err
	}
	return response.BatchID, nil
}

func (r *AIResponseRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	applyTenantFilter(ctx, filter)
	_, err := r.collection.DeleteOne(ctx, filter)
	return err
}
