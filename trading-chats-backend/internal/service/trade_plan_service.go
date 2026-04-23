package service

import (
	"context"
	"errors"
	"strings"

	"trading-chats-backend/internal/models"
	"trading-chats-backend/internal/repository"
)

type TradePlanService struct {
	repo *repository.TradePlanRepository
}

func NewTradePlanService(repo *repository.TradePlanRepository) *TradePlanService {
	return &TradePlanService{repo: repo}
}

func normalizeTradePlanTab(tabTag string) (string, error) {
	if tabTag == "" {
		return models.TabTagFutures, nil
	}
	switch tabTag {
	case models.TabTagFutures, models.TabTagOptions:
		return tabTag, nil
	default:
		return "", errors.New("tab_tag must be futures or options")
	}
}

func normalizeTradePlan(plan *models.TradePlan) error {
	if plan == nil {
		return errors.New("trade plan is required")
	}
	normalizedTab, err := normalizeTradePlanTab(plan.TabTag)
	if err != nil {
		return err
	}
	plan.TabTag = normalizedTab
	plan.Name = strings.TrimSpace(plan.Name)
	plan.Symbol = strings.TrimSpace(plan.Symbol)
	plan.Strategy = strings.TrimSpace(plan.Strategy)
	plan.Direction = strings.TrimSpace(plan.Direction)
	plan.OpenTime = strings.TrimSpace(plan.OpenTime)
	plan.CloseTime = strings.TrimSpace(plan.CloseTime)
	plan.Remark = strings.TrimSpace(plan.Remark)
	if plan.Status == "" {
		plan.Status = models.TradePlanStatusPlanned
	}
	return nil
}

func validateTradePlan(plan *models.TradePlan) error {
	if plan == nil {
		return errors.New("trade plan is required")
	}
	if plan.Symbol == "" {
		return errors.New("symbol is required")
	}
	if plan.Direction == "" {
		return errors.New("direction is required")
	}
	if plan.EntryPrice <= 0 {
		return errors.New("entry_price must be greater than 0")
	}
	if plan.TakeProfit <= 0 {
		return errors.New("take_profit must be greater than 0")
	}
	if plan.StopLoss <= 0 {
		return errors.New("stop_loss must be greater than 0")
	}

	switch plan.Status {
	case models.TradePlanStatusPlanned, models.TradePlanStatusActive, models.TradePlanStatusClosed, models.TradePlanStatusCancelled:
		return nil
	default:
		return errors.New("status must be planned, active, closed or cancelled")
	}
}

func (s *TradePlanService) EnsureIndexes(ctx context.Context) error {
	return s.repo.EnsureIndexes(ctx)
}

func (s *TradePlanService) CreateTradePlan(ctx context.Context, plan *models.TradePlan) error {
	if err := normalizeTradePlan(plan); err != nil {
		return err
	}
	if err := validateTradePlan(plan); err != nil {
		return err
	}
	return s.repo.Create(ctx, plan)
}

func (s *TradePlanService) GetTradePlans(ctx context.Context, tabTag string) ([]models.TradePlan, error) {
	normalizedTab, err := normalizeTradePlanTab(tabTag)
	if err != nil {
		return nil, err
	}
	return s.repo.GetAll(ctx, normalizedTab)
}

func (s *TradePlanService) GetTradePlanByID(ctx context.Context, id string) (*models.TradePlan, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *TradePlanService) UpdateTradePlan(ctx context.Context, plan *models.TradePlan) error {
	if err := normalizeTradePlan(plan); err != nil {
		return err
	}
	if err := validateTradePlan(plan); err != nil {
		return err
	}
	return s.repo.Update(ctx, plan)
}

func (s *TradePlanService) DeleteTradePlan(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
