package service

import (
	"context"
	"errors"
	"log"
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
	taskMap               map[string]cron.EntryID
	mutex                 sync.RWMutex
}

func NewScheduleService(repo *repository.ScheduleRepository, aiResponseService *AIResponseService, promptTemplateService *PromptTemplateService) *ScheduleService {
	c := cron.New(cron.WithSeconds())
	return &ScheduleService{
		repo:                  repo,
		aiResponseService:     aiResponseService,
		promptTemplateService: promptTemplateService,
		cronEngine:            c,
		taskMap:               make(map[string]cron.EntryID),
	}
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
	log.Printf("Executing scheduled task: %s (TemplateID: %s, TriggerType: %s)\n", config.Name, config.TemplateID, triggerType)

	prompt, err := s.promptTemplateService.GeneratePrompt(taskCtx, config.TemplateID)
	if err != nil {
		log.Printf("Failed to generate prompt: %v\n", err)
		logEntry := &models.ScheduleLog{
			TenantID:         config.TenantID,
			ScheduleConfigID: config.ID,
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

	batchID, err := s.aiResponseService.GenerateBatchAIResponses(taskCtx, config.TemplateID)
	logEntry := &models.ScheduleLog{
		TenantID:         config.TenantID,
		ScheduleConfigID: config.ID,
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
	go s.executeTask(ctx, *config, TriggerTypeManual)
	return nil
}

func (s *ScheduleService) CreateConfig(ctx context.Context, config *models.ScheduleConfig) error {
	if config.Status == "" {
		config.Status = "active"
	}
	if err := s.repo.CreateConfig(ctx, config); err != nil {
		return err
	}
	if config.Status == "active" {
		return s.addTaskToCron(*config)
	}
	return nil
}

func (s *ScheduleService) UpdateConfigStatus(ctx context.Context, id string, status string) error {
	if status != "active" && status != "paused" {
		return errors.New("invalid status, must be 'active' or 'paused'")
	}
	if err := s.repo.UpdateConfig(ctx, id, map[string]interface{}{"status": status}); err != nil {
		return err
	}
	config, err := s.repo.GetConfigByID(ctx, id)
	if err != nil {
		return err
	}
	if status == "active" {
		return s.addTaskToCron(*config)
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
	if err := s.repo.UpdateConfig(ctx, id, updateData); err != nil {
		return err
	}
	config, err := s.repo.GetConfigByID(ctx, id)
	if err != nil {
		return err
	}
	if config.Status == "active" {
		return s.addTaskToCron(*config)
	}
	s.removeTaskFromCron(id)
	return nil
}

func (s *ScheduleService) GetLogsByConfigID(ctx context.Context, configID string) ([]models.ScheduleLog, error) {
	return s.repo.GetLogsByConfigID(ctx, configID)
}
