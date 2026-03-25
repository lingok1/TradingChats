package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"trading-chats-backend/internal/models"
	"trading-chats-backend/internal/repository"
	"trading-chats-backend/pkg/utils"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AIResponseService struct {
	repo                  *repository.AIResponseRepository
	modelAPIService       *ModelAPIService
	promptTemplateService *PromptTemplateService
	systemConfigService   SystemConfigService
}

func NewAIResponseService(
	repo *repository.AIResponseRepository,
	modelAPIService *ModelAPIService,
	promptTemplateService *PromptTemplateService,
	systemConfigService SystemConfigService,
) *AIResponseService {
	return &AIResponseService{
		repo:                  repo,
		modelAPIService:       modelAPIService,
		promptTemplateService: promptTemplateService,
		systemConfigService:   systemConfigService,
	}
}

// CreateAIResponse 创建AI响应信息
func (s *AIResponseService) CreateAIResponse(ctx context.Context, response *models.AIResponse) error {
	return s.repo.Create(ctx, response)
}

// GetAIResponseByID 根据ID获取AI响应信息
func (s *AIResponseService) GetAIResponseByID(ctx context.Context, id string) (*models.AIResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}
	return s.repo.GetByID(ctx, objectID)
}

// GetAIResponsesByBatchID 根据批次ID获取AI响应信息
func (s *AIResponseService) GetAIResponsesByBatchID(ctx context.Context, batchID string) ([]models.AIResponse, error) {
	return s.repo.GetByBatchID(ctx, batchID)
}

// GetLatestSuccessfulBatch 获取最近一次成功的完整批次数据
func (s *AIResponseService) GetLatestSuccessfulBatch(ctx context.Context) ([]models.AIResponse, error) {
	batchID, err := s.repo.GetLatestSuccessfulBatchID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find latest successful batch: %w", err)
	}
	return s.repo.GetByBatchID(ctx, batchID)
}

// GetAllAIResponses 获取所有AI响应信息
func (s *AIResponseService) GetAllAIResponses(ctx context.Context) ([]models.AIResponse, error) {
	return s.repo.GetAll(ctx)
}

// UpdateAIResponse 更新AI响应信息
func (s *AIResponseService) UpdateAIResponse(ctx context.Context, response *models.AIResponse) error {
	return s.repo.Update(ctx, response)
}

// DeleteAIResponse 删除AI响应信息
func (s *AIResponseService) DeleteAIResponse(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}
	return s.repo.Delete(ctx, objectID)
}

// GenerateBatchAIResponses 生成批次AI响应
func (s *AIResponseService) GenerateBatchAIResponses(ctx context.Context, templateID string, param1, param2 string) (string, error) {
	// 生成批次ID
	batchID := uuid.New().String()

	// 获取所有模型与API配置
	modelConfigs, err := s.modelAPIService.GetAllModelAPIConfigs(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get model API configs: %w", err)
	}

	// 动态生成提示词
	prompt, err := s.promptTemplateService.GeneratePrompt(ctx, templateID, param1, param2)
	if err != nil {
		return "", fmt.Errorf("failed to generate prompt: %w", err)
	}

	// 对每个模型配置生成AI响应
	for _, config := range modelConfigs {
		for _, modelName := range config.Models {
			// 创建AI响应记录
			response := &models.AIResponse{
				BatchID:   batchID,
				Prompt:    prompt,
				ModelName: modelName,
				Provider:  config.Provider,
				Status:    "pending",
			}

			// 保存到数据库
			if err := s.repo.Create(ctx, response); err != nil {
				continue
			}

			// 异步调用AI模型
			go func(resp *models.AIResponse, cfg models.ModelAPIConfig, model string) {
				s.callAIModel(context.Background(), resp, cfg, model)
			}(response, config, modelName)
		}
	}

	return batchID, nil
}

