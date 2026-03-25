package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"trading-chats-backend/internal/models"
	"trading-chats-backend/internal/repository"
)

type ModelAPIService struct {
	repo *repository.ModelAPIRepository
}

func NewModelAPIService(repo *repository.ModelAPIRepository) *ModelAPIService {
	return &ModelAPIService{
		repo: repo,
	}
}

// CreateModelAPIConfig 创建模型与API配置
func (s *ModelAPIService) CreateModelAPIConfig(ctx context.Context, config *models.ModelAPIConfig) error {
	return s.repo.Create(ctx, config)
}

// GetModelAPIConfigByID 根据ID获取模型与API配置
func (s *ModelAPIService) GetModelAPIConfigByID(ctx context.Context, id string) (*models.ModelAPIConfig, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}
	return s.repo.GetByID(ctx, objectID)
}

// GetAllModelAPIConfigs 获取所有模型与API配置
func (s *ModelAPIService) GetAllModelAPIConfigs(ctx context.Context) ([]models.ModelAPIConfig, error) {
	return s.repo.GetAll(ctx)
}

// GetModelAPIConfigsByProvider 根据提供商获取模型与API配置
func (s *ModelAPIService) GetModelAPIConfigsByProvider(ctx context.Context, provider string) ([]models.ModelAPIConfig, error) {
	return s.repo.GetByProvider(ctx, provider)
}

// UpdateModelAPIConfig 更新模型与API配置
func (s *ModelAPIService) UpdateModelAPIConfig(ctx context.Context, config *models.ModelAPIConfig) error {
	return s.repo.Update(ctx, config)
}

// DeleteModelAPIConfig 删除模型与API配置
func (s *ModelAPIService) DeleteModelAPIConfig(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}
	return s.repo.Delete(ctx, objectID)
}

// TestModelConnectivity 测试模型的连通性
func (s *ModelAPIService) TestModelConnectivity(ctx context.Context, configID string) error {
	// 获取模型与API配置
	config, err := s.GetModelAPIConfigByID(ctx, configID)
	if err != nil {
		return fmt.Errorf("failed to get model API config: %w", err)
	}

	// 根据提供商构建测试请求
	var req *http.Request
	var errRequest error

	switch config.Provider {
	case "anthropic":
		// Anthropic API: /v1/messages POST请求
		testData := map[string]interface{}{
			"model": config.Models[0],
			"messages": []map[string]string{
				{
					"role":    "user",
					"content": "Hello, test message",
				},
			},
			"max_tokens": 100,
		}
		data, err := json.Marshal(testData)
		if err != nil {
			return fmt.Errorf("failed to marshal test data: %w", err)
		}
		req, errRequest = http.NewRequest("POST", config.APIURL, bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", config.APIKey)

	case "openai":
		// OpenAI API: /v1/chat/completions POST请求
		testData := map[string]interface{}{
			"model": config.Models[0],
			"messages": []map[string]string{
				{
					"role":    "user",
					"content": "Hello, test message",
				},
			},
			"max_tokens": 100,
		}
		data, err := json.Marshal(testData)
		if err != nil {
			return fmt.Errorf("failed to marshal test data: %w", err)
		}
		req, errRequest = http.NewRequest("POST", config.APIURL, bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+config.APIKey)

	default:
		return fmt.Errorf("unsupported provider: %s", config.Provider)
	}

	if errRequest != nil {
		return fmt.Errorf("failed to create request: %w", errRequest)
	}

	// 发送请求测试连通性
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned non-OK status: %d", resp.StatusCode)
	}

	return nil
}
