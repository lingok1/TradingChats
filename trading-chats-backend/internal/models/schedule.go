package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ScheduleConfig 定时任务配置模型
type ScheduleConfig struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name       string             `bson:"name" json:"name"`
	CronExpr   string             `bson:"cron_expr" json:"cron_expr"`     // cron表达式，如 "*/10 * * * *"
	TemplateID string             `bson:"template_id" json:"template_id"` // 关联的提示词模版ID
	Param1     string             `bson:"param1" json:"param1"`           // 参数1
	Param2     string             `bson:"param2" json:"param2"`           // 参数2
	Status     string             `bson:"status" json:"status"`           // active: 启用, paused: 暂停
	CreatedAt  interface{}        `bson:"created_at" json:"created_at"`
	UpdatedAt  interface{}        `bson:"updated_at" json:"updated_at"`
}

// ScheduleLog 定时任务执行记录模型
type ScheduleLog struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ScheduleConfigID primitive.ObjectID `bson:"schedule_config_id" json:"schedule_config_id"`
	BatchID          string             `bson:"batch_id" json:"batch_id"` // 关联的 AI Response 批次 ID
	Status           string             `bson:"status" json:"status"`     // success, failed
	Error            string             `bson:"error" json:"error"`       // 错误信息
	ExecutedAt       interface{}        `bson:"executed_at" json:"executed_at"`
}
