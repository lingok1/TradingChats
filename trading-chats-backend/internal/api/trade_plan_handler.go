package api

import (
	"net/http"

	"trading-chats-backend/internal/models"
	"trading-chats-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type TradePlanHandler struct {
	service *service.TradePlanService
}

func NewTradePlanHandler(service *service.TradePlanService) *TradePlanHandler {
	return &TradePlanHandler{service: service}
}

func getTradePlanTabTag(c *gin.Context) string {
	return c.DefaultQuery("tab_tag", models.TabTagFutures)
}

func (h *TradePlanHandler) CreateTradePlan(c *gin.Context) {
	var plan models.TradePlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}

	authCtx := MustGetAuthContext(c)
	if authCtx != nil {
		plan.TenantID = authCtx.TenantID
	}

	if err := h.service.CreateTradePlan(c.Request.Context(), &plan); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(plan))
}

func (h *TradePlanHandler) GetTradePlans(c *gin.Context) {
	tabTag := getTradePlanTabTag(c)
	plans, err := h.service.GetTradePlans(c.Request.Context(), tabTag)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(plans))
}

func (h *TradePlanHandler) GetTradePlanByID(c *gin.Context) {
	id := c.Param("id")
	plan, err := h.service.GetTradePlanByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(plan))
}

func (h *TradePlanHandler) UpdateTradePlan(c *gin.Context) {
	id := c.Param("id")

	var plan models.TradePlan
	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}

	authCtx := MustGetAuthContext(c)
	if authCtx != nil {
		plan.TenantID = authCtx.TenantID
	}

	objectID, err := models.ParseObjectID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}
	plan.ID = objectID

	if err := h.service.UpdateTradePlan(c.Request.Context(), &plan); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("updated successfully"))
}

func (h *TradePlanHandler) DeleteTradePlan(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteTradePlan(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse("deleted successfully"))
}
