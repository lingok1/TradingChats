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

const defaultMaxTokens = 4096

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

func (s *AIResponseService) CreateAIResponse(ctx context.Context, response *models.AIResponse) error {
	return s.repo.Create(ctx, response)
}

func (s *AIResponseService) GetAIResponseByID(ctx context.Context, id string) (*models.AIResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}
	return s.repo.GetByID(ctx, objectID)
}

func (s *AIResponseService) GetAIResponsesByBatchID(ctx context.Context, batchID string) ([]models.AIResponse, error) {
	return s.repo.GetByBatchID(ctx, batchID)
}

func (s *AIResponseService) GetLatestSuccessfulBatch(ctx context.Context) ([]models.AIResponse, error) {
	batchID, err := s.repo.GetLatestSuccessfulBatchID(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to find latest successful batch: %w", err)
	}
	return s.repo.GetByBatchID(ctx, batchID)
}

func (s *AIResponseService) GetAllAIResponses(ctx context.Context) ([]models.AIResponse, error) {
	return s.repo.GetAll(ctx)
}

func (s *AIResponseService) UpdateAIResponse(ctx context.Context, response *models.AIResponse) error {
	return s.repo.Update(ctx, response)
}

func (s *AIResponseService) DeleteAIResponse(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}
	return s.repo.Delete(ctx, objectID)
}

func (s *AIResponseService) GenerateBatchAIResponses(ctx context.Context, templateID string) (string, string, error) {
	batchID := uuid.New().String()

	modelConfigs, err := s.modelAPIService.GetAllModelAPIConfigs(ctx)
	if err != nil {
		return "", "", fmt.Errorf("failed to get model API configs: %w", err)
	}

	prompt, err := s.promptTemplateService.GeneratePrompt(ctx, templateID)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate prompt: %w", err)
	}

	authCtx := models.GetAuthContext(ctx)
	for _, config := range modelConfigs {
		for _, modelName := range config.Models {
			response := &models.AIResponse{
				TenantID:  models.ResolveTenantID(authCtx, config.TenantID),
				BatchID:   batchID,
				Prompt:    prompt,
				ModelName: modelName,
				Provider:  config.Provider,
				Status:    "pending",
			}

			if err := s.repo.Create(ctx, response); err != nil {
				continue
			}

			go func(resp *models.AIResponse, cfg models.ModelAPIConfig, model string) {
				s.callAIModel(context.Background(), resp, cfg, model)
			}(response, config, modelName)
		}
	}

	return batchID, prompt, nil
}

func (s *AIResponseService) callAIModel(ctx context.Context, response *models.AIResponse, config models.ModelAPIConfig, modelName string) {
	var req *http.Request
	var errRequest error

	switch config.Provider {
	case "anthropic":
		reqData := map[string]interface{}{
			"model": modelName,
			"messages": []map[string]string{{
				"role":    "user",
				"content": response.Prompt,
			}},
			"max_tokens": defaultMaxTokens,
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
		reqData := map[string]interface{}{
			"model": modelName,
			"messages": []map[string]string{{
				"role":    "user",
				"content": response.Prompt,
			}},
			"max_tokens": defaultMaxTokens,
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

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		response.Status = "failed"
		response.Error = fmt.Sprintf("failed to read response body: %v", err)
		response.CompletedAt = utils.NowString()
		s.repo.Update(ctx, response)
		return
	}

	bodyStr := string(bodyBytes)
	if bodyStr == "" {
		response.Status = "failed"
		response.Error = "empty response body"
		response.CompletedAt = utils.NowString()
		s.repo.Update(ctx, response)
		return
	}

	if resp.StatusCode != http.StatusOK {
		response.Status = "failed"
		response.Error = fmt.Sprintf("API returned non-OK status: %d", resp.StatusCode)
		response.CompletedAt = utils.NowString()
		s.repo.Update(ctx, response)
		return
	}

	var content string
	var processErr error
	if strings.HasPrefix(strings.TrimSpace(bodyStr), "data:") {
		content, processErr = s.processSSEResponse(bodyStr)
	} else if strings.Contains(strings.ToLower(bodyStr), "<!doctype html") || strings.Contains(strings.ToLower(bodyStr), "<html") {
		response.Status = "failed"
		response.Error = "API returned HTML response instead of JSON. This usually indicates an invalid API URL or authentication error."
		response.CompletedAt = utils.NowString()
		s.repo.Update(ctx, response)
		return
	} else {
		content, processErr = s.processJSONResponse(bodyBytes, config.Provider)
	}

	if processErr != nil {
		response.Status = "failed"
		response.Error = processErr.Error()
		response.CompletedAt = utils.NowString()
		s.repo.Update(ctx, response)
		return
	}

	response.Response = content
	response.Status = "completed"
	response.CompletedAt = utils.NowString()
	s.repo.Update(ctx, response)
}

func (s *AIResponseService) processSSEResponse(body string) (string, error) {
	lines := strings.Split(body, "\n")
	var contentBuilder strings.Builder
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
		if data == "[DONE]" || data == "" {
			continue
		}
		var payload map[string]interface{}
		if err := json.Unmarshal([]byte(data), &payload); err != nil {
			continue
		}
		if choices, ok := payload["choices"].([]interface{}); ok && len(choices) > 0 {
			if choice, ok := choices[0].(map[string]interface{}); ok {
				if delta, ok := choice["delta"].(map[string]interface{}); ok {
					if text, ok := delta["content"].(string); ok {
						contentBuilder.WriteString(text)
					}
				}
			}
		}
	}
	content := strings.TrimSpace(contentBuilder.String())
	if content == "" {
		return "", fmt.Errorf("empty SSE content")
	}
	return content, nil
}

func (s *AIResponseService) processJSONResponse(body []byte, provider string) (string, error) {
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", fmt.Errorf("failed to parse JSON response: %w", err)
	}
	switch provider {
	case "openai":
		choices, ok := payload["choices"].([]interface{})
		if !ok || len(choices) == 0 {
			return "", fmt.Errorf("invalid OpenAI response format")
		}
		choice, ok := choices[0].(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("invalid OpenAI choice format")
		}
		message, ok := choice["message"].(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("invalid OpenAI message format")
		}
		content, ok := message["content"].(string)
		if !ok {
			return "", fmt.Errorf("invalid OpenAI content format")
		}
		return content, nil
	case "anthropic":
		contentItems, ok := payload["content"].([]interface{})
		if !ok || len(contentItems) == 0 {
			return "", fmt.Errorf("invalid Anthropic response format")
		}
		item, ok := contentItems[0].(map[string]interface{})
		if !ok {
			return "", fmt.Errorf("invalid Anthropic content item format")
		}
		content, ok := item["text"].(string)
		if !ok {
			return "", fmt.Errorf("invalid Anthropic text format")
		}
		return content, nil
	default:
		return "", fmt.Errorf("unsupported provider for JSON response: %s", provider)
	}
}
