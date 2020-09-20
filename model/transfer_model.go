package model

import (
	"github.com/shopspring/decimal"
	"time"
)

const (
	UNFINISHED = "0"
	FINISHED   = "1"
)

type Transfer struct {
	//机器人ID
	BotId      string          `json:"bot_id,omitempty" gorm:"type:varchar(50);not null"`
	SnapshotId string          `json:"snapshot_id,omitempty" gorm:"type:varchar(50);default null"`
	MixinId    string          `json:"mixin_id,omitempty" gorm:"type:varchar(50);not null"`
	TraceId    string          `json:"trace_id,omitempty" gorm:"type:varchar(100);not null;primary_key"`
	AssetId    string          `json:"asset_id,omitempty" gorm:"type:varchar(50);not null"`
	Amount     decimal.Decimal `json:"amount,omitempty" gorm:"type:varchar(128);not null"`
	Memo       string          `json:"memo,omitempty" gorm:"type:varchar(120);not null"`
	Status     string          `json:"status,omitempty" gorm:"type:char;not null;index:status_INDEX"`
	//0是用户点击提现后(提现操作未完成) 1机器人完成提现
	CreatedAt time.Time        `json:"created_at,omitempty" gorm:"type:timestamp with time zone"`
}
