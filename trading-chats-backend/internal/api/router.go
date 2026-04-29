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

	r.GET("/", BackendInfoHandler)
	r.GET("/health", Health)

	promptTemplateHandler := NewPromptTemplateHandler(promptTemplateService)
	modelAPIHandler := NewModelAPIHandler(modelAPIService)
	aiResponseHandler := NewAIResponseHandler(aiResponseService, aiResponseEventService)
	tradePlanHandler := NewTradePlanHandler(tradePlanService)
	scheduleHandler := NewScheduleHandler(scheduleService)
	systemConfigHandler := NewSystemConfigHandler(systemConfigService)
	authHandler := NewAuthHandler(authService)

	// 期权功能当前复用通用分析链路，不单独拆 /api/options 前缀：
	// 1. 前端在期权页使用 tab_tag=options 调用 /api/ai-responses/latest 读取最新完成批次。
	// 2. 前端通过 /api/ai-responses/events?tab_tag=options 建立 SSE，模型响应完成后收到 ai_response_updated 再刷新 latest。
	// 3. 期权定时任务写在 /api/schedules，任务字段包含 template_id、cron_expr、tab_tag=options、status。
	// 4. 后端启动或启用任务时只把 status=active 的任务注册到 cron；paused 任务只落库不自动运行。
	// 5. cron/manual 触发后 ScheduleService 先用 template_id 生成提示词，再调用 AIResponseService。
	// 6. AIResponseService 按 tab_tag=options 筛选已启用的模型 API 配置，并将响应落到 ai_responses_options 集合。
	// 7. 每次调度执行结果额外写入 schedule_logs，可通过 /api/schedules/:id/logs 查询批次号、状态、错误和执行时间。
	// 股票功能也沿用同一模式，使用 tab_tag=stock、ai_responses_stock 和模型 tab_settings 中的 stock 开关。
	// 只有股票存在差异化领域行为时再新增专用路由组。
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
			// 提示词模板是期权定时任务的输入来源；ScheduleConfig.template_id 指向这里的模板。
			// GeneratePrompt 会将系统参数中的占位符替换为静态值或 URL 返回的 JSON，并追加当前北京时间。
			promptTemplates.GET("", promptTemplateHandler.GetAllPromptTemplates)
			promptTemplates.GET("/tag", promptTemplateHandler.GetPromptTemplatesByTag)
			promptTemplates.GET("/:id", promptTemplateHandler.GetPromptTemplateByID)
		}

		modelAPIConfigs := api.Group("/model-api-configs")
		{
			// 模型配置通过 tab_settings 控制每个分析页签是否启用。
			// 期权执行时只选择 tab_settings 中 tab_tag=options 且 enabled=true 的模型配置；股票同理使用 tab_tag=stock。
			modelAPIConfigs.GET("", modelAPIHandler.GetAllModelAPIConfigs)
			modelAPIConfigs.GET("/provider", modelAPIHandler.GetModelAPIConfigsByProvider)
			modelAPIConfigs.GET("/:id", modelAPIHandler.GetModelAPIConfigByID)
		}

		aiResponses := api.Group("/ai-responses")
		{
			// AI 响应接口用 query 参数 tab_tag 区分业务页签；期权传 tab_tag=options，股票传 tab_tag=stock。
			// /latest 读取最新完成批次；/batch 按批次查询；/events 是前端实时刷新用的 SSE 通道。
			// 后端根据 tab_tag 将期权响应落库到 ai_responses_options，股票响应落库到 ai_responses_stock。
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
			// 定时任务列表和执行日志。期权/股票任务本身不需要专用接口，通过 tab_tag 标识。
			// 列表用于前端任务面板展示，日志用于查看每次执行的 batch_id、success/failed 和 error。
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
		// 手动生成提示词预览：只返回 prompt，不调用模型、不写 ai_responses。
		protected.POST("/prompt-templates/generate", RequireRoles(models.RoleAdmin, models.RoleTenant), promptTemplateHandler.GeneratePrompt)

		protected.POST("/model-api-configs", RequireRoles(models.RoleAdmin, models.RoleTenant), modelAPIHandler.CreateModelAPIConfig)
		protected.PUT("/model-api-configs/:id", RequireRoles(models.RoleAdmin, models.RoleTenant), modelAPIHandler.UpdateModelAPIConfig)
		protected.DELETE("/model-api-configs/:id", RequireRoles(models.RoleAdmin, models.RoleTenant), modelAPIHandler.DeleteModelAPIConfig)
		protected.POST("/model-api-configs/:id/test", RequireRoles(models.RoleAdmin, models.RoleTenant), modelAPIHandler.TestModelConnectivity)

		// 手动触发 AI 分析：请求体携带 template_id 和 tab_tag=options。
		// 返回 batch_id 后，前端仍通过 SSE/latest 或 /batch 查询结果。
		protected.POST("/ai-responses/generate", RequireRoles(models.RoleAdmin, models.RoleTenant), aiResponseHandler.GenerateBatchAIResponses)

		protected.POST("/trade-plans", RequireRoles(models.RoleAdmin, models.RoleTenant), tradePlanHandler.CreateTradePlan)
		protected.PUT("/trade-plans/:id", RequireRoles(models.RoleAdmin, models.RoleTenant), tradePlanHandler.UpdateTradePlan)
		protected.DELETE("/trade-plans/:id", RequireRoles(models.RoleAdmin, models.RoleTenant), tradePlanHandler.DeleteTradePlan)

		protected.POST("/schedules", RequireRoles(models.RoleAdmin, models.RoleTenant), scheduleHandler.CreateConfig)
		protected.PUT("/schedules/:id", RequireRoles(models.RoleAdmin, models.RoleTenant), scheduleHandler.UpdateConfig)
		// active 会立即注册到 cron；paused 会从 cron 移除但保留配置和历史日志。
		protected.PUT("/schedules/status", RequireRoles(models.RoleAdmin, models.RoleTenant), scheduleHandler.UpdateConfigStatus)
		protected.DELETE("/schedules/:id", RequireRoles(models.RoleAdmin, models.RoleTenant), scheduleHandler.DeleteConfig)
		// 手动触发指定定时任务，走同一条 executeTask 链路，日志 trigger_type=manual。
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
