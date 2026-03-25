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
	return &AIResponseHandler{
		service: service,
	}
}

// GetAIResponseByID 根据ID获取AI响应信息
// @Summary 根据ID获取AI响应信息
// @Description 根据ID获取AI响应信息详情
// @Tags AI响应信息
// @Produce json
// @Param id path string true "AI响应信息ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /ai-responses/{id} [get]
func (h *AIResponseHandler) GetAIResponseByID(c *gin.Context) {
	id := c.Param("id")
	response, err := h.service.GetAIResponseByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(response))
}

// GetAIResponsesByBatchID 根据批次ID获取AI响应信息
// @Summary 根据批次ID获取AI响应信息
// @Description 根据批次ID获取AI响应信息列表
// @Tags AI响应信息
// @Produce json
// @Param batch_id query string true "批次ID"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /ai-responses/batch [get]
func (h *AIResponseHandler) GetAIResponsesByBatchID(c *gin.Context) {
	batchID := c.Query("batch_id")
	responses, err := h.service.GetAIResponsesByBatchID(c.Request.Context(), batchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(responses))
}

// GetLatestSuccessfulBatch 获取最近一次成功的批次数据
// @Summary 获取最近一次成功的批次数据
// @Description 获取最近一次状态为 completed 的完整批次 AI 响应列表
// @Tags AI响应信息
// @Produce json
// @Success 200 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /ai-responses/latest [get]
func (h *AIResponseHandler) GetLatestSuccessfulBatch(c *gin.Context) {
	responses, err := h.service.GetLatestSuccessfulBatch(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse(404, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(responses))
}

// GetAllAIResponses 获取所有AI响应信息
// @Summary 获取所有AI响应信息
// @Description 获取所有AI响应信息列表
// @Tags AI响应信息
// @Produce json
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /ai-responses [get]
func (h *AIResponseHandler) GetAllAIResponses(c *gin.Context) {
	responses, err := h.service.GetAllAIResponses(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(responses))
}

// GenerateBatchAIResponses 生成批次AI响应
// @Summary 生成批次AI响应
// @Description 根据提示词模版和参数生成批次AI响应
// @Tags AI响应信息
// @Accept json
// @Produce json
// @Param request body models.GenerateAIRequest true "生成批次AI响应请求"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /ai-responses/generate [post]
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

	batchID, err := h.service.GenerateBatchAIResponses(c.Request.Context(), req.TemplateID, req.Param1, req.Param2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(map[string]string{"batch_id": batchID}))
}
