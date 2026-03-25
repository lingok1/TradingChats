package api

import (
	"trading-chats-backend/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter 配置路由
func SetupRouter(
	promptTemplateService *service.PromptTemplateService,
	modelAPIService *service.ModelAPIService,
	aiResponseService *service.AIResponseService,
	scheduleService *service.ScheduleService,
	systemConfigService service.SystemConfigService,
) *gin.Engine {
	r := gin.Default()

	// 配置CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// Swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API路由组
	api := r.Group("/api")
	{
		// 提示词模版路由
		promptTemplateHandler := NewPromptTemplateHandler(promptTemplateService)
		promptTemplates := api.Group("/prompt-templates")
		{
			promptTemplates.POST("", promptTemplateHandler.CreatePromptTemplate)
			promptTemplates.GET("", promptTemplateHandler.GetAllPromptTemplates)
			promptTemplates.GET("/tag", promptTemplateHandler.GetPromptTemplatesByTag)
			promptTemplates.GET("/:id", promptTemplateHandler.GetPromptTemplateByID)
			promptTemplates.PUT("/:id", promptTemplateHandler.UpdatePromptTemplate)
			promptTemplates.DELETE("/:id", promptTemplateHandler.DeletePromptTemplate)
			promptTemplates.POST("/generate", promptTemplateHandler.GeneratePrompt)
		}

		// 模型与API配置路由
		modelAPIHandler := NewModelAPIHandler(modelAPIService)
		modelAPIConfigs := api.Group("/model-api-configs")
		{
			modelAPIConfigs.POST("", modelAPIHandler.CreateModelAPIConfig)
			modelAPIConfigs.GET("", modelAPIHandler.GetAllModelAPIConfigs)
			modelAPIConfigs.GET("/provider", modelAPIHandler.GetModelAPIConfigsByProvider)
			modelAPIConfigs.GET("/:id", modelAPIHandler.GetModelAPIConfigByID)
			modelAPIConfigs.PUT("/:id", modelAPIHandler.UpdateModelAPIConfig)
			modelAPIConfigs.DELETE("/:id", modelAPIHandler.DeleteModelAPIConfig)
			modelAPIConfigs.POST("/:id/test", modelAPIHandler.TestModelConnectivity)
		}

		// AI响应信息路由
		aiResponseHandler := NewAIResponseHandler(aiResponseService)
		aiResponses := api.Group("/ai-responses")
		{
			aiResponses.GET("", aiResponseHandler.GetAllAIResponses)
			aiResponses.GET("/batch", aiResponseHandler.GetAIResponsesByBatchID)
			aiResponses.GET("/latest", aiResponseHandler.GetLatestSuccessfulBatch)
			aiResponses.GET("/:id", aiResponseHandler.GetAIResponseByID)
			aiResponses.POST("/generate", aiResponseHandler.GenerateBatchAIResponses)
		}

		// 定时任务配置路由
		scheduleHandler := NewScheduleHandler(scheduleService)
		schedules := api.Group("/schedules")
		{
			schedules.POST("", scheduleHandler.CreateConfig)
			schedules.GET("", scheduleHandler.GetConfigs)
			schedules.PUT("/:id/status", scheduleHandler.UpdateConfigStatus)
			schedules.DELETE("/:id", scheduleHandler.DeleteConfig)
			schedules.GET("/:id/logs", scheduleHandler.GetLogsByConfigID)
		}

		// 系统配置路由
		systemConfigHandler := NewSystemConfigHandler(systemConfigService)
		systemConfig := api.Group("/system-config")
		{
			systemConfig.GET("", systemConfigHandler.GetConfig)
			systemConfig.PUT("/basic", systemConfigHandler.SaveBasicConfig)
			systemConfig.PUT("/parameters", systemConfigHandler.SaveParameters)
		}
	}

	return r
}
