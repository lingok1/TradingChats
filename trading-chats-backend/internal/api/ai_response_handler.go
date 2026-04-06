package api

import (
	"net/http"

	"trading-chats-backend/internal/models"
	"trading-chats-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AIResponseHandler struct {
	service *service.AIResponseService
}

func NewAIResponseHandler(service *service.AIResponseService) *AIResponseHandler {
	return &AIResponseHandler{service: service}
}

// GetAIResponseByID 获取 AI 响应详情
// @Summary 获取 AI 响应详情
// @Description 根据 ID 获取 AI 响应详情
// @Tags AI响应
// @Produce json
// @Param id path string true "AI 响应ID"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/ai-responses/{id} [get]
func (h *AIResponseHandler) GetAIResponseByID(c *gin.Context) {
	id := c.Param("id")
	response, err := h.service.GetAIResponseByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(response))
}

// GetAIResponsesByBatchID 按批次获取 AI 响应
// @Summary 按批次获取 AI 响应
// @Description 根据 batch_id 获取 AI 响应列表
// @Tags AI响应
// @Produce json
// @Param batch_id query string true "批次ID"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/ai-responses/batch [get]
func (h *AIResponseHandler) GetAIResponsesByBatchID(c *gin.Context) {
	batchID := c.Query("batch_id")
	responses, err := h.service.GetAIResponsesByBatchID(c.Request.Context(), batchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(responses))
}

// GetLatestSuccessfulBatch 获取最近成功批次
// @Summary 获取最近成功批次
// @Description 获取最近成功生成的 AI 响应批次
// @Tags AI响应
// @Produce json
// @Success 200 {object} models.Response
// @Failure 404 {object} models.Response
// @Router /api/ai-responses/latest [get]
func (h *AIResponseHandler) GetLatestSuccessfulBatch(c *gin.Context) {
	responses, err := h.service.GetLatestSuccessfulBatch(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse(404, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(responses))
}

// GetAllAIResponses 获取 AI 响应列表
// @Summary 获取 AI 响应列表
// @Description 获取当前可见的 AI 响应列表
// @Tags AI响应
// @Produce json
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/ai-responses [get]
func (h *AIResponseHandler) GetAllAIResponses(c *gin.Context) {
	responses, err := h.service.GetAllAIResponses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(responses))
}

// GenerateBatchAIResponses 生成 AI 响应
// @Summary 生成 AI 响应
// @Description 根据模板生成一批 AI 响应
// @Tags AI响应
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.GenerateAIRequest true "生成请求"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/ai-responses/generate [post]
func (h *AIResponseHandler) GenerateBatchAIResponses(c *gin.Context) {
	var req models.GenerateAIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}

	if req.TemplateID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, "template_id is required"))
		return
	}

	batchID, err := h.service.GenerateBatchAIResponses(c.Request.Context(), req.TemplateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(map[string]string{"batch_id": batchID}))
}
