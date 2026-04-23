package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ScheduleConfig struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TenantID   string             `bson:"tenant_id" json:"tenant_id"`
	Name       string             `bson:"name" json:"name"`
	CronExpr   string             `bson:"cron_expr" json:"cron_expr"`
	TemplateID string             `bson:"template_id" json:"template_id"`
	TabTag     string             `bson:"tab_tag" json:"tab_tag"`
	Status     string             `bson:"status" json:"status"`
	CreatedAt  interface{}        `bson:"created_at" json:"created_at"`
	UpdatedAt  interface{}        `bson:"updated_at" json:"updated_at"`
}

type ScheduleLog struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TenantID         string             `bson:"tenant_id" json:"tenant_id"`
	ScheduleConfigID primitive.ObjectID `bson:"schedule_config_id" json:"schedule_config_id"`
	TabTag           string             `bson:"tab_tag" json:"tab_tag"`
	BatchID          string             `bson:"batch_id" json:"batch_id"`
	Prompt           string             `bson:"prompt" json:"prompt"`
	TriggerType      string             `bson:"trigger_type" json:"trigger_type"`
	Status           string             `bson:"status" json:"status"`
	Error            string             `bson:"error" json:"error"`
	ExecutedAt       interface{}        `bson:"executed_at" json:"executed_at"`
}

type UpdateStatusRequest struct {
	ID     string `json:"id" binding:"required"`
	Status string `json:"status" binding:"required,oneof=active paused"`
}
