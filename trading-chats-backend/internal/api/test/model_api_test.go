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

// ModelAPIService 定义模型与API配置服务接口
type ModelAPIService interface {
	CreateModelAPIConfig(ctx context.Context, config *models.ModelAPIConfig) error
	GetModelAPIConfigByID(ctx context.Context, id string) (*models.ModelAPIConfig, error)
	GetAllModelAPIConfigs(ctx context.Context) ([]models.ModelAPIConfig, error)
	GetModelAPIConfigsByProvider(ctx context.Context, provider string) ([]models.ModelAPIConfig, error)
	UpdateModelAPIConfig(ctx context.Context, config *models.ModelAPIConfig) error
	DeleteModelAPIConfig(ctx context.Context, id string) error
	TestModelConnectivity(ctx context.Context, configID string) error
}

// MockModelAPIService 模拟模型与API配置服务
type MockModelAPIService struct {
	configs map[string]*models.ModelAPIConfig
}

func NewMockModelAPIService() *MockModelAPIService {
	return &MockModelAPIService{
		configs: make(map[string]*models.ModelAPIConfig),
	}
}

func (m *MockModelAPIService) CreateModelAPIConfig(ctx context.Context, config *models.ModelAPIConfig) error {
	config.ID = primitive.NewObjectID()
	m.configs[config.ID.Hex()] = config
	return nil
}

func (m *MockModelAPIService) GetModelAPIConfigByID(ctx context.Context, id string) (*models.ModelAPIConfig, error) {
	if config, ok := m.configs[id]; ok {
		return config, nil
	}
	return nil, fmt.Errorf("config not found")
}

func (m *MockModelAPIService) GetAllModelAPIConfigs(ctx context.Context) ([]models.ModelAPIConfig, error) {
	var configs []models.ModelAPIConfig
	for _, config := range m.configs {
		configs = append(configs, *config)
	}
	return configs, nil
}

func (m *MockModelAPIService) GetModelAPIConfigsByProvider(ctx context.Context, provider string) ([]models.ModelAPIConfig, error) {
	var configs []models.ModelAPIConfig
	for _, config := range m.configs {
		if config.Provider == provider {
			configs = append(configs, *config)
		}
	}
	return configs, nil
}

func (m *MockModelAPIService) UpdateModelAPIConfig(ctx context.Context, config *models.ModelAPIConfig) error {
	m.configs[config.ID.Hex()] = config
	return nil
}

func (m *MockModelAPIService) DeleteModelAPIConfig(ctx context.Context, id string) error {
	delete(m.configs, id)
	return nil
}

func (m *MockModelAPIService) TestModelConnectivity(ctx context.Context, configID string) error {
	return nil
}

// 重新定义ModelAPIHandler以使用接口
type ModelAPIHandler struct {
	service ModelAPIService
}

func NewModelAPIHandler(service ModelAPIService) *ModelAPIHandler {
	return &ModelAPIHandler{
		service: service,
	}
}

// 复制原始处理器方法
func (h *ModelAPIHandler) CreateModelAPIConfig(c *gin.Context) {
	var config models.ModelAPIConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateModelAPIConfig(c.Request.Context(), &config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, config)
}

func (h *ModelAPIHandler) GetModelAPIConfigByID(c *gin.Context) {
	id := c.Param("id")
	config, err := h.service.GetModelAPIConfigByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

func (h *ModelAPIHandler) GetAllModelAPIConfigs(c *gin.Context) {
	configs, err := h.service.GetAllModelAPIConfigs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, configs)
}

func (h *ModelAPIHandler) GetModelAPIConfigsByProvider(c *gin.Context) {
	provider := c.Query("provider")
	configs, err := h.service.GetModelAPIConfigsByProvider(c.Request.Context(), provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, configs)
}

func (h *ModelAPIHandler) UpdateModelAPIConfig(c *gin.Context) {
	id := c.Param("id")
	var config models.ModelAPIConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.ID, _ = primitive.ObjectIDFromHex(id)
	if err := h.service.UpdateModelAPIConfig(c.Request.Context(), &config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, config)
}

func (h *ModelAPIHandler) DeleteModelAPIConfig(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteModelAPIConfig(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Model API config deleted successfully"})
}

func (h *ModelAPIHandler) TestModelConnectivity(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.TestModelConnectivity(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Connectivity test successful"})
}

func TestModelAPIAPI(t *testing.T) {
	// 创建模拟服务
	mockService := NewMockModelAPIService()
	handler := NewModelAPIHandler(mockService)

	// 创建Gin引擎
	r := gin.Default()

	// 注册路由
	modelAPIConfigs := r.Group("/api/model-api-configs")
	{
		modelAPIConfigs.POST("", handler.CreateModelAPIConfig)
		modelAPIConfigs.GET("", handler.GetAllModelAPIConfigs)
		modelAPIConfigs.GET("/provider", handler.GetModelAPIConfigsByProvider)
		modelAPIConfigs.GET("/:id", handler.GetModelAPIConfigByID)
		modelAPIConfigs.PUT("/:id", handler.UpdateModelAPIConfig)
		modelAPIConfigs.DELETE("/:id", handler.DeleteModelAPIConfig)
		modelAPIConfigs.POST("/:id/test", handler.TestModelConnectivity)
	}

	// 测试创建模型与API配置
	t.Run("CreateModelAPIConfig", func(t *testing.T) {
		config := models.ModelAPIConfig{
			Name:     "Test Config",
			APIURL:   "https://api.openai.com/v1/chat/completions",
			APIKey:   "test-key",
			Models:   []string{"gpt-3.5-turbo", "gpt-4"},
			Provider: "openai",
		}

		data, _ := json.Marshal(config)
		req, _ := http.NewRequest("POST", "/api/model-api-configs", bytes.NewBuffer(data))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status code %d, got %d", http.StatusCreated, w.Code)
		}

		var response models.ModelAPIConfig
		json.Unmarshal(w.Body.Bytes(), &response)

		if response.Name != config.Name {
			t.Errorf("Expected name %s, got %s", config.Name, response.Name)
		}
	})

	// 测试获取所有模型与API配置
	t.Run("GetAllModelAPIConfigs", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/model-api-configs", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		var response []models.ModelAPIConfig
		json.Unmarshal(w.Body.Bytes(), &response)

		if len(response) == 0 {
			t.Error("Expected at least one config, got none")
		}
	})

	// 测试根据提供商获取模型与API配置
	t.Run("GetModelAPIConfigsByProvider", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/model-api-configs/provider?provider=openai", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		var response []models.ModelAPIConfig
		json.Unmarshal(w.Body.Bytes(), &response)

		if len(response) == 0 {
			t.Error("Expected at least one config with provider 'openai', got none")
		}
	})

	// 测试测试模型连通性
	t.Run("TestModelConnectivity", func(t *testing.T) {
		req, _ := http.NewRequest("POST", "/api/model-api-configs/test-id/test", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		var response map[string]string
		json.Unmarshal(w.Body.Bytes(), &response)

		if response["message"] == "" {
			t.Error("Expected message to be returned, got empty")
		}
	})
}
