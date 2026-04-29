package service

import (
	"testing"

	"trading-chats-backend/internal/models"
)

func TestNormalizeTradePlanTabAllowsStock(t *testing.T) {
	tabTag, err := normalizeTradePlanTab(models.TabTagStock)
	if err != nil {
		t.Fatalf("normalizeTradePlanTab returned error: %v", err)
	}

	if tabTag != models.TabTagStock {
		t.Fatalf("expected %q, got %q", models.TabTagStock, tabTag)
	}
}

func TestValidateTradePlanAllowsStockPlan(t *testing.T) {
	plan := &models.TradePlan{
		TabTag:     models.TabTagStock,
		Symbol:     "600519",
		Direction:  "buy",
		EntryPrice: 1500,
		TakeProfit: 1680,
		StopLoss:   1420,
		Status:     models.TradePlanStatusPlanned,
	}

	if err := normalizeTradePlan(plan); err != nil {
		t.Fatalf("normalizeTradePlan returned error: %v", err)
	}

	if err := validateTradePlan(plan); err != nil {
		t.Fatalf("validateTradePlan returned error: %v", err)
	}
}

func TestNormalizeTradePlanTabRejectsUnsupportedTab(t *testing.T) {
	_, err := normalizeTradePlanTab(models.TabTagNews)
	if err == nil {
		t.Fatal("expected unsupported tab_tag error")
	}
}
