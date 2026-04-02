package models

import "time"

const GlobalSystemConfigID = "global_config"

type SystemConfig struct {
	ID          string            `bson:"_id,omitempty" json:"id"`
	SystemTitle string            `bson:"system_title" json:"system_title"`
	SystemLogo  string            `bson:"system_logo" json:"system_logo"`
	Param1      string            `bson:"param1" json:"param1"`
	Param2      string            `bson:"param2" json:"param2"`
	Parameters  map[string]string `bson:"parameters" json:"parameters"`
	UpdatedAt   time.Time         `bson:"updated_at" json:"updated_at"`
}

type SaveSystemBasicConfigRequest struct {
	SystemTitle string `json:"system_title"`
	SystemLogo  string `json:"system_logo"`
}

type SaveSystemParametersRequest struct {
	Parameters map[string]string `json:"parameters"`
}

type SaveSystemRuntimeConfigRequest struct {
	Param1 string `json:"param1"`
	Param2 string `json:"param2"`
}
