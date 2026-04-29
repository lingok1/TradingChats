package api

import (
	"net/http"
	"os"
	"runtime"
	"time"
	"trading-chats-backend/internal/models"

	"github.com/gin-gonic/gin"
)

const backendName = "trading-chats-backend"

type SystemInfo struct {
	OS           string `json:"os"`
	Arch         string `json:"arch"`
	GoVersion    string `json:"go_version"`
	CPUCount     int    `json:"cpu_count"`
	Hostname     string `json:"hostname"`
	Compiler     string `json:"compiler"`
	NumGoroutine int    `json:"num_goroutine"`
}

type BackendInfo struct {
	Status      string     `json:"status"`
	Name        string     `json:"name"`
	CurrentTime string     `json:"current_time"`
	System      SystemInfo `json:"system"`
}

// BackendInfoHandler 后端服务信息
// @Summary 后端服务信息
// @Description 返回后端状态、名称、当前时间和系统信息
// @Tags 系统
// @Produce json
// @Success 200 {object} models.Response
// @Router / [get]
func BackendInfoHandler(c *gin.Context) {
	hostname, _ := os.Hostname()

	c.JSON(http.StatusOK, models.SuccessResponse(BackendInfo{
		Status:      "ok",
		Name:        backendName,
		CurrentTime: time.Now().Format(time.RFC3339),
		System: SystemInfo{
			OS:           runtime.GOOS,
			Arch:         runtime.GOARCH,
			GoVersion:    runtime.Version(),
			CPUCount:     runtime.NumCPU(),
			Hostname:     hostname,
			Compiler:     runtime.Compiler,
			NumGoroutine: runtime.NumGoroutine(),
		},
	}))
}

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
