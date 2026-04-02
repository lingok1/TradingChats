package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PromptTemplate struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TenantID  string             `bson:"tenant_id" json:"tenant_id"`
	Name      string             `bson:"name" json:"name"`
	Content   string             `bson:"content" json:"content"`
	Tags      []string           `bson:"tags" json:"tags"`
	CreatedAt interface{}        `bson:"created_at" json:"created_at"`
	UpdatedAt interface{}        `bson:"updated_at" json:"updated_at"`
}

type ModelAPIConfig struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TenantID  string             `bson:"tenant_id" json:"tenant_id"`
	Name      string             `bson:"name" json:"name"`
	APIURL    string             `bson:"api_url" json:"api_url"`
	APIKey    string             `bson:"api_key" json:"api_key"`
	Models    []string           `bson:"models" json:"models"`
	Provider  string             `bson:"provider" json:"provider"`
	CreatedAt interface{}        `bson:"created_at" json:"created_at"`
	UpdatedAt interface{}        `bson:"updated_at" json:"updated_at"`
}

type AIResponse struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TenantID    string             `bson:"tenant_id" json:"tenant_id"`
	BatchID     string             `bson:"batch_id" json:"batch_id"`
	Prompt      string             `bson:"prompt" json:"prompt"`
	Response    string             `bson:"response" json:"response"`
	ModelName   string             `bson:"model_name" json:"model_name"`
	Provider    string             `bson:"provider" json:"provider"`
	Status      string             `bson:"status" json:"status"`
	Error       string             `bson:"error" json:"error"`
	CreatedAt   interface{}        `bson:"created_at" json:"created_at"`
	UpdatedAt   interface{}        `bson:"updated_at" json:"updated_at"`
	CompletedAt interface{}        `bson:"completed_at" json:"completed_at"`
}

type TaskStatus struct {
	TaskID    string `json:"task_id"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func SuccessResponse(data interface{}) Response {
	return Response{Code: 200, Msg: "success", Data: data}
}

func ErrorResponse(code int, msg string) Response {
	return Response{Code: code, Msg: msg, Data: nil}
}

func ParseObjectID(id string) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id)
}

type GenerateAIRequest struct {
	TemplateID string `json:"template_id" example:"69b4e08fbd4f24a37ff0a0cc"`
}
