package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"trading-chats-backend/internal/models"
	"trading-chats-backend/internal/repository"
	"trading-chats-backend/pkg/utils"

	"github.com/robfig/cron/v3"
)

const (
	TriggerTypeManual = "manual"
	TriggerTypeAuto   = "auto"
)

type ScheduleService struct {
	repo                  *repository.ScheduleRepository
	aiResponseService     *AIResponseService
	promptTemplateService *PromptTemplateService
	cronEngine            *cron.Cron
	cronParser            cron.Parser
	taskMap               map[string]cron.EntryID
	mutex                 sync.RWMutex
}

func NewScheduleService(repo *repository.ScheduleRepository, aiResponseService *AIResponseService, promptTemplateService *PromptTemplateService) *ScheduleService {
	parser := cron.NewParser(
		cron.SecondOptional |
			cron.Minute |
			cron.Hour |
			cron.Dom |
			cron.Month |
			cron.Dow |
			cron.Descriptor,
	)
	c := cron.New(cron.WithParser(parser))
	return &ScheduleService{
		repo:                  repo,
		aiResponseService:     aiResponseService,
		promptTemplateService: promptTemplateService,
		cronEngine:            c,
		cronParser:            parser,
		taskMap:               make(map[string]cron.EntryID),
	}
}

func normalizeScheduleConfig(config *models.ScheduleConfig) {
	if config == nil {
		return
	}
	config.TabTag = models.NormalizeTabTag(config.TabTag)
	config.CronExpr = strings.TrimSpace(config.CronExpr)
}

func (s *ScheduleService) validateCronExpr(expr string) error {
	if strings.TrimSpace(expr) == "" {
		return errors.New("cron_expr is required")
	}
	if _, err := s.cronParser.Parse(expr); err != nil {
		return fmt.Errorf("invalid cron_expr: %w", err)
	}
	return nil
}

func scheduleConfigSnapshot(config models.ScheduleConfig) map[string]interface{} {
	return map[string]interface{}{
		"name":        config.Name,
		"cron_expr":   config.CronExpr,
		"template_id": config.TemplateID,
		"tab_tag":     config.TabTag,
		"status":      config.Status,
	}
}

func mergeScheduleConfig(base models.ScheduleConfig, updateData map[string]interface{}) models.ScheduleConfig {
	updated := base

	if name, ok := updateData["name"].(string); ok {
		updated.Name = name
	}
	if cronExpr, ok := updateData["cron_expr"].(string); ok {
		updated.CronExpr = cronExpr
	}
	if templateID, ok := updateData["template_id"].(string); ok {
		updated.TemplateID = templateID
	}
	if tabTag, ok := updateData["tab_tag"].(string); ok {
		updated.TabTag = tabTag
	}
	if status, ok := updateData["status"].(string); ok {
		updated.Status = status
	}

	return updated
}

func (s *ScheduleService) Start(ctx context.Context) error {
	s.cronEngine.Start()
	configs, err := s.repo.GetActiveConfigs(ctx)
	if err != nil {
		return err
	}
	for _, config := range configs {
		if err := s.addTaskToCron(config); err != nil {
			log.Printf("Failed to add task %s to cron: %v\n", config.ID.Hex(), err)
		}
	}
	log.Println("Schedule Service started successfully.")
	return nil
}

func (s *ScheduleService) Stop() {
	s.cronEngine.Stop()
}

func (s *ScheduleService) addTaskToCron(config models.ScheduleConfig) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	config.CronExpr = strings.TrimSpace(config.CronExpr)
	if err := s.validateCronExpr(config.CronExpr); err != nil {
		return err
	}
	if entryID, exists := s.taskMap[config.ID.Hex()]; exists {
		s.cronEngine.Remove(entryID)
	}
	entryID, err := s.cronEngine.AddFunc(config.CronExpr, func() {
		authCtx := &models.AuthContext{TenantID: config.TenantID, Role: models.RoleTenant}
		ctx := models.WithAuthContext(context.Background(), authCtx)
		s.executeTask(ctx, config, TriggerTypeAuto)
	})
	if err != nil {
		return err
	}
	s.taskMap[config.ID.Hex()] = entryID
	return nil
}

