package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"trading-chats-backend/internal/models"
	"trading-chats-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type AIResponseHandler struct {
	service      *service.AIResponseService
	eventService *service.AIResponseEventService
}

func NewAIResponseHandler(
	service *service.AIResponseService,
	eventService *service.AIResponseEventService,
) *AIResponseHandler {
	return &AIResponseHandler{
		service:      service,
		eventService: eventService,
	}
}

func getTabTagFromQuery(c *gin.Context) string {
	return models.NormalizeTabTag(c.DefaultQuery("tab_tag", models.TabTagFutures))
}

func (h *AIResponseHandler) StreamAIResponseEvents(c *gin.Context) {
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, "streaming is not supported"))
		return
	}

	tabTag := getTabTagFromQuery(c)
	eventCh, unsubscribe := h.eventService.Subscribe(tabTag)
	defer unsubscribe()

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Status(http.StatusOK)
	flusher.Flush()

	keepAlive := time.NewTicker(20 * time.Second)
	defer keepAlive.Stop()

	for {
		select {
		case <-c.Request.Context().Done():
			return
		case event, ok := <-eventCh:
			if !ok {
				return
			}

			payload, err := json.Marshal(event)
			if err != nil {
				continue
			}

			_, _ = fmt.Fprintf(c.Writer, "event: ai_response_updated\n")
			_, _ = fmt.Fprintf(c.Writer, "data: %s\n\n", payload)
			flusher.Flush()
		case <-keepAlive.C:
			_, _ = fmt.Fprint(c.Writer, ": keep-alive\n\n")
			flusher.Flush()
		}
	}
}

func (h *AIResponseHandler) GetAIResponseByID(c *gin.Context) {
	id := c.Param("id")
	tabTag := getTabTagFromQuery(c)
	response, err := h.service.GetAIResponseByID(c.Request.Context(), tabTag, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(response))
}

func (h *AIResponseHandler) GetAIResponsesByBatchID(c *gin.Context) {
	batchID := c.Query("batch_id")
	tabTag := getTabTagFromQuery(c)
	responses, err := h.service.GetAIResponsesByBatchID(c.Request.Context(), tabTag, batchID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(responses))
}

func (h *AIResponseHandler) GetLatestBatch(c *gin.Context) {
	tabTag := getTabTagFromQuery(c)
	responses, err := h.service.GetLatestBatch(c.Request.Context(), tabTag)
	if err != nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse(404, err.Error()))
		return
	}

	completed := make([]models.AIResponse, 0, len(responses))
	for _, response := range responses {
		if response.Status == "completed" {
			completed = append(completed, response)
		}
	}

	c.JSON(http.StatusOK, models.SuccessResponse(completed))
}

func (h *AIResponseHandler) GetAllAIResponses(c *gin.Context) {
	tabTag := getTabTagFromQuery(c)
	responses, err := h.service.GetAllAIResponses(c.Request.Context(), tabTag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(responses))
}

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

	tabTag := models.NormalizeTabTag(req.TabTag)
	batchID, err := h.service.GenerateBatchAIResponses(c.Request.Context(), req.TemplateID, tabTag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(map[string]string{
		"batch_id": batchID,
		"tab_tag":  tabTag,
	}))
}
