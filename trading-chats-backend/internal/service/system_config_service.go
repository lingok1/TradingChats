package service

import (
	"context"
	"encoding/json"
	"strings"
	"time"
	"trading-chats-backend/internal/db"
	"trading-chats-backend/internal/models"
	"trading-chats-backend/internal/repository"
)

const systemConfigCacheKey = "system_config:global"
const systemConfigCacheTTL = 24 * time.Hour

type SystemConfigService interface {
	GetConfig(ctx context.Context) (*models.SystemConfig, error)
	SaveBasicConfig(ctx context.Context, input *models.SaveSystemBasicConfigRequest) error
	SaveParameters(ctx context.Context, input *models.SaveSystemParametersRequest) error
	SaveRuntimeConfig(ctx context.Context, input *models.SaveSystemRuntimeConfigRequest) error
}

type systemConfigService struct {
	repo repository.SystemConfigRepository
}

func NewSystemConfigService(repo repository.SystemConfigRepository) SystemConfigService {
	return &systemConfigService{repo: repo}
}

func (s *systemConfigService) GetConfig(ctx context.Context) (*models.SystemConfig, error) {
	if cached, ok := s.getCachedConfig(ctx); ok {
		return normalizeSystemConfig(cached), nil
	}

	config, err := s.repo.GetConfig(ctx)
	if err != nil {
		return nil, err
	}

	normalized := normalizeSystemConfig(config)
	s.cacheConfig(ctx, normalized)
	return normalized, nil
}

func (s *systemConfigService) SaveBasicConfig(ctx context.Context, input *models.SaveSystemBasicConfigRequest) error {
	if input == nil {
		input = &models.SaveSystemBasicConfigRequest{}
	}

	current, err := s.repo.GetConfig(ctx)
	if err != nil {
		return err
	}

	next := normalizeSystemConfig(current)
	next.SystemTitle = input.SystemTitle
	next.SystemLogo = input.SystemLogo
	next.UpdatedAt = time.Now()

	if err := s.repo.SaveConfig(ctx, next); err != nil {
		return err
	}

	s.invalidateCache(ctx)
	return nil
}

func (s *systemConfigService) SaveParameters(ctx context.Context, input *models.SaveSystemParametersRequest) error {
	if input == nil {
		input = &models.SaveSystemParametersRequest{}
	}

	current, err := s.repo.GetConfig(ctx)
	if err != nil {
		return err
	}

	next := normalizeSystemConfig(current)
	next.Parameters = cloneParameters(input.Parameters)
	next.UpdatedAt = time.Now()

	if err := s.repo.SaveConfig(ctx, next); err != nil {
		return err
	}

	s.invalidateCache(ctx)
	return nil
}

func (s *systemConfigService) SaveRuntimeConfig(ctx context.Context, input *models.SaveSystemRuntimeConfigRequest) error {
	if input == nil {
		input = &models.SaveSystemRuntimeConfigRequest{}
	}

	current, err := s.repo.GetConfig(ctx)
	if err != nil {
		return err
	}

	next := normalizeSystemConfig(current)
	next.Param1 = normalizeRuntimeValue(input.Param1)
	next.Param2 = normalizeRuntimeValue(input.Param2)
	next.UpdatedAt = time.Now()

	if err := s.repo.SaveConfig(ctx, next); err != nil {
		return err
	}

	s.invalidateCache(ctx)
	return nil
}

func (s *systemConfigService) getCachedConfig(ctx context.Context) (*models.SystemConfig, bool) {
	if db.RedisClient == nil {
		return nil, false
	}

	payload, err := db.RedisClient.Get(ctx, systemConfigCacheKey).Result()
	if err != nil || payload == "" {
		return nil, false
	}

	var config models.SystemConfig
	if err := json.Unmarshal([]byte(payload), &config); err != nil {
		return nil, false
	}

	return &config, true
}

func (s *systemConfigService) cacheConfig(ctx context.Context, config *models.SystemConfig) {
	if db.RedisClient == nil || config == nil {
		return
	}

	payload, err := json.Marshal(config)
	if err != nil {
		return
	}

	_ = db.RedisClient.Set(ctx, systemConfigCacheKey, payload, systemConfigCacheTTL).Err()
}

func (s *systemConfigService) invalidateCache(ctx context.Context) {
	if db.RedisClient == nil {
		return
	}

	_ = db.RedisClient.Del(ctx, systemConfigCacheKey).Err()
}

func normalizeSystemConfig(config *models.SystemConfig) *models.SystemConfig {
	if config == nil {
		return &models.SystemConfig{
			ID:          models.GlobalSystemConfigID,
			SystemTitle: "Trading Chats",
			SystemLogo:  "",
			Param1:      "",
			Param2:      "",
			Parameters:  map[string]string{},
			UpdatedAt:   time.Now(),
		}
	}

	normalized := &models.SystemConfig{
		ID:          config.ID,
		SystemTitle: config.SystemTitle,
		SystemLogo:  config.SystemLogo,
		Param1:      normalizeRuntimeValue(config.Param1),
		Param2:      normalizeRuntimeValue(config.Param2),
		Parameters:  cloneParameters(config.Parameters),
		UpdatedAt:   config.UpdatedAt,
	}
	if normalized.ID == "" {
		normalized.ID = models.GlobalSystemConfigID
	}
	if normalized.Parameters == nil {
		normalized.Parameters = map[string]string{}
	}

	return normalized
}

func cloneParameters(source map[string]string) map[string]string {
	if source == nil {
		return map[string]string{}
	}

	cloned := make(map[string]string, len(source))
	for key, value := range source {
		cloned[key] = value
	}

	return cloned
}

func normalizeRuntimeValue(value string) string {
	return strings.TrimSpace(value)
}
