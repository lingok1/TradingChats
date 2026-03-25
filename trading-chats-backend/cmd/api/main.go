package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "trading-chats-backend/docs"

	"trading-chats-backend/internal/api"
	"trading-chats-backend/internal/config"
	"trading-chats-backend/internal/db"
	"trading-chats-backend/internal/repository"
	"trading-chats-backend/internal/service"
)

// 生成Swagger文档
// $ swag init -g cmd/api/main.go

// @title 凌期AI辅助期货挑选品种和持仓止盈止损分析程序
// @version 1.0
// @description 凌期AI辅助期货挑选品种和持仓止盈止损分析程序API文档
// @host localhost:8080
// @BasePath /api
func main() {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 连接数据库
	if err := db.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to databases: %v", err)
	}
	defer db.Disconnect()

	// 初始化仓库
	mongoDB := db.MongoClient.Database(cfg.MongoDB.Database)
	promptTemplateRepo := repository.NewPromptTemplateRepository(mongoDB)
	modelAPIRepo := repository.NewModelAPIRepository(mongoDB)
	aiResponseRepo := repository.NewAIResponseRepository(mongoDB)
	scheduleRepo := repository.NewScheduleRepository(mongoDB)
	systemConfigRepo := repository.NewSystemConfigRepository(mongoDB)

	// 初始化服务
	systemConfigService := service.NewSystemConfigService(systemConfigRepo)
	promptTemplateService := service.NewPromptTemplateService(promptTemplateRepo, systemConfigService)
	modelAPIService := service.NewModelAPIService(modelAPIRepo)
	aiResponseService := service.NewAIResponseService(aiResponseRepo, modelAPIService, promptTemplateService, systemConfigService)
	scheduleService := service.NewScheduleService(scheduleRepo, aiResponseService)

	// 启动定时任务调度器
	if err := scheduleService.Start(context.Background()); err != nil {
		log.Printf("Failed to start schedule service: %v", err)
	}
	defer scheduleService.Stop()

	// 设置路由
	router := api.SetupRouter(promptTemplateService, modelAPIService, aiResponseService, scheduleService, systemConfigService)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", cfg.Server.Port),
		Handler: router,
	}

	// 启动服务器
	go func() {
		log.Printf("Server is running on port %s", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 优雅关闭服务器
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
