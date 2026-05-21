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
)

const recommendationPromptTemplate = `你是一位专业的期货交易分析师。以下是多个AI模型对当前期货市场的分析结果：

%s

请根据以上所有分析结果，综合判断并推荐3个最具交易价值的期货品种。

输出格式必须严格按照以下JSON数组，不要输出任何其他内容：
[
  {
    "symbol": "品种名称和合约代码，例如：聚丙烯(pp2609)",
    "direction": "做多或做空",
    "entry_range": "入场价格区间，例如：9050-9080",
    "take_profit": "止盈价格",
    "stop_loss": "止损价格",
    "reason": "推荐理由，不超过50字"
  }
]`

type FuturesRecommendationService struct {
	repo            *repository.FuturesRecommendationRepository
	aiResponseRepo  *repository.AIResponseRepository
	modelAPIService *ModelAPIService
	eventService    *AIResponseEventService
}

func NewFuturesRecommendationService(
	repo *repository.FuturesRecommendationRepository,
	aiResponseRepo *repository.AIResponseRepository,
	modelAPIService *ModelAPIService,
) *FuturesRecommendationService {
	return &FuturesRecommendationService{
		repo:            repo,
		aiResponseRepo:  aiResponseRepo,
		modelAPIService: modelAPIService,
	}
}

// SetEventService 注入事件服务（用于推荐生成完成后通知前端）
func (s *FuturesRecommendationService) SetEventService(eventService *AIResponseEventService) {
	s.eventService = eventService
}

func (s *FuturesRecommendationService) GetLatest(ctx context.Context) (*models.FuturesRecommendation, error) {
	return s.repo.GetLatest(ctx)
}

func (s *FuturesRecommendationService) GetLatestByTab(ctx context.Context, tabTag string) (*models.FuturesRecommendation, error) {
	if tabTag == "" {
		return s.repo.GetLatest(ctx)
	}
	return s.repo.GetLatestByTab(ctx, tabTag)
}

func (s *FuturesRecommendationService) GetList(ctx context.Context, limit int64) ([]models.FuturesRecommendation, error) {
	return s.repo.GetList(ctx, limit)
}

// Generate 保留向后兼容
func (s *FuturesRecommendationService) Generate(ctx context.Context, modelAPIID, modelName string) error {
	// 从最新分析批次构建默认 prompt
	batchID, err := s.aiResponseRepo.GetLatestCompletedBatchID(ctx, models.TabTagFutures)
	if err != nil {
		return fmt.Errorf("no completed futures analysis found: %w", err)
	}
	responses, err := s.aiResponseRepo.GetCompletedByBatchID(ctx, models.TabTagFutures, batchID)
	if err != nil || len(responses) == 0 {
		return fmt.Errorf("failed to load futures analysis: %w", err)
	}
	var sb strings.Builder
	for i, r := range responses {
		sb.WriteString(fmt.Sprintf("=== 模型%d（%s）分析结果 ===\n%s\n\n", i+1, r.ModelName, r.Response))
	}
	prompt := fmt.Sprintf(recommendationPromptTemplate, sb.String())
	return s.GenerateWithPrompt(ctx, modelAPIID, modelName, models.TabTagFutures, prompt)
}

