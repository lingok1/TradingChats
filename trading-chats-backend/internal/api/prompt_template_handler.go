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
	return &PromptTemplateHandler{
		service: service,
	}
}

// CreatePromptTemplate 创建提示词模版
// @Summary 创建提示词模版
// @Description 创建新的提示词模版
// @Tags 提示词模版
// @Accept json
// @Produce json
// @Param template body models.PromptTemplate true "提示词模版信息"
// @Success 201 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /prompt-templates [post]
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

	c.JSON(http.StatusCreated, models.SuccessResponse(template))
}

// GetPromptTemplateByID 根据ID获取提示词模版
// @Summary 根据ID获取提示词模版
// @Description 根据ID获取提示词模版详情
// @Tags 提示词模版
// @Produce json
// @Param id path string true "提示词模版ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 404 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /prompt-templates/{id} [get]
func (h *PromptTemplateHandler) GetPromptTemplateByID(c *gin.Context) {
	id := c.Param("id")
	template, err := h.service.GetPromptTemplateByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(template))
}

// GetAllPromptTemplates 获取所有提示词模版
// @Summary 获取所有提示词模版
// @Description 获取所有提示词模版列表
// @Tags 提示词模版
// @Produce json
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /prompt-templates [get]
func (h *PromptTemplateHandler) GetAllPromptTemplates(c *gin.Context) {
	templates, err := h.service.GetAllPromptTemplates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(templates))
}

// GetPromptTemplatesByTag 根据标签获取提示词模版
// @Summary 根据标签获取提示词模版
// @Description 根据标签获取提示词模版列表
// @Tags 提示词模版
// @Produce json
// @Param tag query string true "标签"
// @Success 200 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /prompt-templates/tag [get]
func (h *PromptTemplateHandler) GetPromptTemplatesByTag(c *gin.Context) {
	tag := c.Query("tag")
	templates, err := h.service.GetPromptTemplatesByTag(c.Request.Context(), tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(templates))
}

// UpdatePromptTemplate 更新提示词模版
// @Summary 更新提示词模版
// @Description 更新提示词模版信息
// @Tags 提示词模版
// @Accept json
// @Produce json
// @Param id path string true "提示词模版ID"
// @Param template body models.PromptTemplate true "提示词模版信息"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /prompt-templates/{id} [put]
func (h *PromptTemplateHandler) UpdatePromptTemplate(c *gin.Context) {
	id := c.Param("id")
	var template models.PromptTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse(400, err.Error()))
		return
	}

	template.ID, _ = models.ParseObjectID(id)
	if err := h.service.UpdatePromptTemplate(c.Request.Context(), &template); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(template))
}

// DeletePromptTemplate 删除提示词模版
// @Summary 删除提示词模版
// @Description 删除提示词模版
// @Tags 提示词模版
// @Produce json
// @Param id path string true "提示词模版ID"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /prompt-templates/{id} [delete]
func (h *PromptTemplateHandler) DeletePromptTemplate(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeletePromptTemplate(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(map[string]string{"message": "Prompt template deleted successfully"}))
}

// GeneratePrompt 动态生成提示词
// @Summary 动态生成提示词
// @Description 根据提示词模版和参数动态生成提示词
// @Tags 提示词模版
// @Accept json
// @Produce json
// @Param request body models.GenerateAIRequest true "生成提示词请求"
// @Success 200 {object} models.Response
// @Failure 400 {object} models.Response
// @Failure 500 {object} models.Response
// @Router /prompt-templates/generate [post]
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

	prompt, err := h.service.GeneratePrompt(c.Request.Context(), req.TemplateID, req.Param1, req.Param2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.ErrorResponse(500, err.Error()))
		return
	}

	c.JSON(http.StatusOK, models.SuccessResponse(map[string]string{"prompt": prompt}))
}