func (s *ScheduleService) removeTaskFromCron(configID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if entryID, exists := s.taskMap[configID]; exists {
		s.cronEngine.Remove(entryID)
		delete(s.taskMap, configID)
	}
}

func (s *ScheduleService) executeTask(ctx context.Context, config models.ScheduleConfig, triggerType string) {
	taskCtx, cancel := context.WithTimeout(ctx, 8*time.Minute)
	defer cancel()
	config.TabTag = models.NormalizeTabTag(config.TabTag)
	log.Printf("Executing scheduled task: %s (TemplateID: %s, TriggerType: %s)\n", config.Name, config.TemplateID, triggerType)

	prompt, err := s.promptTemplateService.GeneratePrompt(taskCtx, config.TemplateID)
	if err != nil {
		log.Printf("Failed to generate prompt: %v\n", err)
		logEntry := &models.ScheduleLog{
			TenantID:         config.TenantID,
			ScheduleConfigID: config.ID,
			TabTag:           config.TabTag,
			Prompt:           "",
			TriggerType:      triggerType,
			Status:           "failed",
			Error:            err.Error(),
			ExecutedAt:       utils.NowString(),
		}
		if err := s.repo.CreateLog(taskCtx, logEntry); err != nil {
			log.Printf("Failed to save schedule log: %v\n", err)
		}
		return
	}

	batchID, err := s.aiResponseService.GenerateBatchAIResponses(taskCtx, config.TemplateID, config.TabTag)
	logEntry := &models.ScheduleLog{
		TenantID:         config.TenantID,
		ScheduleConfigID: config.ID,
		TabTag:           config.TabTag,
		Prompt:           prompt,
		TriggerType:      triggerType,
		ExecutedAt:       utils.NowString(),
	}
	if err != nil {
		log.Printf("Task execution failed: %v\n", err)
		logEntry.Status = "failed"
		logEntry.Error = err.Error()
	} else {
		log.Printf("Task execution succeeded, BatchID: %s\n", batchID)
		logEntry.Status = "success"
		logEntry.BatchID = batchID
	}
	if err := s.repo.CreateLog(taskCtx, logEntry); err != nil {
		log.Printf("Failed to save schedule log: %v\n", err)
	}
}

func (s *ScheduleService) TriggerNow(ctx context.Context, id string) error {
	config, err := s.repo.GetConfigByID(ctx, id)
	if err != nil {
		return err
	}
	normalizeScheduleConfig(config)
	authCtx := cloneAuthContext(models.GetAuthContext(ctx))
	go func(cfg models.ScheduleConfig, auth *models.AuthContext) {
		bg := context.Background()
		if auth != nil {
			bg = models.WithAuthContext(bg, auth)
		}
		s.executeTask(bg, cfg, TriggerTypeManual)
	}(*config, authCtx)
	return nil
}

func (s *ScheduleService) CreateConfig(ctx context.Context, config *models.ScheduleConfig) error {
	if config.Status == "" {
		config.Status = "active"
	}
	normalizeScheduleConfig(config)
	if config.Status == "active" {
		if err := s.validateCronExpr(config.CronExpr); err != nil {
			return err
		}
	}
	if err := s.repo.CreateConfig(ctx, config); err != nil {
		return err
	}
	if config.Status == "active" {
		if err := s.addTaskToCron(*config); err != nil {
			if rollbackErr := s.repo.DeleteConfig(ctx, config.ID.Hex()); rollbackErr != nil {
				return fmt.Errorf("failed to register schedule task: %v; rollback delete failed: %w", err, rollbackErr)
			}
			return fmt.Errorf("failed to register schedule task: %w", err)
		}
	}
	return nil
}

