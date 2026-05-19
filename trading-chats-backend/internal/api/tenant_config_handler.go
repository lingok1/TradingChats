package api

import (
	"net/http"
	"trading-chats-backend/internal/models"
	"trading-chats-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type TenantConfigHandler struct {
	service *service.TenantConfigService
}

func NewTenantConfigHandler(svc *service.TenantConfigService) *TenantConfigHandler {
	return &TenantConfigHandler{service: svc}
}

// GetConfig 获取租户配置
// @Summary 获取租户配置
// @Tags 租户配置
// @Produce json
// @Security BearerAuth
// @Param tenant_id query string false "租户ID（管理员专用）"
// @Success 200 {object} models.Response
// @Router /api/tenant-config [get]
func (h *TenantConfigHandler) GetConfig(c *gin.Context) {
	authCtx := MustGetAuthContext(c)
	tenantID := authCtx.TenantID
	if models.IsAdmin(authCtx) && c.Query("tenant_id") != "" {
		tenantID = c.Query("tenant_id")
	}
	cfg, err := h.service.GetConfig(c.Request.Context(), tenantID)
	if err != nil {
		c.JSON(http.StatusForbidden, models.ErrorResponse(403, err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse(cfg))
}

// SaveMenu 管理员设置租户菜单可见性
// @Summary 设置租户菜单
// @Tags 租户配置
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tenant_id query string true "租户ID"
// @Param body body models.TenantMenuConfig true "菜单配置"
// @Success 200 {object} models.Response
// @Router /api/tenant-config/menu [put]
func (h *TenantConfigHandler) SaveMenu(c *gin.Context) {
	tenantID := c.Query("tenant_id")
	if tenantID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, "tenant_id is required"))
		return
	}
	var menu models.TenantMenuConfig
	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}
	if err := h.service.SaveMenu(c.Request.Context(), tenantID, menu); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("ok"))
}

// SaveParameters 保存租户动态参数
// @Summary 保存租户动态参数
// @Tags 租户配置
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param tenant_id query string false "租户ID（管理员专用，不传则用当前租户）"
// @Param body body models.SaveSystemParametersRequest true "动态参数"
// @Success 200 {object} models.Response
// @Router /api/tenant-config/parameters [put]
func (h *TenantConfigHandler) SaveParameters(c *gin.Context) {
	authCtx := MustGetAuthContext(c)
	tenantID := authCtx.TenantID
	if models.IsAdmin(authCtx) && c.Query("tenant_id") != "" {
		tenantID = c.Query("tenant_id")
	}
	var req models.SaveSystemParametersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}
	if err := h.service.SaveParameters(c.Request.Context(), tenantID, req.Parameters); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("ok"))
}

// ListConfigs 管理员获取所有租户配置
// @Summary 获取所有租户配置
// @Tags 租户配置
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.Response
// @Router /api/tenant-config/list [get]
func (h *TenantConfigHandler) ListConfigs(c *gin.Context) {
	configs, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse(configs))
}
