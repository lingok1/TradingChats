package api

import (
	"time"
	"trading-chats-backend/internal/config"
	"trading-chats-backend/internal/models"
	"trading-chats-backend/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRoutes(
	r *gin.Engine,
	promptTemplateService *service.PromptTemplateService,
	modelAPIService *service.ModelAPIService,
	aiResponseService *service.AIResponseService,
	aiResponseEventService *service.AIResponseEventService,
	tradePlanService *service.TradePlanService,
	scheduleService *service.ScheduleService,
	systemConfigService service.SystemConfigService,
	authService *service.AuthService,
) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", Health)

	promptTemplateHandler := NewPromptTemplateHandler(promptTemplateService)
	modelAPIHandler := NewModelAPIHandler(modelAPIService)
	aiResponseHandler := NewAIResponseHandler(aiResponseService, aiResponseEventService)
	tradePlanHandler := NewTradePlanHandler(tradePlanService)
	scheduleHandler := NewScheduleHandler(scheduleService)
	systemConfigHandler := NewSystemConfigHandler(systemConfigService)
	authHandler := NewAuthHandler(authService)

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/reset-password", authHandler.ResetPassword)
		}

		promptTemplates := api.Group("/prompt-templates")
		{
			promptTemplates.GET("", promptTemplateHandler.GetAllPromptTemplates)
			promptTemplates.GET("/tag", promptTemplateHandler.GetPromptTemplatesByTag)
			promptTemplates.GET("/:id", promptTemplateHandler.GetPromptTemplateByID)
		}

		modelAPIConfigs := api.Group("/model-api-configs")
		{
			modelAPIConfigs.GET("", modelAPIHandler.GetAllModelAPIConfigs)
			modelAPIConfigs.GET("/provider", modelAPIHandler.GetModelAPIConfigsByProvider)
			modelAPIConfigs.GET("/:id", modelAPIHandler.GetModelAPIConfigByID)
		}

		aiResponses := api.Group("/ai-responses")
		{
			aiResponses.GET("", aiResponseHandler.GetAllAIResponses)
			aiResponses.GET("/batch", aiResponseHandler.GetAIResponsesByBatchID)
			aiResponses.GET("/events", aiResponseHandler.StreamAIResponseEvents)
			aiResponses.GET("/latest", aiResponseHandler.GetLatestBatch)
			aiResponses.GET("/:id", aiResponseHandler.GetAIResponseByID)
		}

		tradePlans := api.Group("/trade-plans")
		{
			tradePlans.GET("", AuthMiddleware(authService), tradePlanHandler.GetTradePlans)
			tradePlans.GET("/:id", AuthMiddleware(authService), tradePlanHandler.GetTradePlanByID)
		}

		schedules := api.Group("/schedules")
		{
			schedules.GET("", scheduleHandler.GetConfigs)
			schedules.GET("/:id/logs", scheduleHandler.GetLogsByConfigID)
		}

		systemConfig := api.Group("/system-config")
		{
			systemConfig.GET("", systemConfigHandler.GetConfig)
		}
	}

	protected := r.Group("/api")
	protected.Use(AuthMiddleware(authService))
	protected.Use(TenantIDMiddleware())
	{
		protected.POST("/prompt-templates", RequireRoles(models.RoleAdmin, models.RoleTenant), promptTemplateHandler.CreatePromptTemplate)
		protected.PUT("/prompt-templates/:id", RequireRoles(models.RoleAdmin, models.RoleTenant), promptTemplateHandler.UpdatePromptTemplate)
		protected.DELETE("/prompt-templates/:id", RequireRoles(models.RoleAdmin, models.RoleTenant), promptTemplateHandler.DeletePromptTemplate)
		protected.POST("/prompt-templates/generate", RequireRoles(models.RoleAdmin, models.RoleTenant), promptTemplateHandler.GeneratePrompt)

		protected.POST("/model-api-configs", RequireRoles(models.RoleAdmin, models.RoleTenant), modelAPIHandler.CreateModelAPIConfig)
		protected.PUT("/model-api-configs/:id", RequireRoles(models.RoleAdmin, models.RoleTenant), modelAPIHandler.UpdateModelAPIConfig)
		protected.DELETE("/model-api-configs/:id", RequireRoles(models.RoleAdmin, models.RoleTenant), modelAPIHandler.DeleteModelAPIConfig)
		protected.POST("/model-api-configs/:id/test", RequireRoles(models.RoleAdmin, models.RoleTenant), modelAPIHandler.TestModelConnectivity)

		protected.POST("/ai-responses/generate", RequireRoles(models.RoleAdmin, models.RoleTenant), aiResponseHandler.GenerateBatchAIResponses)

		protected.POST("/trade-plans", RequireRoles(models.RoleAdmin, models.RoleTenant), tradePlanHandler.CreateTradePlan)
		protected.PUT("/trade-plans/:id", RequireRoles(models.RoleAdmin, models.RoleTenant), tradePlanHandler.UpdateTradePlan)
		protected.DELETE("/trade-plans/:id", RequireRoles(models.RoleAdmin, models.RoleTenant), tradePlanHandler.DeleteTradePlan)

		protected.POST("/schedules", RequireRoles(models.RoleAdmin, models.RoleTenant), scheduleHandler.CreateConfig)
		protected.PUT("/schedules/:id", RequireRoles(models.RoleAdmin, models.RoleTenant), scheduleHandler.UpdateConfig)
		protected.PUT("/schedules/status", RequireRoles(models.RoleAdmin, models.RoleTenant), scheduleHandler.UpdateConfigStatus)
		protected.DELETE("/schedules/:id", RequireRoles(models.RoleAdmin, models.RoleTenant), scheduleHandler.DeleteConfig)
		protected.POST("/schedules/:id/trigger", RequireRoles(models.RoleAdmin, models.RoleTenant), scheduleHandler.TriggerNow)

		protected.PUT("/system-config/basic", RequireRoles(models.RoleAdmin), systemConfigHandler.SaveBasicConfig)
		protected.PUT("/system-config/parameters", RequireRoles(models.RoleAdmin, models.RoleTenant), systemConfigHandler.SaveParameters)
	}
}

func SetupRouter(
	promptTemplateService *service.PromptTemplateService,
	modelAPIService *service.ModelAPIService,
	aiResponseService *service.AIResponseService,
	aiResponseEventService *service.AIResponseEventService,
	tradePlanService *service.TradePlanService,
	scheduleService *service.ScheduleService,
	systemConfigService service.SystemConfigService,
	authService *service.AuthService,
	swaggerConfig *config.SwaggerConfig,
) *gin.Engine {
	r := gin.Default()
	SetupRoutes(r, promptTemplateService, modelAPIService, aiResponseService, aiResponseEventService, tradePlanService, scheduleService, systemConfigService, authService)

	r.GET("/swagger/*any", SwaggerBasicAuth(swaggerConfig), ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/swagger/doc.json"),
	))

	return r
}
