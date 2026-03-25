package service

import (
	"context"
	"time"
	"trading-chats-backend/internal/models"
	"trading-chats-backend/internal/repository"
)

type SystemConfigService interface {
	GetConfig(ctx context.Context) (*models.SystemConfig, error)
	SaveBasicConfig(ctx context.Context, input *models.SaveSystemBasicConfigRequest) error
	SaveParameters(ctx context.Context, input *models.SaveSystemParametersRequest) error
}

type systemConfigService struct {
	repo repository.SystemConfigRepository
}

func NewSystemConfigService(repo repository.SystemConfigRepository) SystemConfigService {
	return &systemConfigService{repo: repo}
}

func (s *systemConfigService) GetConfig(ctx context.Context) (*models.SystemConfig, error) {
	config, err := s.repo.GetConfig(ctx)
	if err != nil {
		return nil, err
	}

	return normalizeSystemConfig(config), nil
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

	return s.repo.SaveConfig(ctx, next)
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

	return s.repo.SaveConfig(ctx, next)
}

func normalizeSystemConfig(config *models.SystemConfig) *models.SystemConfig {
	if config == nil {
		return &models.SystemConfig{
			ID:          models.GlobalSystemConfigID,
			SystemTitle: "Trading Chats",
			SystemLogo:  "",
			Parameters:  map[string]string{},
			UpdatedAt:   time.Now(),
		}
	}

	normalized := &models.SystemConfig{
		ID:          config.ID,
		SystemTitle: config.SystemTitle,
		SystemLogo:  config.SystemLogo,
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
