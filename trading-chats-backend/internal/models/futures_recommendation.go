package models

import "time"

type FuturesRecommendation struct {
	ID           string               `bson:"_id,omitempty" json:"id"`
	BatchID      string               `bson:"batch_id" json:"batch_id"` // 引用的 ai_responses batch_id
	TabTag       string               `bson:"tab_tag" json:"tab_tag"`   // futures | options
	Items        []RecommendationItem `bson:"items" json:"items"`
	RawResponse  string               `bson:"raw_response" json:"raw_response"`
	ModelName    string               `bson:"model_name" json:"model_name"`
	ModelAPIName string               `bson:"model_api_name" json:"model_api_name"`
	CreatedAt    time.Time            `bson:"created_at" json:"created_at"`
}

type RecommendationItem struct {
	Symbol     string `bson:"symbol" json:"symbol"`
	Direction  string `bson:"direction" json:"direction"`
	EntryRange string `bson:"entry_range" json:"entry_range"`
	TakeProfit string `bson:"take_profit" json:"take_profit"`
	StopLoss   string `bson:"stop_loss" json:"stop_loss"`
	Reason     string `bson:"reason" json:"reason"`
}