// GenerateWithPrompt 使用完整 prompt 生成推荐（prompt 中应已包含分析数据）
func (s *FuturesRecommendationService) GenerateWithPrompt(ctx context.Context, modelAPIID, modelName, tabTag, prompt string) error {
	if tabTag == "" {
		tabTag = models.TabTagFutures
	}

	// 关联的 batch_id（用于追溯当前推荐来自哪批分析）
	batchID, _ := s.aiResponseRepo.GetLatestCompletedBatchID(ctx, tabTag)

	// 获取指定模型配置
	modelConfig, err := s.modelAPIService.GetModelAPIConfigByID(ctx, modelAPIID)
	if err != nil {
		return fmt.Errorf("model config not found: %w", err)
	}

	// 调用 AI
	rawResponse, err := s.callModel(ctx, *modelConfig, modelName, prompt)
	if err != nil {
		return fmt.Errorf("AI call failed: %w", err)
	}

	// 解析 JSON
	items, err := parseRecommendationJSON(rawResponse)
	if err != nil {
		return fmt.Errorf("failed to parse recommendation: %w", err)
	}

	// 存库
	rec := &models.FuturesRecommendation{
		BatchID:      batchID,
		TabTag:       tabTag,
		Items:        items,
		RawResponse:  rawResponse,
		ModelName:    modelName,
		ModelAPIName: modelConfig.Name,
	}
	if err := s.repo.Save(ctx, rec); err != nil {
		return err
	}

	// 发布事件，通知前端刷新（复用 AIResponseEventService，事件类型为 recommendation_updated）
	if s.eventService != nil {
		event := models.AIResponseEvent{
			Type:         "recommendation_updated",
			TabTag:       tabTag,
			BatchID:      batchID,
			Status:       "completed",
			ModelName:    modelName,
			ModelAPIName: modelConfig.Name,
			ResponseID:   rec.ID,
		}
		// 同时发布到对应 tab 频道和 home 频道（事件内容不变，仍含原 tab_tag）
		s.eventService.Publish(event)
		s.eventService.PublishToChannel("home", event)
	}

	return nil
}

func (s *FuturesRecommendationService) callModel(ctx context.Context, config models.ModelAPIConfig, modelName, prompt string) (string, error) {
	reqData := map[string]interface{}{
		"model": modelName,
		"messages": []map[string]string{{
			"role":    "user",
			"content": prompt,
		}},
		"max_tokens": 8192,
	}
	data, _ := json.Marshal(reqData)

	req, err := http.NewRequestWithContext(ctx, "POST", config.APIURL, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	switch config.Provider {
	case "anthropic":
		req.Header.Set("x-api-key", config.APIKey)
		req.Header.Set("anthropic-version", "2023-06-01")
	default:
		req.Header.Set("Authorization", "Bearer "+config.APIKey)
	}

	client := &http.Client{Timeout: 4 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API status %d: %s", resp.StatusCode, string(body))
	}

	content, err := extractContent(string(body), config.Provider)
	if err != nil {
		return "", fmt.Errorf("cannot extract content: %w, raw body: %.500s", err, string(body))
	}
	return content, nil
}

func extractContent(body, provider string) (string, error) {
	var result map[string]interface{}
	if err := json.Unmarshal([]byte(body), &result); err != nil {
		return "", err
	}
	switch provider {
	case "anthropic":
		if content, ok := result["content"].([]interface{}); ok && len(content) > 0 {
			if block, ok := content[0].(map[string]interface{}); ok {
				if text, ok := block["text"].(string); ok {
					return text, nil
				}
			}
		}
	default: // openai compatible
		if choices, ok := result["choices"].([]interface{}); ok && len(choices) > 0 {
			if choice, ok := choices[0].(map[string]interface{}); ok {
				if msg, ok := choice["message"].(map[string]interface{}); ok {
					if text, ok := msg["content"].(string); ok && text != "" {
						return text, nil
					}
					// fallback: some models put content in reasoning_content
					if text, ok := msg["reasoning_content"].(string); ok && text != "" {
						return text, nil
					}
				}
			}
		}
	}
	return "", fmt.Errorf("cannot extract content from response")
}

func parseRecommendationJSON(raw string) ([]models.RecommendationItem, error) {
	// 提取 JSON 数组部分
	start := strings.Index(raw, "[")
	end := strings.LastIndex(raw, "]")
	if start < 0 || end <= start {
		return nil, fmt.Errorf("no JSON array found in response")
	}
	jsonStr := raw[start : end+1]

	var items []models.RecommendationItem
	if err := json.Unmarshal([]byte(jsonStr), &items); err != nil {
		return nil, err
	}
	return items, nil
}
