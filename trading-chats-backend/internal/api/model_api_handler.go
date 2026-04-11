package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"trading-chats-backend/internal/models"
	"trading-chats-backend/internal/service"
)

type ModelAPIHandler struct {
	service *service.ModelAPIService
}

func NewModelAPIHandler(service *service.ModelAPIService) *ModelAPIHandler {
	return &ModelAPIHandler{
		service: service,
	}
}

// CreateModelAPIConfig 创建模型与API配置
// @Summary 创建模型与API配置
// @Description 创建新的模型与API配置
// @Tags 模型与API配置
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param config body models.ModelAPIConfig true "模型与API配置信息"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/model-api-configs [post]
func (h *ModelAPIHandler) CreateModelAPIConfig(c *gin.Context) {
	var config models.ModelAPIConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}

	// 设置租户ID
	authCtx := MustGetAuthContext(c)
	if authCtx != nil {
		config.TenantID = authCtx.TenantID
	}

	if err := h.service.CreateModelAPIConfig(c.Request.Context(), &config); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusCreated, models.SuccessResponse(config))
}

// GetModelAPIConfigByID 根据ID获取模型与API配置
// @Summary 根据ID获取模型与API配置
// @Description 根据ID获取模型与API配置详情
// @Tags 模型与API配置
// @Produce json
// @Param id path string true "模型与API配置ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/model-api-configs/{id} [get]
func (h *ModelAPIHandler) GetModelAPIConfigByID(c *gin.Context) {
	id := c.Param("id")
	config, err := h.service.GetModelAPIConfigByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(config))
}

// GetAllModelAPIConfigs 获取所有模型与API配置
// @Summary 获取所有模型与API配置
// @Description 获取所有模型与API配置列表
// @Tags 模型与API配置
// @Produce json
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/model-api-configs [get]
func (h *ModelAPIHandler) GetAllModelAPIConfigs(c *gin.Context) {
	configs, err := h.service.GetAllModelAPIConfigs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(configs))
}

// GetModelAPIConfigsByProvider 根据提供商获取模型与API配置
// @Summary 根据提供商获取模型与API配置
// @Description 根据提供商获取模型与API配置列表
// @Tags 模型与API配置
// @Produce json
// @Param provider query string true "提供商"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/model-api-configs/provider [get]
func (h *ModelAPIHandler) GetModelAPIConfigsByProvider(c *gin.Context) {
	provider := c.Query("provider")
	configs, err := h.service.GetModelAPIConfigsByProvider(c.Request.Context(), provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(configs))
}

// UpdateModelAPIConfig 更新模型与API配置
// @Summary 更新模型与API配置
// @Description 更新模型与API配置信息
// @Tags 模型与API配置
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "模型与API配置ID"
// @Param config body models.ModelAPIConfig true "模型与API配置信息"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/model-api-configs/{id} [put]
func (h *ModelAPIHandler) UpdateModelAPIConfig(c *gin.Context) {
	id := c.Param("id")
	var config models.ModelAPIConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}

	config.ID, _ = models.ParseObjectID(id)

	// 设置租户ID
	authCtx := MustGetAuthContext(c)
	if authCtx != nil {
		config.TenantID = authCtx.TenantID
	}

	if err := h.service.UpdateModelAPIConfig(c.Request.Context(), &config); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(config))
}

// DeleteModelAPIConfig 删除模型与API配置
// @Summary 删除模型与API配置
// @Description 删除模型与API配置
// @Tags 模型与API配置
// @Produce json
// @Security BearerAuth
// @Param id path string true "模型与API配置ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/model-api-configs/{id} [delete]
func (h *ModelAPIHandler) DeleteModelAPIConfig(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteModelAPIConfig(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(map[string]string{"message": "Model API config deleted successfully"}))
}

// TestModelConnectivity 测试模型的连通性
// @Summary 测试模型的连通性
// @Description 测试模型与API的连通性
// @Tags 模型与API配置
// @Produce json
// @Security BearerAuth
// @Param id path string true "模型与API配置ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/model-api-configs/{id}/test [post]
func (h *ModelAPIHandler) TestModelConnectivity(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.TestModelConnectivity(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
	
	return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(map[string]string{"message": "Model connectivity test passed"}))
}
