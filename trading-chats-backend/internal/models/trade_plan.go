package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	TradePlanStatusPlanned   = "planned"
	TradePlanStatusActive    = "active"
	TradePlanStatusClosed    = "closed"
	TradePlanStatusCancelled = "cancelled"
)

type TradePlan struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TenantID   string             `bson:"tenant_id" json:"tenant_id"`
	TabTag     string             `bson:"tab_tag" json:"tab_tag"`
	Name       string             `bson:"name" json:"name"`
	Symbol     string             `bson:"symbol" json:"symbol"`
	Strategy   string             `bson:"strategy" json:"strategy"`
	Direction  string             `bson:"direction" json:"direction"`
	EntryPrice float64            `bson:"entry_price" json:"entry_price"`
	TakeProfit float64            `bson:"take_profit" json:"take_profit"`
	StopLoss   float64            `bson:"stop_loss" json:"stop_loss"`
	OpenTime   string             `bson:"open_time" json:"open_time"`
	CloseTime  string             `bson:"close_time" json:"close_time"`
	Status     string             `bson:"status" json:"status"`
	Remark     string             `bson:"remark" json:"remark"`
	CreatedAt  interface{}        `bson:"created_at" json:"created_at"`
	UpdatedAt  interface{}        `bson:"updated_at" json:"updated_at"`
}
