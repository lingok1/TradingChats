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

// GetConfig godoc
// @Summary 获取系统配置
// @Description 获取全局系统配置（包括标题、Logo、动态参数等）
// @Tags system-config
// @Accept json
// @Produce json
// @Success 200 {object} models.Response{data=models.SystemConfig}
// @Failure 500 {object} models.Response
// @Router /system-config [get]
func (h *SystemConfigHandler) GetConfig(c *gin.Context) {
	config, err := h.service.GetConfig(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, "Failed to get system config: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(config))
}

// SaveBasicConfig godoc
// @Summary 保存系统基础配置
// @Description 更新全局系统标题和Logo
// @Tags system-config
// @Accept json
// @Produce json
// @Param config body models.SaveSystemBasicConfigRequest true "系统基础配置"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /system-config/basic [put]
func (h *SystemConfigHandler) SaveBasicConfig(c *gin.Context) {
	var config models.SaveSystemBasicConfigRequest
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, "Invalid request parameters"))
		return
	}

	err := h.service.SaveBasicConfig(c.Request.Context(), &config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, "Failed to save system basic config: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(map[string]string{"message": "System basic config saved successfully"}))
}

// SaveParameters godoc
// @Summary 保存系统动态参数
// @Description 更新全局动态参数配置
// @Tags system-config
// @Accept json
// @Produce json
// @Param config body models.SaveSystemParametersRequest true "系统动态参数配置"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /system-config/parameters [put]
func (h *SystemConfigHandler) SaveParameters(c *gin.Context) {
	var config models.SaveSystemParametersRequest
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, "Invalid request parameters"))
		return
	}

	err := h.service.SaveParameters(c.Request.Context(), &config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, "Failed to save system parameters: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(map[string]string{"message": "System parameters saved successfully"}))
}
