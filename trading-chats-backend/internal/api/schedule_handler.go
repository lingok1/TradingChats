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
	return &ScheduleHandler{
		service: service,
	}
}

// CreateConfig 创建定时任务配置
// @Summary 创建定时任务配置
// @Description 创建新的定时任务配置
// @Tags 定时任务
// @Accept json
// @Produce json
// @Param config body models.ScheduleConfig true "定时任务配置"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /schedules [post]
func (h *ScheduleHandler) CreateConfig(c *gin.Context) {
	var config models.ScheduleConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}

	if err := h.service.CreateConfig(c.Request.Context(), &config); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(config))
}

// GetConfigs 获取所有定时任务配置
// @Summary 获取所有定时任务配置
// @Description 获取所有定时任务配置列表
// @Tags 定时任务
// @Produce json
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /schedules [get]
func (h *ScheduleHandler) GetConfigs(c *gin.Context) {
	configs, err := h.service.GetConfigs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(configs))
}

// UpdateConfigStatus 更新任务状态 (暂停/恢复)
// @Summary 更新任务状态
// @Description 更新定时任务的状态（active/paused）
// @Tags 定时任务
// @Accept json
// @Produce json
// @Param id path string true "定时任务ID"
// @Param status body string true "状态 (active/paused)"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /schedules/{id}/status [put]
func (h *ScheduleHandler) UpdateConfigStatus(c *gin.Context) {
	id := c.Param("id")

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}

	if err := h.service.UpdateConfigStatus(c.Request.Context(), id, req.Status); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("status updated successfully"))
}

// DeleteConfig 删除定时任务配置
// @Summary 删除定时任务配置
// @Description 根据ID删除定时任务配置
// @Tags 定时任务
// @Produce json
// @Param id path string true "定时任务ID"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /schedules/{id} [delete]
func (h *ScheduleHandler) DeleteConfig(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteConfig(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("deleted successfully"))
}

// GetLogsByConfigID 获取某个任务的执行日志
// @Summary 获取任务执行日志
// @Description 根据定时任务ID获取执行记录列表
// @Tags 定时任务
// @Produce json
// @Param id path string true "定时任务ID"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /schedules/{id}/logs [get]
func (h *ScheduleHandler) GetLogsByConfigID(c *gin.Context) {
	id := c.Param("id")
	logs, err := h.service.GetLogsByConfigID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(logs))
}
