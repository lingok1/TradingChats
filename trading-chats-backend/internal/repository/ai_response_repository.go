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
	db *mongo.Database
}

func NewAIResponseRepository(db *mongo.Database) *AIResponseRepository {
	return &AIResponseRepository{db: db}
}

func (r *AIResponseRepository) collection(tabTag string) *mongo.Collection {
	return r.db.Collection(models.AIResponseCollectionName(tabTag))
}

func (r *AIResponseRepository) Create(ctx context.Context, tabTag string, response *models.AIResponse) error {
	authCtx := models.GetAuthContext(ctx)
	tenantID := models.ResolveTenantID(authCtx, response.TenantID)
	if tenantID == "" {
		return errors.New("tenant_id is required")
	}
	response.TenantID = tenantID
	response.CreatedAt = utils.NowString()
	response.UpdatedAt = utils.NowString()
	result, err := r.collection(tabTag).InsertOne(ctx, response)
	if err != nil {
		return err
	}
	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		response.ID = id
	}
	return nil
}

func (r *AIResponseRepository) GetByID(ctx context.Context, tabTag string, id primitive.ObjectID) (*models.AIResponse, error) {
	filter := bson.M{"_id": id}
	applyTenantFilter(ctx, filter)
	var response models.AIResponse
	if err := r.collection(tabTag).FindOne(ctx, filter).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (r *AIResponseRepository) GetByBatchID(ctx context.Context, tabTag string, batchID string) ([]models.AIResponse, error) {
	filter := bson.M{"batch_id": batchID}
	applyTenantFilter(ctx, filter)
	cursor, err := r.collection(tabTag).Find(ctx, filter)
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

func (r *AIResponseRepository) GetCompletedByBatchID(ctx context.Context, tabTag string, batchID string) ([]models.AIResponse, error) {
	filter := bson.M{
		"batch_id": batchID,
		"status":   "completed",
	}
	applyTenantFilter(ctx, filter)
	opts := options.Find().SetSort(bson.D{
		{Key: "completed_at", Value: -1},
		{Key: "updated_at", Value: -1},
		{Key: "_id", Value: -1},
	})
	cursor, err := r.collection(tabTag).Find(ctx, filter, opts)
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

func (r *AIResponseRepository) GetAll(ctx context.Context, tabTag string) ([]models.AIResponse, error) {
	filter := bson.M{}
	applyTenantFilter(ctx, filter)
	cursor, err := r.collection(tabTag).Find(ctx, filter)
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

func (r *AIResponseRepository) Update(ctx context.Context, tabTag string, response *models.AIResponse) error {
	response.UpdatedAt = utils.NowString()
	filter := bson.M{"_id": response.ID}
	applyTenantFilter(ctx, filter)
	_, err := r.collection(tabTag).UpdateOne(ctx, filter, bson.M{"$set": response})
	return err
}

func (r *AIResponseRepository) GetLatestBatchID(ctx context.Context, tabTag string) (string, error) {
	filter := bson.M{}
	applyTenantFilter(ctx, filter)
	var response models.AIResponse
	opts := options.FindOne().SetSort(bson.D{
		{Key: "created_at", Value: -1},
		{Key: "_id", Value: -1},
	})
	if err := r.collection(tabTag).FindOne(ctx, filter, opts).Decode(&response); err != nil {
		return "", err
	}
	return response.BatchID, nil
}

func (r *AIResponseRepository) GetLatestCompletedBatchID(ctx context.Context, tabTag string) (string, error) {
	filter := bson.M{"status": "completed"}
	applyTenantFilter(ctx, filter)
	var response models.AIResponse
	opts := options.FindOne().SetSort(bson.D{
		{Key: "completed_at", Value: -1},
		{Key: "updated_at", Value: -1},
		{Key: "_id", Value: -1},
	})
	if err := r.collection(tabTag).FindOne(ctx, filter, opts).Decode(&response); err != nil {
		return "", err
	}
	return response.BatchID, nil
}

func (r *AIResponseRepository) Delete(ctx context.Context, tabTag string, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	applyTenantFilter(ctx, filter)
	_, err := r.collection(tabTag).DeleteOne(ctx, filter)
	return err
}
