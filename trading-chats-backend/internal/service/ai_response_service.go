package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
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
	"go.mongodb.org/mongo-driver/mongo"
)

const defaultMaxTokens = 4096

type AIResponseService struct {
	repo                  *repository.AIResponseRepository
	modelAPIService       *ModelAPIService
	promptTemplateService *PromptTemplateService
	systemConfigService   SystemConfigService
	eventService          *AIResponseEventService
}

func NewAIResponseService(
	repo *repository.AIResponseRepository,
	modelAPIService *ModelAPIService,
	promptTemplateService *PromptTemplateService,
	systemConfigService SystemConfigService,
	eventService *AIResponseEventService,
) *AIResponseService {
	return &AIResponseService{
		repo:                  repo,
		modelAPIService:       modelAPIService,
		promptTemplateService: promptTemplateService,
		systemConfigService:   systemConfigService,
		eventService:          eventService,
	}
}

func normalizeAIResponseTab(tabTag string) string {
	return models.NormalizeTabTag(tabTag)
}

func cloneAuthContext(authCtx *models.AuthContext) *models.AuthContext {
	if authCtx == nil {
		return nil
	}
	return &models.AuthContext{
		UserID:    authCtx.UserID,
		TenantID:  authCtx.TenantID,
		Role:      authCtx.Role,
		Username:  authCtx.Username,
		SessionID: authCtx.SessionID,
	}
}

func (s *AIResponseService) CreateAIResponse(ctx context.Context, tabTag string, response *models.AIResponse) error {
	return s.repo.Create(ctx, normalizeAIResponseTab(tabTag), response)
}

func (s *AIResponseService) GetAIResponseByID(ctx context.Context, tabTag string, id string) (*models.AIResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}
	return s.repo.GetByID(ctx, normalizeAIResponseTab(tabTag), objectID)
}

func (s *AIResponseService) GetAIResponsesByBatchID(ctx context.Context, tabTag string, batchID string) ([]models.AIResponse, error) {
	return s.repo.GetByBatchID(ctx, normalizeAIResponseTab(tabTag), batchID)
}

func (s *AIResponseService) GetLatestBatch(ctx context.Context, tabTag string) ([]models.AIResponse, error) {
	normalizedTabTag := normalizeAIResponseTab(tabTag)
	batchID, err := s.repo.GetLatestCompletedBatchID(ctx, normalizedTabTag)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []models.AIResponse{}, nil
		}
		return nil, fmt.Errorf("failed to find latest batch: %w", err)
	}
	responses, err := s.repo.GetCompletedByBatchID(ctx, normalizedTabTag, batchID)
	if err != nil {
		return nil, fmt.Errorf("failed to load latest completed batch: %w", err)
	}
	if responses == nil {
		return []models.AIResponse{}, nil
	}
	return responses, nil
}

func (s *AIResponseService) GetAllAIResponses(ctx context.Context, tabTag string) ([]models.AIResponse, error) {
	return s.repo.GetAll(ctx, normalizeAIResponseTab(tabTag))
}

func (s *AIResponseService) UpdateAIResponse(ctx context.Context, tabTag string, response *models.AIResponse) error {
	return s.repo.Update(ctx, normalizeAIResponseTab(tabTag), response)
}

func (s *AIResponseService) DeleteAIResponse(ctx context.Context, tabTag string, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}
	return s.repo.Delete(ctx, normalizeAIResponseTab(tabTag), objectID)
}

func (s *AIResponseService) GenerateBatchAIResponses(ctx context.Context, templateID string, tabTag string) (string, error) {
	normalizedTabTag := normalizeAIResponseTab(tabTag)
	batchID := uuid.New().String()

	modelConfigs, err := s.modelAPIService.GetEnabledModelAPIConfigsByTab(ctx, normalizedTabTag)
	if err != nil {
		return "", fmt.Errorf("failed to get model API configs: %w", err)
	}
	if len(modelConfigs) == 0 {
		return "", fmt.Errorf("no enabled model API configs found for tab %s", normalizedTabTag)
	}

	prompt, err := s.promptTemplateService.GeneratePrompt(ctx, templateID)
	if err != nil {
		return "", fmt.Errorf("failed to generate prompt: %w", err)
	}

	authCtx := cloneAuthContext(models.GetAuthContext(ctx))

	for _, config := range modelConfigs {
		for _, modelName := range config.Models {
			response := &models.AIResponse{
				TenantID:     models.ResolveTenantID(authCtx, config.TenantID),
				BatchID:      batchID,
				ModelAPIID:   config.ID,
				ModelAPIName: config.Name,
				ModelName:    modelName,
				Provider:     config.Provider,
				Status:       "pending",
			}

			if err := s.repo.Create(ctx, normalizedTabTag, response); err != nil {
				continue
			}

			go func(resp *models.AIResponse, cfg models.ModelAPIConfig, model string, p string, auth *models.AuthContext) {
				bg := context.Background()
				if auth != nil {
					bg = models.WithAuthContext(bg, auth)
				}
				s.callAIModel(bg, normalizedTabTag, resp, cfg, model, p)
			}(response, config, modelName, prompt, cloneAuthContext(authCtx))
		}
	}

	return batchID, nil
}

