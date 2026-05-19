package service

import (
	"context"
	"errors"
	"trading-chats-backend/internal/models"
	"trading-chats-backend/internal/repository"
)

type TenantConfigService struct {
	repo             *repository.TenantConfigRepository
	systemConfigSvc  SystemConfigService
}

func NewTenantConfigService(repo *repository.TenantConfigRepository, systemConfigSvc SystemConfigService) *TenantConfigService {
	return &TenantConfigService{repo: repo, systemConfigSvc: systemConfigSvc}
}

func (s *TenantConfigService) GetConfig(ctx context.Context, tenantID string) (*models.TenantConfig, error) {
	authCtx := models.GetAuthContext(ctx)
	if !models.IsAdmin(authCtx) {
		if authCtx == nil || authCtx.TenantID != tenantID {
			return nil, errors.New("permission denied")
		}
	}
	return s.repo.GetByTenantID(ctx, tenantID)
}

func (s *TenantConfigService) SaveMenu(ctx context.Context, tenantID string, menu models.TenantMenuConfig) error {
	cfg, err := s.repo.GetByTenantID(ctx, tenantID)
	if err != nil {
		return err
	}
	cfg.MenuConfig = menu
	return s.repo.Save(ctx, cfg)
}

func (s *TenantConfigService) SaveParameters(ctx context.Context, tenantID string, params map[string]string) error {
	cfg, err := s.repo.GetByTenantID(ctx, tenantID)
	if err != nil {
		return err
	}
	cfg.Parameters = params
	return s.repo.Save(ctx, cfg)
}

// GetParameters 优先取租户参数，无则 fallback 全局参数
func (s *TenantConfigService) GetParameters(ctx context.Context) (map[string]string, error) {
	authCtx := models.GetAuthContext(ctx)
	if authCtx != nil && authCtx.TenantID != "" {
		cfg, err := s.repo.GetByTenantID(ctx, authCtx.TenantID)
		if err == nil && len(cfg.Parameters) > 0 {
			return cfg.Parameters, nil
		}
	}
	sysConfig, err := s.systemConfigSvc.GetConfig(ctx)
	if err != nil {
		return map[string]string{}, nil
	}
	return sysConfig.Parameters, nil
}

func (s *TenantConfigService) GetAll(ctx context.Context) ([]models.TenantConfig, error) {
	return s.repo.GetAll(ctx)
}
