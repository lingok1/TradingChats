package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"trading-chats-backend/internal/models"
)

// PromptTemplateService 定义提示词模版服务接口
type PromptTemplateService interface {
	CreatePromptTemplate(ctx context.Context, template *models.PromptTemplate) error
	GetPromptTemplateByID(ctx context.Context, id string) (*models.PromptTemplate, error)
	GetAllPromptTemplates(ctx context.Context) ([]models.PromptTemplate, error)
	GetPromptTemplatesByTag(ctx context.Context, tag string) ([]models.PromptTemplate, error)
	UpdatePromptTemplate(ctx context.Context, template *models.PromptTemplate) error
	DeletePromptTemplate(ctx context.Context, id string) error
	GeneratePrompt(ctx context.Context, templateID string, param1, param2 string) (string, error)
}

// MockPromptTemplateService 模拟提示词模版服务
type MockPromptTemplateService struct {
	templates map[string]*models.PromptTemplate
}

func NewMockPromptTemplateService() *MockPromptTemplateService {
	return &MockPromptTemplateService{
		templates: make(map[string]*models.PromptTemplate),
	}
}

func (m *MockPromptTemplateService) CreatePromptTemplate(ctx context.Context, template *models.PromptTemplate) error {
	template.ID = primitive.NewObjectID()
	m.templates[template.ID.Hex()] = template
	return nil
}

func (m *MockPromptTemplateService) GetPromptTemplateByID(ctx context.Context, id string) (*models.PromptTemplate, error) {
	if template, ok := m.templates[id]; ok {
		return template, nil
	}
	return nil, fmt.Errorf("template not found")
}

func (m *MockPromptTemplateService) GetAllPromptTemplates(ctx context.Context) ([]models.PromptTemplate, error) {
	var templates []models.PromptTemplate
	for _, template := range m.templates {
		templates = append(templates, *template)
	}
	return templates, nil
}

func (m *MockPromptTemplateService) GetPromptTemplatesByTag(ctx context.Context, tag string) ([]models.PromptTemplate, error) {
	var templates []models.PromptTemplate
	for _, template := range m.templates {
		for _, t := range template.Tags {
			if t == tag {
				templates = append(templates, *template)
				break
			}
		}
	}
	return templates, nil
}

func (m *MockPromptTemplateService) UpdatePromptTemplate(ctx context.Context, template *models.PromptTemplate) error {
	m.templates[template.ID.Hex()] = template
	return nil
}

func (m *MockPromptTemplateService) DeletePromptTemplate(ctx context.Context, id string) error {
	delete(m.templates, id)
	return nil
}

func (m *MockPromptTemplateService) GeneratePrompt(ctx context.Context, templateID string, param1, param2 string) (string, error) {
	return "Generated prompt with params: " + param1 + " " + param2, nil
}

// 重新定义PromptTemplateHandler以使用接口
type PromptTemplateHandler struct {
	service PromptTemplateService
}

func NewPromptTemplateHandler(service PromptTemplateService) *PromptTemplateHandler {
	return &PromptTemplateHandler{
		service: service,
	}
}

// 复制原始处理器方法
func (h *PromptTemplateHandler) CreatePromptTemplate(c *gin.Context) {
	var template models.PromptTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreatePromptTemplate(c.Request.Context(), &template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, template)
}

func (h *PromptTemplateHandler) GetPromptTemplateByID(c *gin.Context) {
	id := c.Param("id")
	template, err := h.service.GetPromptTemplateByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

func (h *PromptTemplateHandler) GetAllPromptTemplates(c *gin.Context) {
	templates, err := h.service.GetAllPromptTemplates(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, templates)
}

func (h *PromptTemplateHandler) GetPromptTemplatesByTag(c *gin.Context) {
	tag := c.Query("tag")
	templates, err := h.service.GetPromptTemplatesByTag(c.Request.Context(), tag)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, templates)
}

func (h *PromptTemplateHandler) UpdatePromptTemplate(c *gin.Context) {
	id := c.Param("id")
	var template models.PromptTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	template.ID, _ = primitive.ObjectIDFromHex(id)
	if err := h.service.UpdatePromptTemplate(c.Request.Context(), &template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, template)
}

func (h *PromptTemplateHandler) DeletePromptTemplate(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeletePromptTemplate(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Prompt template deleted successfully"})
}

func (h *PromptTemplateHandler) GeneratePrompt(c *gin.Context) {
	var req map[string]string
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	templateID := req["template_id"]
	param1 := req["param1"]
	param2 := req["param2"]

	prompt, err := h.service.GeneratePrompt(c.Request.Context(), templateID, param1, param2)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"prompt": prompt})
}

func TestPromptTemplateAPI(t *testing.T) {
	// 创建模拟服务
	mockService := NewMockPromptTemplateService()
	handler := NewPromptTemplateHandler(mockService)

	// 创建Gin引擎
	r := gin.Default()

	// 注册路由
	promptTemplates := r.Group("/api/prompt-templates")
	{
		promptTemplates.POST("", handler.CreatePromptTemplate)
		promptTemplates.GET("", handler.GetAllPromptTemplates)
		promptTemplates.GET("/tag", handler.GetPromptTemplatesByTag)
		promptTemplates.GET("/:id", handler.GetPromptTemplateByID)
		promptTemplates.PUT("/:id", handler.UpdatePromptTemplate)
		promptTemplates.DELETE("/:id", handler.DeletePromptTemplate)
		promptTemplates.POST("/generate", handler.GeneratePrompt)
	}

	// 测试创建提示词模版
	t.Run("CreatePromptTemplate", func(t *testing.T) {
		template := models.PromptTemplate{
			Name:    "Test Template",
			Content: "Test content",
			Tags:    []string{"期货", "股票"},
		}

		data, _ := json.Marshal(template)
		req, _ := http.NewRequest("POST", "/api/prompt-templates", bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
		}

		var response models.PromptTemplate
		json.Unmarshal(w.Body.Bytes(), &response)

		if response.Name != template.Name {
			t.Errorf("Expected name %s, got %s", template.Name, response.Name)
		}
	})

	// 测试获取所有提示词模版
	t.Run("GetAllPromptTemplates", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/prompt-templates", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		var response []models.PromptTemplate
		json.Unmarshal(w.Body.Bytes(), &response)

		if len(response) == 0 {
			t.Error("Expected at least one template, got none")
		}
	})

	// 测试根据标签获取提示词模版
	t.Run("GetPromptTemplatesByTag", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/prompt-templates/tag?tag=期货", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		var response []models.PromptTemplate
		json.Unmarshal(w.Body.Bytes(), &response)

		if len(response) == 0 {
			t.Error("Expected at least one template with tag '期货', got none")
		}
	})

	// 测试动态生成提示词
	t.Run("GeneratePrompt", func(t *testing.T) {
		reqData := map[string]string{
			"template_id": "test-id",
			"param1":      "param1",
			"param2":      "param2",
		}

		data, _ := json.Marshal(reqData)
		req, _ := http.NewRequest("POST", "/api/prompt-templates/generate", bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["prompt"] == "" {
			t.Error("Expected prompt to be generated, got empty")
		}
	})
}