func (s *AIResponseService) markResponseFailed(ctx context.Context, tabTag string, response *models.AIResponse, err error) {
	response.Status = "failed"
	response.Error = err.Error()
	response.CompletedAt = utils.NowString()
	_ = s.repo.Update(ctx, tabTag, response)
}

func (s *AIResponseService) publishResponseEvent(tabTag string, response *models.AIResponse) {
	if s.eventService == nil || response == nil {
		return
	}
	if response.Status != "completed" {
		return
	}

	event := models.AIResponseEvent{
		Type:         "ai_response_updated",
		TabTag:       normalizeAIResponseTab(tabTag),
		BatchID:      response.BatchID,
		Status:       response.Status,
		ModelName:    response.ModelName,
		ModelAPIName: response.ModelAPIName,
		TenantID:     response.TenantID,
	}
	if !response.ID.IsZero() {
		event.ResponseID = response.ID.Hex()
	}
	if !response.ModelAPIID.IsZero() {
		event.ModelAPIID = response.ModelAPIID.Hex()
	}

	s.eventService.Publish(event)
}

func (s *AIResponseService) callAIModel(ctx context.Context, tabTag string, response *models.AIResponse, config models.ModelAPIConfig, modelName string, prompt string) {
	var req *http.Request
	var errRequest error

	switch config.Provider {
	case "anthropic":
		reqData := map[string]interface{}{
			"model": modelName,
			"messages": []map[string]string{{
				"role":    "user",
				"content": prompt,
			}},
			"max_tokens": defaultMaxTokens,
		}
		data, err := json.Marshal(reqData)
		if err != nil {
			s.markResponseFailed(ctx, tabTag, response, fmt.Errorf("failed to marshal request data: %v", err))
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
				"content": prompt,
			}},
			"max_tokens": defaultMaxTokens,
		}
		data, err := json.Marshal(reqData)
		if err != nil {
			s.markResponseFailed(ctx, tabTag, response, fmt.Errorf("failed to marshal request data: %v", err))
			return
		}
		req, errRequest = http.NewRequest("POST", config.APIURL, bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+config.APIKey)

	default:
		s.markResponseFailed(ctx, tabTag, response, fmt.Errorf("unsupported provider: %s", config.Provider))
		return
	}

	if errRequest != nil {
		s.markResponseFailed(ctx, tabTag, response, fmt.Errorf("failed to create request: %v", errRequest))
		return
	}

	client := &http.Client{Timeout: 4 * 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		s.markResponseFailed(ctx, tabTag, response, fmt.Errorf("failed to send request: %v", err))
		return
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		s.markResponseFailed(ctx, tabTag, response, fmt.Errorf("failed to read response body: %v", err))
		return
	}

	bodyStr := string(bodyBytes)
	if bodyStr == "" {
		s.markResponseFailed(ctx, tabTag, response, fmt.Errorf("empty response body"))
		return
	}

	if resp.StatusCode != http.StatusOK {
		s.markResponseFailed(ctx, tabTag, response, fmt.Errorf("API returned non-OK status: %d", resp.StatusCode))
		return
	}

	var content string
	var processErr error
	if strings.HasPrefix(strings.TrimSpace(bodyStr), "data:") {
		content, processErr = s.processSSEResponse(bodyStr)
	} else if strings.Contains(strings.ToLower(bodyStr), "<!doctype html") || strings.Contains(strings.ToLower(bodyStr), "<html") {
		s.markResponseFailed(ctx, tabTag, response, fmt.Errorf("API returned HTML response instead of JSON. This usually indicates an invalid API URL or authentication error."))
		return
	} else {
		content, processErr = s.processJSONResponse(bodyBytes, config.Provider)
	}

	if processErr != nil {
		s.markResponseFailed(ctx, tabTag, response, processErr)
		return
	}

	if index := strings.Index(content, "| 序号 |"); index > 0 {
		content = content[index:]
	}

	response.Response = content
	response.Status = "completed"
	response.CompletedAt = utils.NowString()
	if err := s.repo.Update(ctx, tabTag, response); err != nil {
		return
	}
	s.publishResponseEvent(tabTag, response)
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
