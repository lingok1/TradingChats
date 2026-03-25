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

// CreatePromptTemplate 创建提示词模版
func (s *PromptTemplateService) CreatePromptTemplate(ctx context.Context, template *models.PromptTemplate) error {
	return s.repo.Create(ctx, template)
}

// GetPromptTemplateByID 根据ID获取提示词模版
func (s *PromptTemplateService) GetPromptTemplateByID(ctx context.Context, id string) (*models.PromptTemplate, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}
	return s.repo.GetByID(ctx, objectID)
}

// GetAllPromptTemplates 获取所有提示词模版
func (s *PromptTemplateService) GetAllPromptTemplates(ctx context.Context) ([]models.PromptTemplate, error) {
	return s.repo.GetAll(ctx)
}

// GetPromptTemplatesByTag 根据标签获取提示词模版
func (s *PromptTemplateService) GetPromptTemplatesByTag(ctx context.Context, tag string) ([]models.PromptTemplate, error) {
	return s.repo.GetByTag(ctx, tag)
}

// UpdatePromptTemplate 更新提示词模版
func (s *PromptTemplateService) UpdatePromptTemplate(ctx context.Context, template *models.PromptTemplate) error {
	return s.repo.Update(ctx, template)
}

// DeletePromptTemplate 删除提示词模版
func (s *PromptTemplateService) DeletePromptTemplate(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}
	return s.repo.Delete(ctx, objectID)
}

// isValidURL 校验是否是有效的http或https URL
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

// fetchJSONData 通过GET请求获取JSON数据
func fetchJSONData(ctx context.Context, urlStr string) (map[string]interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/json")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

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

// formatJSONData 格式化JSON数据为字符串
func formatJSONData(data map[string]interface{}) string {
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Sprintf("{error: %s}", err.Error())
	}
	return string(jsonBytes)
}

// GeneratePrompt 动态生成提示词
func (s *PromptTemplateService) GeneratePrompt(ctx context.Context, templateID string, param1, param2 string) (string, error) {
	// 获取提示词模版
	template, err := s.GetPromptTemplateByID(ctx, templateID)
	if err != nil {
		return "", fmt.Errorf("failed to get prompt template: %w", err)
	}

	// 获取系统配置
	sysConfig, err := s.systemConfigService.GetConfig(ctx)
	if err != nil {
		// 如果获取系统配置失败，仅打印日志，不中断生成
		fmt.Printf("Warning: failed to get system config: %v\n", err)
	}

	// 获取当前北京时间
	beijingTime := time.Now().Add(8 * time.Hour)
	// 开仓时间为当前北京时间5分钟后
	openTime := beijingTime.Add(5 * time.Minute)

	// 生成最终提示词
	prompt := template.Content

	// 1. 处理系统动态参数（隐式替换）
	if sysConfig != nil && sysConfig.Parameters != nil {
		for key, val := range sysConfig.Parameters {
			placeholder := fmt.Sprintf("{{.%s}}", key)
			if strings.Contains(prompt, placeholder) {
				if isValidURL(val) {
					data, err := fetchJSONData(ctx, val)
					if err != nil {
						prompt = strings.ReplaceAll(prompt, placeholder, fmt.Sprintf("[参数%s数据获取失败: %s]", key, err.Error()))
					} else {
						prompt = strings.ReplaceAll(prompt, placeholder, formatJSONData(data))
					}
				} else {
					prompt = strings.ReplaceAll(prompt, placeholder, val)
				}
			}
		}
	}

	// 2. 兼容旧的 param1 和 param2 逻辑（追加到末尾）
	if param1 != "" {
		if isValidURL(param1) {
			data, err := fetchJSONData(ctx, param1)
			if err != nil {
				prompt += fmt.Sprintf("\n\n[参数1数据获取失败: %s]", err.Error())
			} else {
				prompt += fmt.Sprintf("\n\n参数1数据:\n%s", formatJSONData(data))
			}
		} else {
			prompt += fmt.Sprintf("\n\n%s", param1)
		}
	}

	if param2 != "" {
		if isValidURL(param2) {
			data, err := fetchJSONData(ctx, param2)
			if err != nil {
				prompt += fmt.Sprintf("\n\n[参数2数据获取失败: %s]", err.Error())
			} else {
				prompt += fmt.Sprintf("\n\n参数2数据:\n%s", formatJSONData(data))
			}
		} else {
			prompt += fmt.Sprintf("\n\n%s", param2)
		}
	}

	prompt += fmt.Sprintf("\n\n当前北京时间: %s", beijingTime.Format("2006-01-02 15:04:05"))
	prompt += fmt.Sprintf("\n开仓时间: %s", openTime.Format("2006-01-02 15:04:05"))

	return prompt, nil
}
