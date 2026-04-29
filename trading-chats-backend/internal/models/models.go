package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	TabTagFutures  = "futures"
	TabTagOptions  = "options"
	TabTagStock    = "stock"
	TabTagNews     = "news"
	TabTagPosition = "position"
)

func NormalizeTabTag(tabTag string) string {
	switch tabTag {
	case TabTagOptions:
		return TabTagOptions
	case TabTagStock:
		return TabTagStock
	case TabTagNews:
		return TabTagNews
	case TabTagPosition:
		return TabTagPosition
	case TabTagFutures:
		fallthrough
	default:
		return TabTagFutures
	}
}

func AIResponseCollectionName(tabTag string) string {
	switch NormalizeTabTag(tabTag) {
	case TabTagOptions:
		return "ai_responses_options"
	case TabTagStock:
		return "ai_responses_stock"
	case TabTagNews:
		return "ai_responses_news"
	case TabTagPosition:
		return "ai_responses_position"
	default:
		return "ai_responses"
	}
}

type PromptTemplate struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TenantID  string             `bson:"tenant_id" json:"tenant_id"`
	Name      string             `bson:"name" json:"name"`
	Content   string             `bson:"content" json:"content"`
	Tags      []string           `bson:"tags" json:"tags"`
	CreatedAt interface{}        `bson:"created_at" json:"created_at"`
	UpdatedAt interface{}        `bson:"updated_at" json:"updated_at"`
}

type ModelAPITabSetting struct {
	TabTag  string `bson:"tab_tag" json:"tab_tag"`
	Enabled bool   `bson:"enabled" json:"enabled"`
}

type ModelAPIConfig struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	TenantID    string               `bson:"tenant_id" json:"tenant_id"`
	Name        string               `bson:"name" json:"name"`
	APIURL      string               `bson:"api_url" json:"api_url"`
	APIKey      string               `bson:"api_key" json:"api_key"`
	Models      []string             `bson:"models" json:"models"`
	Provider    string               `bson:"provider" json:"provider"`
	TabSettings []ModelAPITabSetting `bson:"tab_settings,omitempty" json:"tab_settings,omitempty"`
	CreatedAt   interface{}          `bson:"created_at" json:"created_at"`
	UpdatedAt   interface{}          `bson:"updated_at" json:"updated_at"`
}

type AIResponse struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TenantID     string             `bson:"tenant_id" json:"tenant_id"`
	BatchID      string             `bson:"batch_id" json:"batch_id"`
	Response     string             `bson:"response" json:"response"`
	ModelAPIID   primitive.ObjectID `bson:"model_api_id,omitempty" json:"model_api_id"`
	ModelAPIName string             `bson:"model_api_name" json:"model_api_name"`
	ModelName    string             `bson:"model_name" json:"model_name"`
	Provider     string             `bson:"provider" json:"provider"`
	Status       string             `bson:"status" json:"status"`
	Error        string             `bson:"error" json:"error"`
	CreatedAt    interface{}        `bson:"created_at" json:"created_at"`
	UpdatedAt    interface{}        `bson:"updated_at" json:"updated_at"`
	CompletedAt  interface{}        `bson:"completed_at" json:"completed_at"`
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
	TabTag     string `json:"tab_tag,omitempty" example:"futures"`
}

type AIResponseEvent struct {
	Type         string `json:"type"`
	TabTag       string `json:"tab_tag"`
	BatchID      string `json:"batch_id"`
	Status       string `json:"status"`
	ModelName    string `json:"model_name"`
	ModelAPIID   string `json:"model_api_id,omitempty"`
	ModelAPIName string `json:"model_api_name,omitempty"`
	TenantID     string `json:"tenant_id,omitempty"`
	ResponseID   string `json:"response_id,omitempty"`
}
