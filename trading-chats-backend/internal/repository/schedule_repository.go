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

type ScheduleRepository struct {
	configCollection *mongo.Collection
	logCollection    *mongo.Collection
}

func NewScheduleRepository(db *mongo.Database) *ScheduleRepository {
	return &ScheduleRepository{
		configCollection: db.Collection("schedule_configs"),
		logCollection:    db.Collection("schedule_logs"),
	}
}

func (r *ScheduleRepository) CreateConfig(ctx context.Context, config *models.ScheduleConfig) error {
	authCtx := models.GetAuthContext(ctx)
	tenantID := models.ResolveTenantID(authCtx, config.TenantID)
	if tenantID == "" {
		return errors.New("tenant_id is required")
	}
	config.TenantID = tenantID
	config.CreatedAt = utils.NowString()
	config.UpdatedAt = utils.NowString()
	if config.Status == "" {
		config.Status = "active"
	}
	result, err := r.configCollection.InsertOne(ctx, config)
	if err != nil {
		return err
	}
	config.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ScheduleRepository) GetConfigByID(ctx context.Context, id string) (*models.ScheduleConfig, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}
	applyTenantFilter(ctx, filter)
	var config models.ScheduleConfig
	if err := r.configCollection.FindOne(ctx, filter).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *ScheduleRepository) GetAllConfigs(ctx context.Context) ([]models.ScheduleConfig, error) {
	filter := bson.M{}
	applyTenantFilter(ctx, filter)
	cursor, err := r.configCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var configs []models.ScheduleConfig
	if err := cursor.All(ctx, &configs); err != nil {
		return nil, err
	}
	return configs, nil
}

func (r *ScheduleRepository) GetActiveConfigs(ctx context.Context) ([]models.ScheduleConfig, error) {
	filter := bson.M{"status": "active"}
	applyTenantFilter(ctx, filter)
	cursor, err := r.configCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var configs []models.ScheduleConfig
	if err := cursor.All(ctx, &configs); err != nil {
		return nil, err
	}
	return configs, nil
}

func (r *ScheduleRepository) UpdateConfig(ctx context.Context, id string, updateData map[string]interface{}) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	updateData["updated_at"] = utils.NowString()
	filter := bson.M{"_id": objID}
	applyTenantFilter(ctx, filter)
	_, err = r.configCollection.UpdateOne(ctx, filter, bson.M{"$set": updateData})
	return err
}

func (r *ScheduleRepository) DeleteConfig(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	filter := bson.M{"_id": objID}
	applyTenantFilter(ctx, filter)
	_, err = r.configCollection.DeleteOne(ctx, filter)
	return err
}

func (r *ScheduleRepository) CreateLog(ctx context.Context, log *models.ScheduleLog) error {
	authCtx := models.GetAuthContext(ctx)
	if log.TenantID == "" && authCtx != nil {
		log.TenantID = authCtx.TenantID
	}
	log.ExecutedAt = utils.NowString()
	result, err := r.logCollection.InsertOne(ctx, log)
	if err != nil {
		return err
	}
	log.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

func (r *ScheduleRepository) GetLogsByConfigID(ctx context.Context, configID string) ([]models.ScheduleLog, error) {
	objID, err := primitive.ObjectIDFromHex(configID)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"schedule_config_id": objID}
	applyTenantFilter(ctx, filter)
	cursor, err := r.logCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	var logs []models.ScheduleLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, err
	}
	return logs, nil
}
