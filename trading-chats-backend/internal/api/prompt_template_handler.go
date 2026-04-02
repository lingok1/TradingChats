package api

import (
	"net/http"

	"trading-chats-backend/internal/models"
	"trading-chats-backend/internal/service"

	"github.com/gin-gonic/gin"
)

type PromptTemplateHandler struct {
	service *service.PromptTemplateService
}

func NewPromptTemplateHandler(service *service.PromptTemplateService) *PromptTemplateHandler {
	return &PromptTemplateHandler{service: service}
}

// CreatePromptTemplate 创建提示词模板
// @Summary 创建提示词模板
// @Description 创建新的提示词模板
// @Tags 提示词模板
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.PromptTemplate true "提示词模板"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/prompt-templates [post]
func (h *PromptTemplateHandler) CreatePromptTemplate(c *gin.Context) {
	var template models.PromptTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}

	if err := h.service.CreatePromptTemplate(c.Request.Context(), &template); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(template))
}

// GetAllPromptTemplates 获取提示词模板列表
// @Summary 获取提示词模板列表
// @Description 获取当前可见的提示词模板列表
// @Tags 提示词模板
// @Produce json
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/prompt-templates [get]
func (h *PromptTemplateHandler) GetAllPromptTemplates(c *gin.Context) {
	templates, err := h.service.GetAllPromptTemplates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(templates))
}

// GetPromptTemplatesByTag 按标签获取提示词模板
// @Summary 按标签获取提示词模板
// @Description 根据 tag 查询提示词模板
// @Tags 提示词模板
// @Produce json
// @Param tag query string true "标签"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/prompt-templates/tag [get]
func (h *PromptTemplateHandler) GetPromptTemplatesByTag(c *gin.Context) {
	tag := c.Query("tag")
	templates, err := h.service.GetPromptTemplatesByTag(c.Request.Context(), tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(templates))
}

// GetPromptTemplateByID 获取提示词模板详情
// @Summary 获取提示词模板详情
// @Description 根据 ID 获取提示词模板
// @Tags 提示词模板
// @Produce json
// @Param id path string true "提示词模板ID"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/prompt-templates/{id} [get]
func (h *PromptTemplateHandler) GetPromptTemplateByID(c *gin.Context) {
	id := c.Param("id")
	template, err := h.service.GetPromptTemplateByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(template))
}

// UpdatePromptTemplate 更新提示词模板
// @Summary 更新提示词模板
// @Description 根据 ID 更新提示词模板
// @Tags 提示词模板
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "提示词模板ID"
// @Param body body models.PromptTemplate true "提示词模板"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/prompt-templates/{id} [put]
func (h *PromptTemplateHandler) UpdatePromptTemplate(c *gin.Context) {
	id := c.Param("id")
	var template models.PromptTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}

	parsedID, err := models.ParseObjectID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}
	template.ID = parsedID

	if err := h.service.UpdatePromptTemplate(c.Request.Context(), &template); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(template))
}

// DeletePromptTemplate 删除提示词模板
// @Summary 删除提示词模板
// @Description 根据 ID 删除提示词模板
// @Tags 提示词模板
// @Produce json
// @Security BearerAuth
// @Param id path string true "提示词模板ID"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/prompt-templates/{id} [delete]
func (h *PromptTemplateHandler) DeletePromptTemplate(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeletePromptTemplate(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(map[string]string{"message": "Prompt template deleted successfully"}))
}

// GeneratePrompt 生成 Prompt
// @Summary 生成 Prompt
// @Description 基于模板生成 Prompt 内容
// @Tags 提示词模板
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param body body models.GenerateAIRequest true "生成请求"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /api/prompt-templates/generate [post]
func (h *PromptTemplateHandler) GeneratePrompt(c *gin.Context) {
	var req models.GenerateAIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}

	if req.TemplateID == "" {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, "template_id is required"))
		return
	}

	prompt, err := h.service.GeneratePrompt(c.Request.Context(), req.TemplateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(map[string]string{"prompt": prompt}))
}
