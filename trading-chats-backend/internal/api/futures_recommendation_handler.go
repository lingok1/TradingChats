package api

import (
	"net/http"
	"strconv"

	"trading-chats-backend/internal/models"
	"trading-chats-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type FuturesRecommendationHandler struct {
	service *service.FuturesRecommendationService
}

func NewFuturesRecommendationHandler(svc *service.FuturesRecommendationService) *FuturesRecommendationHandler {
	return &FuturesRecommendationHandler{service: svc}
}

// GetLatest 获取最新推荐
// @Summary 获取最新优选推荐
// @Tags 推荐
// @Produce json
// @Param tab_tag query string false "tab类型（futures/options/stock），不传则返回所有类型中最新一条"
// @Success 200 {object} models.Response
// @Router /api/futures-recommendation/latest [get]
func (h *FuturesRecommendationHandler) GetLatest(c *gin.Context) {
	tabTag := c.Query("tab_tag")
	rec, err := h.service.GetLatestByTab(c.Request.Context(), tabTag)
	if err != nil {
		c.JSON(http.StatusOK, models.SuccessResponse(nil))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse(rec))
}

// GetList 获取推荐历史列表
// @Summary 获取期货推荐历史
// @Tags 期货推荐
// @Produce json
// @Param limit query int false "条数限制，默认20"
// @Success 200 {object} models.Response
// @Router /api/futures-recommendation [get]
func (h *FuturesRecommendationHandler) GetList(c *gin.Context) {
	limit := int64(20)
	if l := c.Query("limit"); l != "" {
		if n, err := strconv.ParseInt(l, 10, 64); err == nil && n > 0 {
			limit = n
		}
	}
	list, err := h.service.GetList(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse(list))
}

// Generate 手动触发生成推荐
// @Summary 手动触发期货推荐生成
// @Tags 期货推荐
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body GenerateRecommendationRequest true "模型配置"
// @Success 200 {object} models.Response
// @Router /api/futures-recommendation/generate [post]
func (h *FuturesRecommendationHandler) Generate(c *gin.Context) {
	var req GenerateRecommendationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}
	if err := h.service.Generate(c.Request.Context(), req.ModelAPIID, req.ModelName); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, models.SuccessResponse("ok"))
}

type GenerateRecommendationRequest struct {
	ModelAPIID string `json:"model_api_id" binding:"required"`
	ModelName  string `json:"model_name" binding:"required"`
}
