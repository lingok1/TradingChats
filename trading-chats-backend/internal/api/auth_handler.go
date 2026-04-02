package api

import (
	"net/http"
	"trading-chats-backend/internal/models"
	"trading-chats-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Register 用户注册
// @Summary 用户注册
// @Description 创建租户与租户管理员账号，并直接返回登录令牌
// @Tags 鉴权
// @Accept json
// @Produce json
// @Param body body models.RegisterRequest true "注册请求"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}

	resp, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(resp))
}

// Login 用户登录
// @Summary 用户登录
// @Description 使用用户名与密码获取访问令牌
// @Tags 鉴权
// @Accept json
// @Produce json
// @Param body body models.LoginRequest true "登录请求"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}

	resp, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse(401, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(resp))
}

// RefreshToken 刷新令牌
// @Summary 刷新令牌
// @Description 使用 refresh token 获取新的访问令牌
// @Tags 鉴权
// @Accept json
// @Produce json
// @Param body body models.RefreshTokenRequest true "刷新令牌请求"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 401 {object} models.Response
// @Router /api/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}

	resp, err := h.service.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, models.ErrorResponse(401, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(resp))
}

// Logout 用户登出
// @Summary 用户登出
// @Description 使用 refresh token 注销当前会话
// @Tags 鉴权
// @Accept json
// @Produce json
// @Param body body models.LogoutRequest true "登出请求"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req models.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}

	if err := h.service.Logout(c.Request.Context(), &req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("logout success"))
}
