package api

import (
	"net/http"
	"trading-chats-backend/internal/models"

	"github.com/gin-gonic/gin"
)

// Health 健康检查
// @Summary 健康检查
// @Description 返回服务健康状态
// @Tags 系统
// @Produce json
// @Success 200 {object} models.Response
// @Router /health [get]
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, models.SuccessResponse(map[string]string{"status": "ok"}))
}
