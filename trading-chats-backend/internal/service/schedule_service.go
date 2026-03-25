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

type ScheduleService struct {
	repo              *repository.ScheduleRepository
	aiResponseService *AIResponseService
	cronEngine        *cron.Cron
	taskMap           map[string]cron.EntryID // 映射 ConfigID -> EntryID
	mutex             sync.RWMutex
}

func NewScheduleService(repo *repository.ScheduleRepository, aiResponseService *AIResponseService) *ScheduleService {
	// 使用秒级 cron (如果不加这个选项，默认是分钟级别)
	c := cron.New(cron.WithSeconds())
	return &ScheduleService{
		repo:              repo,
		aiResponseService: aiResponseService,
		cronEngine:        c,
		taskMap:           make(map[string]cron.EntryID),
	}
}

// Start 启动定时任务调度器，并加载所有激活的任务
func (s *ScheduleService) Start(ctx context.Context) error {
	s.cronEngine.Start()

	// 从数据库加载激活状态的任务
	configs, err := s.repo.GetActiveConfigs(ctx)
	if err != nil {
		return err
	}

	for _, config := range configs {
		err := s.addTaskToCron(config)
		if err != nil {
			log.Printf("Failed to add task %s to cron: %v\n", config.ID.Hex(), err)
		}
	}

	log.Println("Schedule Service started successfully.")
	return nil
}

// Stop 停止调度器
func (s *ScheduleService) Stop() {
	s.cronEngine.Stop()
}

// addTaskToCron 将任务添加到执行引擎中
func (s *ScheduleService) addTaskToCron(config models.ScheduleConfig) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	// 如果已经在运行，先移除
	if entryID, exists := s.taskMap[config.ID.Hex()]; exists {
		s.cronEngine.Remove(entryID)
	}

	entryID, err := s.cronEngine.AddFunc(config.CronExpr, func() {
		s.executeTask(config)
	})
	if err != nil {
		return err
	}

	s.taskMap[config.ID.Hex()] = entryID
	return nil
}

// removeTaskFromCron 从执行引擎中移除任务
func (s *ScheduleService) removeTaskFromCron(configID string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if entryID, exists := s.taskMap[configID]; exists {
		s.cronEngine.Remove(entryID)
		delete(s.taskMap, configID)
	}
}

// executeTask 实际执行任务的逻辑
func (s *ScheduleService) executeTask(config models.ScheduleConfig) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	log.Printf("Executing scheduled task: %s (TemplateID: %s)\n", config.Name, config.TemplateID)

	// 调用 AI Response 服务生成内容
	batchID, err := s.aiResponseService.GenerateBatchAIResponses(ctx, config.TemplateID, config.Param1, config.Param2)

	// 记录执行日志
	logEntry := &models.ScheduleLog{
		ScheduleConfigID: config.ID,
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

	// 保存到数据库
	err = s.repo.CreateLog(context.Background(), logEntry)
	if err != nil {
		log.Printf("Failed to save schedule log: %v\n", err)
	}
}

// CreateConfig 创建并可能启动定时任务
func (s *ScheduleService) CreateConfig(ctx context.Context, config *models.ScheduleConfig) error {
	if config.Status == "" {
		config.Status = "active"
	}

	err := s.repo.CreateConfig(ctx, config)
	if err != nil {
		return err
	}

	if config.Status == "active" {
		return s.addTaskToCron(*config)
	}
	return nil
}

// UpdateConfigStatus 更改任务状态（暂停/恢复）
func (s *ScheduleService) UpdateConfigStatus(ctx context.Context, id string, status string) error {
	if status != "active" && status != "paused" {
		return errors.New("invalid status, must be 'active' or 'paused'")
	}

	err := s.repo.UpdateConfig(ctx, id, map[string]interface{}{"status": status})
	if err != nil {
		return err
	}

	config, err := s.repo.GetConfigByID(ctx, id)
	if err != nil {
		return err
	}

	if status == "active" {
		return s.addTaskToCron(*config)
	} else {
		s.removeTaskFromCron(id)
	}

	return nil
}

// DeleteConfig 删除配置及对应的引擎任务
func (s *ScheduleService) DeleteConfig(ctx context.Context, id string) error {
	err := s.repo.DeleteConfig(ctx, id)
	if err != nil {
		return err
	}

	s.removeTaskFromCron(id)
	return nil
}

// GetConfigs 获取所有配置
func (s *ScheduleService) GetConfigs(ctx context.Context) ([]models.ScheduleConfig, error) {
	return s.repo.GetAllConfigs(ctx)
}

// GetLogsByConfigID 获取任务的执行日志
func (s *ScheduleService) GetLogsByConfigID(ctx context.Context, configID string) ([]models.ScheduleLog, error) {
	return s.repo.GetLogsByConfigID(ctx, configID)
}