// callAIModel 调用AI模型
func (s *AIResponseService) callAIModel(ctx context.Context, response *models.AIResponse, config models.ModelAPIConfig, modelName string) {
	// 根据提供商构建请求
	var req *http.Request
	var errRequest error

	switch config.Provider {
	case "anthropic":
		// Anthropic API: /v1/messages POST请求
		reqData := map[string]interface{}{
			"model": modelName,
			"messages": []map[string]string{
				{
					"role":    "user",
					"content": response.Prompt,
				},
			},
			"max_tokens": 1000000000,
		}
		data, err := json.Marshal(reqData)
		if err != nil {
			response.Status = "failed"
			response.Error = fmt.Sprintf("failed to marshal request data: %v", err)
			response.CompletedAt = utils.NowString()
			s.repo.Update(ctx, response)
			return
		}
		req, errRequest = http.NewRequest("POST", config.APIURL, bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", config.APIKey)

	case "openai":
		// OpenAI API: /v1/chat/completions POST请求
		reqData := map[string]interface{}{
			"model": modelName,
			"messages": []map[string]string{
				{
					"role":    "user",
					"content": response.Prompt,
				},
			},
			"max_tokens": 1000000000,
		}
		data, err := json.Marshal(reqData)
		if err != nil {
			response.Status = "failed"
			response.Error = fmt.Sprintf("failed to marshal request data: %v", err)
			response.CompletedAt = utils.NowString()
			s.repo.Update(ctx, response)
			return
		}
		req, errRequest = http.NewRequest("POST", config.APIURL, bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+config.APIKey)

	default:
		response.Status = "failed"
		response.Error = fmt.Sprintf("unsupported provider: %s", config.Provider)
		response.CompletedAt = utils.NowString()
		s.repo.Update(ctx, response)
		return
	}

	if errRequest != nil {
		response.Status = "failed"
		response.Error = fmt.Sprintf("failed to create request: %v", errRequest)
		response.CompletedAt = utils.NowString()
		s.repo.Update(ctx, response)
		return
	}

	// 发送请求
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		response.Status = "failed"
		response.Error = fmt.Sprintf("failed to send request: %v", err)
		response.CompletedAt = utils.NowString()
		s.repo.Update(ctx, response)
		return
	}
	defer resp.Body.Close()

	// 读取并检查响应内容
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		response.Status = "failed"
		response.Error = fmt.Sprintf("failed to read response body: %v", err)
		response.CompletedAt = utils.NowString()
		s.repo.Update(ctx, response)
		return
	}

	// 检查响应内容是否为空
	bodyStr := string(bodyBytes)
	if bodyStr == "" {
		response.Status = "failed"
		response.Error = "empty response body"
		response.CompletedAt = utils.NowString()
		s.repo.Update(ctx, response)
		return
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		response.Status = "failed"
		response.Error = fmt.Sprintf("API returned non-OK status: %d", resp.StatusCode)
		response.CompletedAt = utils.NowString()
		s.repo.Update(ctx, response)
		return
	}

	// 处理不同格式的响应
	var content string
	var processErr error

	// 检查是否是 SSE 格式（以 "data:" 开头）
	if strings.HasPrefix(strings.TrimSpace(bodyStr), "data:") {
		content, processErr = s.processSSEResponse(bodyStr)
	} else if strings.Contains(strings.ToLower(bodyStr), "<!doctype html") || strings.Contains(strings.ToLower(bodyStr), "<html") {
		// 处理 HTML 响应
		response.Status = "failed"
		response.Error = "API returned HTML response instead of JSON. This usually indicates an invalid API URL or authentication error."
		response.CompletedAt = utils.NowString()
		s.repo.Update(ctx, response)
		return
	} else {
		// 尝试解析普通 JSON 响应
		var responseData map[string]interface{}
		if err := json.Unmarshal(bodyBytes, &responseData); err != nil {
			// 记录实际响应内容的前200个字符，便于排查
			actualContent := bodyStr
			if len(actualContent) > 200 {
				actualContent = actualContent[:200] + "..."
			}
			response.Status = "failed"
			response.Error = fmt.Sprintf("failed to decode response: %v. Actual content: %s", err, actualContent)
			response.CompletedAt = utils.NowString()
			s.repo.Update(ctx, response)
			return
		}

		// 提取响应内容
		content, processErr = s.extractContentFromJSON(responseData, config.Provider)
	}

	if processErr != nil {
		response.Status = "failed"
		response.Error = fmt.Sprintf("failed to process response: %v", processErr)
		response.CompletedAt = utils.NowString()
		s.repo.Update(ctx, response)
		return
	}

	// 检查是否获取到内容
	if content == "" {
		response.Status = "failed"
		response.Error = "no content found in response"
		response.CompletedAt = utils.NowString()
		s.repo.Update(ctx, response)
		return
	}

	// 更新响应信息
	response.Response = content
	response.Status = "completed"
	response.CompletedAt = utils.NowString()
	s.repo.Update(ctx, response)
}

// processSSEResponse 处理 SSE 格式的响应
func (s *AIResponseService) processSSEResponse(sseContent string) (string, error) {
	var content strings.Builder

	// 按行分割 SSE 内容
	lines := strings.Split(sseContent, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "data:") {
			// 提取 data: 后面的内容
			dataPart := strings.TrimPrefix(line, "data:")
			dataPart = strings.TrimSpace(dataPart)

			// 跳过空数据和结束标记
			if dataPart == "" || dataPart == "[DONE]" {
				continue
			}

			// 解析 JSON
			var chunkData map[string]interface{}
			if err := json.Unmarshal([]byte(dataPart), &chunkData); err != nil {
				return "", fmt.Errorf("failed to parse SSE chunk: %v", err)
			}

			// 提取 delta 内容
			if choices, ok := chunkData["choices"].([]interface{}); ok && len(choices) > 0 {
				if choice, ok := choices[0].(map[string]interface{}); ok {
					if delta, ok := choice["delta"].(map[string]interface{}); ok {
						if text, ok := delta["content"].(string); ok {
							content.WriteString(text)
						}
					}
				}
			}
		}
	}

	return content.String(), nil
}

// extractContentFromJSON 从 JSON 响应中提取内容
func (s *AIResponseService) extractContentFromJSON(responseData map[string]interface{}, provider string) (string, error) {
	switch provider {
	case "anthropic":
		if messages, ok := responseData["messages"].([]interface{}); ok && len(messages) > 0 {
			if message, ok := messages[0].(map[string]interface{}); ok {
				if c, ok := message["content"].(string); ok {
					return c, nil
				}
			}
		}

	case "openai":
		if choices, ok := responseData["choices"].([]interface{}); ok && len(choices) > 0 {
			if choice, ok := choices[0].(map[string]interface{}); ok {
				if message, ok := choice["message"].(map[string]interface{}); ok {
					if c, ok := message["content"].(string); ok {
						return c, nil
					}
				}
			}
		}
	}

	return "", fmt.Errorf("no content found in response for provider: %s", provider)
}
