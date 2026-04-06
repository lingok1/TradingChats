package api

import (
	"net/http"

	"trading-chats-backend/internal/models"
	"trading-chats-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type SystemConfigHandler struct {
	service service.SystemConfigService
}

func NewSystemConfigHandler(service service.SystemConfigService) *SystemConfigHandler {
	return &SystemConfigHandler{service: service}
}

// GetConfig 获取系统配置
// @Summary 获取系统配置
// @Description 获取当前系统配置
// @Tags 系统配置
// @Produce json
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/system-config [get]
func (h *SystemConfigHandler) GetConfig(c *gin.Context) {
	config, err := h.service.GetConfig(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse(config))
}

// SaveBasicConfig 保存基础系统配置
// @Summary 保存基础系统配置
// @Description 保存系统标题与 Logo 等基础配置
// @Tags 系统配置
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.SaveSystemBasicConfigRequest true "基础配置"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/system-config/basic [put]
func (h *SystemConfigHandler) SaveBasicConfig(c *gin.Context) {
	var req models.SaveSystemBasicConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}

	if err := h.service.SaveBasicConfig(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("ok"))
}

// SaveParameters 保存动态参数配置
// @Summary 保存动态参数配置
// @Description 保存系统动态参数字典
// @Tags 系统配置
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.SaveSystemParametersRequest true "动态参数配置"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/system-config/parameters [put]
func (h *SystemConfigHandler) SaveParameters(c *gin.Context) {
	var req models.SaveSystemParametersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}

	if err := h.service.SaveParameters(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("ok"))
}
