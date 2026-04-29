package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"trading-chats-backend/internal/models"
	"trading-chats-backend/internal/repository"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PromptTemplateService struct {
	repo                *repository.PromptTemplateRepository
	systemConfigService SystemConfigService
}

func NewPromptTemplateService(repo *repository.PromptTemplateRepository, systemConfigService SystemConfigService) *PromptTemplateService {
	return &PromptTemplateService{
		repo:                repo,
		systemConfigService: systemConfigService,
	}
}

func (s *PromptTemplateService) CreatePromptTemplate(ctx context.Context, template *models.PromptTemplate) error {
	return s.repo.Create(ctx, template)
}

func (s *PromptTemplateService) GetPromptTemplateByID(ctx context.Context, id string) (*models.PromptTemplate, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}
	return s.repo.GetByID(ctx, objectID)
}

func (s *PromptTemplateService) GetAllPromptTemplates(ctx context.Context) ([]models.PromptTemplate, error) {
	return s.repo.GetAll(ctx)
}

func (s *PromptTemplateService) GetPromptTemplatesByTag(ctx context.Context, tag string) ([]models.PromptTemplate, error) {
	return s.repo.GetByTag(ctx, tag)
}

func (s *PromptTemplateService) UpdatePromptTemplate(ctx context.Context, template *models.PromptTemplate) error {
	return s.repo.Update(ctx, template)
}

func (s *PromptTemplateService) DeletePromptTemplate(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}
	return s.repo.Delete(ctx, objectID)
}

func isValidURL(str string) bool {
	if str == "" {
		return false
	}
	u, err := url.Parse(str)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}

func fetchJSONData(ctx context.Context, urlStr string) (map[string]interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from %s: %w", urlStr, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from %s: %d", urlStr, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON data: %w", err)
	}

	return data, nil
}

func formatJSONData(data map[string]interface{}) string {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("{error: %s}", err.Error())
	}
	return string(jsonBytes)
}

func (s *PromptTemplateService) GeneratePrompt(ctx context.Context, templateID string) (string, error) {
	template, err := s.GetPromptTemplateByID(ctx, templateID)
	if err != nil {
		return "", fmt.Errorf("failed to get prompt template: %w", err)
	}

	sysConfig, err := s.systemConfigService.GetConfig(ctx)
	if err != nil {
		fmt.Printf("Warning: failed to get system config: %v\n", err)
	}

	beijingTime := time.Now().Add(8 * time.Hour)
	openTime := beijingTime.Add(5 * time.Minute)
	prompt := template.Content

	if sysConfig != nil {
		runtimeParameters := cloneParameters(sysConfig.Parameters)
		if runtimeParameters == nil {
			runtimeParameters = map[string]string{}
		}

		for key, val := range runtimeParameters {
			placeholder := fmt.Sprintf("{{.%s}}", key)
			if !strings.Contains(prompt, placeholder) {
				continue
			}
			if isValidURL(val) {
				data, fetchErr := fetchJSONData(ctx, val)
				if fetchErr != nil {
					prompt = strings.ReplaceAll(prompt, placeholder, fmt.Sprintf("[参数%s数据获取失败: %s]", key, fetchErr.Error()))
				} else {
					prompt = strings.ReplaceAll(prompt, placeholder, formatJSONData(data))
				}
			} else {
				prompt = strings.ReplaceAll(prompt, placeholder, val)
			}
		}
	}

	prompt += fmt.Sprintf("\n- 当前北京时间: %s"+"，开仓时间: %s", beijingTime.Format("2006-01-02 15:04:05"),openTime.Format("2006-01-02 15:04:05"))
	

	return prompt, nil
}
