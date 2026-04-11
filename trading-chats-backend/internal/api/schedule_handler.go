package api

import (
	"net/http"

	"trading-chats-backend/internal/models"
	"trading-chats-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	service *service.ScheduleService
}

func NewScheduleHandler(service *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{service: service}
}

// CreateConfig 创建定时任务
// @Summary 创建定时任务
// @Description 创建新的定时任务配置
// @Tags 定时任务
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.ScheduleConfig true "定时任务配置"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/schedules [post]
func (h *ScheduleHandler) CreateConfig(c *gin.Context) {
	var config models.ScheduleConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}

	// 设置租户ID
	authCtx := MustGetAuthContext(c)
	if authCtx != nil {
		config.TenantID = authCtx.TenantID
	}

	if err := h.service.CreateConfig(c.Request.Context(), &config); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(config))
}

// GetConfigs 获取定时任务列表
// @Summary 获取定时任务列表
// @Description 获取当前可见的定时任务配置列表
// @Tags 定时任务
// @Produce json
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/schedules [get]
func (h *ScheduleHandler) GetConfigs(c *gin.Context) {
	configs, err := h.service.GetConfigs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(configs))
}

// UpdateConfigStatus 更新定时任务状态
// @Summary 更新定时任务状态
// @Description 将定时任务切换为 active 或 paused
// @Tags 定时任务
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.UpdateStatusRequest true "状态更新请求"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/schedules/status [put]
func (h *ScheduleHandler) UpdateConfigStatus(c *gin.Context) {
	var req models.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}

	if err := h.service.UpdateConfigStatus(c.Request.Context(), req.ID, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("status updated successfully"))
}

// DeleteConfig 删除定时任务
// @Summary 删除定时任务
// @Description 根据 ID 删除定时任务配置
// @Tags 定时任务
// @Produce json
// @Security BearerAuth
// @Param id path string true "定时任务ID"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/schedules/{id} [delete]
func (h *ScheduleHandler) DeleteConfig(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteConfig(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("deleted successfully"))
}

// GetLogsByConfigID 获取定时任务日志
// @Summary 获取定时任务日志
// @Description 根据任务 ID 获取执行日志
// @Tags 定时任务
// @Produce json
// @Param id path string true "定时任务ID"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/schedules/{id}/logs [get]
func (h *ScheduleHandler) GetLogsByConfigID(c *gin.Context) {
	id := c.Param("id")
	logs, err := h.service.GetLogsByConfigID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(logs))
}

// UpdateConfig 更新定时任务
// @Summary 更新定时任务
// @Description 更新指定定时任务的配置
// @Tags 定时任务
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "定时任务ID"
// @Param body body models.ScheduleConfig true "定时任务配置"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/schedules/{id} [put]
func (h *ScheduleHandler) UpdateConfig(c *gin.Context) {
	id := c.Param("id")
	var config models.ScheduleConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}

	// 更新配置
	updateData := map[string]interface{}{
		"name":        config.Name,
		"cron_expr":   config.CronExpr,
		"template_id": config.TemplateID,
		"status":      config.Status,
	}

	if err := h.service.UpdateConfig(c.Request.Context(), id, updateData); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("updated successfully"))
}

// TriggerNow 手动触发定时任务
// @Summary 手动触发定时任务
// @Description 立即执行指定定时任务
// @Tags 定时任务
// @Produce json
// @Security BearerAuth
// @Param id path string true "定时任务ID"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/schedules/{id}/trigger [post]
func (h *ScheduleHandler) TriggerNow(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.TriggerNow(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("triggered successfully"))
}