func (s *ScheduleService) UpdateConfigStatus(ctx context.Context, id string, status string) error {
	if status != "active" && status != "paused" {
		return errors.New("invalid status, must be 'active' or 'paused'")
	}
	config, err := s.repo.GetConfigByID(ctx, id)
	if err != nil {
		return err
	}
	normalizeScheduleConfig(config)
	previousStatus := config.Status
	if status == "active" {
		if err := s.validateCronExpr(config.CronExpr); err != nil {
			return err
		}
	}
	if err := s.repo.UpdateConfig(ctx, id, map[string]interface{}{"status": status}); err != nil {
		return err
	}
	config.Status = status
	if status == "active" {
		if err := s.addTaskToCron(*config); err != nil {
			rollbackStatus := previousStatus
			if rollbackStatus == "" {
				rollbackStatus = "paused"
			}
			if rollbackErr := s.repo.UpdateConfig(ctx, id, map[string]interface{}{"status": rollbackStatus}); rollbackErr != nil {
				return fmt.Errorf("failed to register schedule task: %v; rollback failed: %w", err, rollbackErr)
			}
			if rollbackStatus != "active" {
				s.removeTaskFromCron(id)
			}
			return fmt.Errorf("failed to register schedule task: %w", err)
		}
		return nil
	}
	s.removeTaskFromCron(id)
	return nil
}

func (s *ScheduleService) DeleteConfig(ctx context.Context, id string) error {
	if err := s.repo.DeleteConfig(ctx, id); err != nil {
		return err
	}
	s.removeTaskFromCron(id)
	return nil
}

func (s *ScheduleService) GetConfigs(ctx context.Context) ([]models.ScheduleConfig, error) {
	return s.repo.GetAllConfigs(ctx)
}

func (s *ScheduleService) UpdateConfig(ctx context.Context, id string, updateData map[string]interface{}) error {
	if tabTag, ok := updateData["tab_tag"].(string); ok {
		updateData["tab_tag"] = models.NormalizeTabTag(tabTag)
	}
	if cronExpr, ok := updateData["cron_expr"].(string); ok {
		updateData["cron_expr"] = strings.TrimSpace(cronExpr)
	}
	if status, ok := updateData["status"].(string); ok && status != "active" && status != "paused" {
		return errors.New("invalid status, must be 'active' or 'paused'")
	}
	currentConfig, err := s.repo.GetConfigByID(ctx, id)
	if err != nil {
		return err
	}
	normalizeScheduleConfig(currentConfig)

	nextConfig := mergeScheduleConfig(*currentConfig, updateData)
	normalizeScheduleConfig(&nextConfig)
	if nextConfig.Status == "active" {
		if err := s.validateCronExpr(nextConfig.CronExpr); err != nil {
			return err
		}
	}

	if err := s.repo.UpdateConfig(ctx, id, updateData); err != nil {
		return err
	}

	if nextConfig.Status == "active" {
		if err := s.addTaskToCron(nextConfig); err != nil {
			if rollbackErr := s.repo.UpdateConfig(ctx, id, scheduleConfigSnapshot(*currentConfig)); rollbackErr != nil {
				return fmt.Errorf("failed to register schedule task: %v; rollback failed: %w", err, rollbackErr)
			}
			if currentConfig.Status == "active" {
				if restoreErr := s.addTaskToCron(*currentConfig); restoreErr != nil {
					return fmt.Errorf("failed to restore previous schedule task after rollback: %w", restoreErr)
				}
			} else {
				s.removeTaskFromCron(id)
			}
			return fmt.Errorf("failed to register schedule task: %w", err)
		}
		return nil
	}
	s.removeTaskFromCron(id)
	return nil
}

func (s *ScheduleService) GetLogsByConfigID(ctx context.Context, configID string) ([]models.ScheduleLog, error) {
	return s.repo.GetLogsByConfigID(ctx, configID)
}
