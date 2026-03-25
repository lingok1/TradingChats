package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PromptTemplate 提示词模版模型
type PromptTemplate struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Content   string             `bson:"content" json:"content"`
	Tags      []string           `bson:"tags" json:"tags"` // 股票、期货、期权等
	CreatedAt interface{}        `bson:"created_at" json:"created_at"`
	UpdatedAt interface{}        `bson:"updated_at" json:"updated_at"`
}

// ModelAPIConfig 模型与API配置模型
type ModelAPIConfig struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	APIURL    string             `bson:"api_url" json:"api_url"`
	APIKey    string             `bson:"api_key" json:"api_key"`
	Models    []string           `bson:"models" json:"models"`     // 一次可添加多个模型名称
	Provider  string             `bson:"provider" json:"provider"` // anthropic, openai
	CreatedAt interface{}        `bson:"created_at" json:"created_at"`
	UpdatedAt interface{}        `bson:"updated_at" json:"updated_at"`
}

// AIResponse AI模型响应信息模型
type AIResponse struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	BatchID     string             `bson:"batch_id" json:"batch_id"` // 批次ID
	Prompt      string             `bson:"prompt" json:"prompt"`
	Response    string             `bson:"response" json:"response"`
	ModelName   string             `bson:"model_name" json:"model_name"`
	Provider    string             `bson:"provider" json:"provider"`
	Status      string             `bson:"status" json:"status"` // pending, completed, failed
	Error       string             `bson:"error" json:"error"`   // 错误信息
	CreatedAt   interface{}        `bson:"created_at" json:"created_at"`
	UpdatedAt   interface{}        `bson:"updated_at" json:"updated_at"`
	CompletedAt interface{}        `bson:"completed_at" json:"completed_at"`
}

// TaskStatus 任务状态模型（用于Redis）
type TaskStatus struct {
	TaskID    string `json:"task_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Response 通用API响应结构
type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

// SuccessResponse 成功响应
func SuccessResponse(data interface{}) Response {
	return Response{
		Code: 200,
		Msg:  "success",
		Data: data,
	}
}

// ErrorResponse 错误响应
func ErrorResponse(code int, msg string) Response {
	return Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	}
}

// ParseObjectID 将字符串ID转换为ObjectID
func ParseObjectID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}

// GenerateAIRequest 生成AI响应请求体
type GenerateAIRequest struct {
	TemplateID string `json:"template_id" example:"69b4e08fbd4f24a37ff0a0cc"`
	Param1     string `json:"param1" example:"https://api.example.com/data1"`
	Param2     string `json:"param2" example:"https://api.example.com/data2"`
}
