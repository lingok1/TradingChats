package main

import (
	"context"
	"log"
	"trading-chats-backend/internal/api"
	"trading-chats-backend/internal/config"
	"trading-chats-backend/internal/db"
	"trading-chats-backend/internal/repository"
	"trading-chats-backend/internal/service"

	_ "trading-chats-backend/docs"
)

// @title Trading Chats Backend API
// @version 1.0
// @description Trading Chats 后端服务接口文档
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 请输入 JWT Token，格式为: Bearer <token>
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := db.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Disconnect()

	promptTemplateRepo := repository.NewPromptTemplateRepository(db.MongoDB)
	modelAPIRepo := repository.NewModelAPIRepository(db.MongoDB)
	aiResponseRepo := repository.NewAIResponseRepository(db.MongoDB)
	scheduleRepo := repository.NewScheduleRepository(db.MongoDB)
	systemConfigRepo := repository.NewSystemConfigRepository(db.MongoDB)
	authRepo := repository.NewAuthRepository(db.MongoDB)

	systemConfigService := service.NewSystemConfigService(systemConfigRepo)
	promptTemplateService := service.NewPromptTemplateService(promptTemplateRepo, systemConfigService)
	modelAPIService := service.NewModelAPIService(modelAPIRepo)
	aiResponseService := service.NewAIResponseService(aiResponseRepo, modelAPIService, promptTemplateService, systemConfigService)
	scheduleService := service.NewScheduleService(scheduleRepo, aiResponseService, promptTemplateService)
	authService := service.NewAuthService(authRepo, cfg.JWT)

	if err := authService.EnsureBootstrapData(context.Background()); err != nil {
		log.Fatalf("Failed to bootstrap auth data: %v", err)
	}

	if err := scheduleService.Start(context.Background()); err != nil {
		log.Printf("Failed to start schedule service: %v", err)
	}
	defer scheduleService.Stop()

	r := api.SetupRouter(promptTemplateService, modelAPIService, aiResponseService, scheduleService, systemConfigService, authService, &cfg.Swagger)

	log.Printf("Server running at http://localhost:%s", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
