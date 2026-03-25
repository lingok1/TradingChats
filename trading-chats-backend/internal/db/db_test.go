package db

import (
	"testing"

	"trading-chats-backend/internal/config"
)

func TestConnect(t *testing.T) {
	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// 连接数据库
	if err := Connect(cfg); err != nil {
		t.Fatalf("Failed to connect to databases: %v", err)
	}

	// 断开连接
	if err := Disconnect(); err != nil {
		t.Fatalf("Failed to disconnect from databases: %v", err)
	}

	t.Log("Database connection test passed")
}
