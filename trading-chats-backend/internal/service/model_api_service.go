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

func normalizeModelAPITabSettings(settings []models.ModelAPITabSetting) []models.ModelAPITabSetting {
	if len(settings) == 0 {
		return nil
	}

	normalized := make([]models.ModelAPITabSetting, 0, len(settings))
	indexByTag := make(map[string]int, len(settings))

	for _, setting := range settings {
		tabTag := models.NormalizeTabTag(setting.TabTag)
		if index, ok := indexByTag[tabTag]; ok {
			normalized[index].Enabled = setting.Enabled
			continue
		}

		indexByTag[tabTag] = len(normalized)
		normalized = append(normalized, models.ModelAPITabSetting{
			TabTag:  tabTag,
			Enabled: setting.Enabled,
		})
	}

	return normalized
}

func allTabSettings() []models.ModelAPITabSetting {
	return []models.ModelAPITabSetting{
		{TabTag: models.TabTagFutures, Enabled: true},
		{TabTag: models.TabTagOptions, Enabled: true},
		{TabTag: models.TabTagNews, Enabled: true},
		{TabTag: models.TabTagPosition, Enabled: true},
	}
}

func normalizeModelAPIConfig(config *models.ModelAPIConfig) {
	if config == nil {
		return
	}

	config.TabSettings = normalizeModelAPITabSettings(config.TabSettings)
	if len(config.TabSettings) == 0 {
		config.TabSettings = allTabSettings()
	}
}

func normalizeModelAPIConfigs(configs []models.ModelAPIConfig) []models.ModelAPIConfig {
	for index := range configs {
		normalizeModelAPIConfig(&configs[index])
	}
	return configs
}

func isTabEnabled(config *models.ModelAPIConfig, tabTag string) bool {
	if config == nil {
		return false
	}

	normalizedTab := models.NormalizeTabTag(tabTag)
	for _, setting := range config.TabSettings {
		if setting.TabTag == normalizedTab && setting.Enabled {
			return true
		}
	}

	return false
}

func (s *ModelAPIService) EnsureTabSettings(ctx context.Context) error {
	configs, err := s.repo.GetAll(ctx)
	if err != nil {
		return err
	}

	for index := range configs {
		needsBackfill := len(configs[index].TabSettings) == 0
		normalizeModelAPIConfig(&configs[index])
		if needsBackfill {
			if err := s.repo.BackfillTabSettings(ctx, &configs[index]); err != nil {
				return err
			}
			continue
		}

		existing := make(map[string]bool, len(configs[index].TabSettings))
		for _, setting := range configs[index].TabSettings {
			existing[setting.TabTag] = true
		}

		missing := false
		for _, tag := range []string{models.TabTagFutures, models.TabTagOptions, models.TabTagNews, models.TabTagPosition} {
			if !existing[tag] {
				configs[index].TabSettings = append(configs[index].TabSettings, models.ModelAPITabSetting{
					TabTag:  tag,
					Enabled: false,
				})
				missing = true
			}
		}
		if missing {
			if err := s.repo.BackfillTabSettings(ctx, &configs[index]); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *ModelAPIService) CreateModelAPIConfig(ctx context.Context, config *models.ModelAPIConfig) error {
	normalizeModelAPIConfig(config)
	return s.repo.Create(ctx, config)
}

func (s *ModelAPIService) GetModelAPIConfigByID(ctx context.Context, id string) (*models.ModelAPIConfig, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %w", err)
	}
	config, err := s.repo.GetByID(ctx, objectID)
	if err != nil {
		return nil, err
	}
	normalizeModelAPIConfig(config)
	return config, nil
}

func (s *ModelAPIService) GetAllModelAPIConfigs(ctx context.Context) ([]models.ModelAPIConfig, error) {
	configs, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return normalizeModelAPIConfigs(configs), nil
}

func (s *ModelAPIService) GetEnabledModelAPIConfigsByTab(ctx context.Context, tabTag string) ([]models.ModelAPIConfig, error) {
	normalizedTab := models.NormalizeTabTag(tabTag)
	configs, err := s.repo.GetEnabledByTabTag(ctx, normalizedTab)
	if err != nil {
		return nil, err
	}

	configs = normalizeModelAPIConfigs(configs)
	enabled := make([]models.ModelAPIConfig, 0, len(configs))
	for _, config := range configs {
		if isTabEnabled(&config, normalizedTab) {
			enabled = append(enabled, config)
		}
	}

	return enabled, nil
}

func (s *ModelAPIService) GetModelAPIConfigsByProvider(ctx context.Context, provider string) ([]models.ModelAPIConfig, error) {
	configs, err := s.repo.GetByProvider(ctx, provider)
	if err != nil {
		return nil, err
	}
	return normalizeModelAPIConfigs(configs), nil
}

func (s *ModelAPIService) UpdateModelAPIConfig(ctx context.Context, config *models.ModelAPIConfig) error {
	normalizeModelAPIConfig(config)
	return s.repo.Update(ctx, config)
}

func (s *ModelAPIService) DeleteModelAPIConfig(ctx context.Context, id string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("invalid ID format: %w", err)
	}
	return s.repo.Delete(ctx, objectID)
}

func (s *ModelAPIService) TestModelConnectivity(ctx context.Context, configID string) error {
	config, err := s.GetModelAPIConfigByID(ctx, configID)
	if err != nil {
		return fmt.Errorf("failed to get model API config: %w", err)
	}

	var req *http.Request
	var errRequest error

	switch config.Provider {
	case "anthropic":
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
