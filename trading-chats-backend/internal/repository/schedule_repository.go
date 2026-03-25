package repository

import (
	"context"
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

// CreateConfig 创建定时任务配置
func (r *ScheduleRepository) CreateConfig(ctx context.Context, config *models.ScheduleConfig) error {
	config.CreatedAt = utils.NowString()
	config.UpdatedAt = utils.NowString()
	if config.Status == "" {
		config.Status = "active" // 默认启用
	}

	result, err := r.configCollection.InsertOne(ctx, config)
	if err != nil {
		return err
	}
	config.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetConfigByID 根据ID获取定时任务配置
func (r *ScheduleRepository) GetConfigByID(ctx context.Context, id string) (*models.ScheduleConfig, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var config models.ScheduleConfig
	err = r.configCollection.FindOne(ctx, bson.M{"_id": objID}).Decode(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// GetAllConfigs 获取所有定时任务配置
func (r *ScheduleRepository) GetAllConfigs(ctx context.Context) ([]models.ScheduleConfig, error) {
	cursor, err := r.configCollection.Find(ctx, bson.M{})
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

// GetActiveConfigs 获取所有处于启用状态的配置
func (r *ScheduleRepository) GetActiveConfigs(ctx context.Context) ([]models.ScheduleConfig, error) {
	cursor, err := r.configCollection.Find(ctx, bson.M{"status": "active"})
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

// UpdateConfig 更新定时任务配置
func (r *ScheduleRepository) UpdateConfig(ctx context.Context, id string, updateData map[string]interface{}) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	updateData["updated_at"] = utils.NowString()
	update := bson.M{"$set": updateData}

	_, err = r.configCollection.UpdateOne(ctx, bson.M{"_id": objID}, update)
	return err
}

// DeleteConfig 删除定时任务配置
func (r *ScheduleRepository) DeleteConfig(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	_, err = r.configCollection.DeleteOne(ctx, bson.M{"_id": objID})
	return err
}

// CreateLog 创建定时任务执行记录
func (r *ScheduleRepository) CreateLog(ctx context.Context, log *models.ScheduleLog) error {
	log.ExecutedAt = utils.NowString()
	result, err := r.logCollection.InsertOne(ctx, log)
	if err != nil {
		return err
	}
	log.ID = result.InsertedID.(primitive.ObjectID)
	return nil
}

// GetLogsByConfigID 根据配置ID获取执行记录
func (r *ScheduleRepository) GetLogsByConfigID(ctx context.Context, configID string) ([]models.ScheduleLog, error) {
	objID, err := primitive.ObjectIDFromHex(configID)
	if err != nil {
		return nil, err
	}

	cursor, err := r.logCollection.Find(ctx, bson.M{"schedule_config_id": objID})
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
