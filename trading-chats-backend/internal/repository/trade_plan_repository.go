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

type TradePlanRepository struct {
	collection *mongo.Collection
}

func NewTradePlanRepository(db *mongo.Database) *TradePlanRepository {
	return &TradePlanRepository{collection: db.Collection("trade_plans")}
}

func (r *TradePlanRepository) EnsureIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "tenant_id", Value: 1},
				{Key: "tab_tag", Value: 1},
				{Key: "status", Value: 1},
				{Key: "updated_at", Value: -1},
			},
		},
		{
			Keys: bson.D{
				{Key: "tenant_id", Value: 1},
				{Key: "symbol", Value: 1},
				{Key: "tab_tag", Value: 1},
			},
		},
	}
	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}

func (r *TradePlanRepository) Create(ctx context.Context, plan *models.TradePlan) error {
	authCtx := models.GetAuthContext(ctx)
	tenantID := models.ResolveTenantID(authCtx, plan.TenantID)
	if tenantID == "" {
		return errors.New("tenant_id is required")
	}

	plan.TenantID = tenantID
	plan.CreatedAt = utils.NowString()
	plan.UpdatedAt = utils.NowString()
	result, err := r.collection.InsertOne(ctx, plan)
	if err != nil {
		return err
	}
	if id, ok := result.InsertedID.(primitive.ObjectID); ok {
		plan.ID = id
	}
	return nil
}

func (r *TradePlanRepository) GetByID(ctx context.Context, id string) (*models.TradePlan, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objID}
	applyTenantFilter(ctx, filter)

	var plan models.TradePlan
	if err := r.collection.FindOne(ctx, filter).Decode(&plan); err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *TradePlanRepository) GetAll(ctx context.Context, tabTag string) ([]models.TradePlan, error) {
	filter := bson.M{
		"tab_tag": tabTag,
	}
	applyTenantFilter(ctx, filter)

	opts := options.Find().SetSort(bson.D{{Key: "updated_at", Value: -1}, {Key: "_id", Value: -1}})
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var plans []models.TradePlan
	if err := cursor.All(ctx, &plans); err != nil {
		return nil, err
	}
	return plans, nil
}

func (r *TradePlanRepository) Update(ctx context.Context, plan *models.TradePlan) error {
	plan.UpdatedAt = utils.NowString()
	filter := bson.M{"_id": plan.ID}
	applyTenantFilter(ctx, filter)
	_, err := r.collection.UpdateOne(ctx, filter, bson.M{"$set": plan})
	return err
}

func (r *TradePlanRepository) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	applyTenantFilter(ctx, filter)
	_, err = r.collection.DeleteOne(ctx, filter)
	return err
}
