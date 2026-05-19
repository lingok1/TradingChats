package models

import "time"

type TenantConfig struct {
	ID         string            `bson:"_id" json:"id"`
	Parameters map[string]string `bson:"parameters" json:"parameters"`
	MenuConfig TenantMenuConfig  `bson:"menu_config" json:"menu_config"`
	UpdatedAt  time.Time         `bson:"updated_at" json:"updated_at"`
}

type TenantMenuConfig struct {
	VisibleTabs     []string `bson:"visible_tabs" json:"visible_tabs"`
	VisibleSettings []string `bson:"visible_settings" json:"visible_settings"`
}

var DefaultVisibleTabs = []string{"futures", "options", "stock", "plan", "about"}
var DefaultVisibleSettings = []string{"schedules", "models", "templates", "parameters"}
