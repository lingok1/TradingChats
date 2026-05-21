package repository

import (
	"context"
	"time"
	"trading-chats-backend/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FuturesRecommendationRepository struct {
	collection *mongo.Collection
}

func NewFuturesRecommendationRepository(db *mongo.Database) *FuturesRecommendationRepository {
	return &FuturesRecommendationRepository{collection: db.Collection("futures_recommendations")}
}

func (r *FuturesRecommendationRepository) Save(ctx context.Context, rec *models.FuturesRecommendation) error {
	rec.ID = primitive.NewObjectID().Hex()
	rec.CreatedAt = time.Now()
	_, err := r.collection.InsertOne(ctx, rec)
	return err
}

func (r *FuturesRecommendationRepository) GetLatest(ctx context.Context) (*models.FuturesRecommendation, error) {
	opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})
	var rec models.FuturesRecommendation
	if err := r.collection.FindOne(ctx, bson.M{}, opts).Decode(&rec); err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *FuturesRecommendationRepository) GetLatestByTab(ctx context.Context, tabTag string) (*models.FuturesRecommendation, error) {
	opts := options.FindOne().SetSort(bson.D{{Key: "created_at", Value: -1}})
	filter := bson.M{"tab_tag": tabTag}
	var rec models.FuturesRecommendation
	if err := r.collection.FindOne(ctx, filter, opts).Decode(&rec); err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *FuturesRecommendationRepository) GetList(ctx context.Context, limit int64) ([]models.FuturesRecommendation, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(limit)
	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var recs []models.FuturesRecommendation
	if err := cursor.All(ctx, &recs); err != nil {
		return nil, err
	}
	if recs == nil {
		recs = []models.FuturesRecommendation{}
	}
	return recs, nil
}
