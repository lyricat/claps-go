package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type Wallet struct {
	BotId     string          `json:"bot_id,omitempty" gorm:"type:varchar(50);primary_key;not null"`
	AssetId   string          `json:"asset_id,omitempty" gorm:"type:varchar(50);primary_key;not null"`
	ProjectId int64           `json:"project_id,omitempty" gorm:"type:bigint;not null"`
	Total     decimal.Decimal `json:"total,omitempty" gorm:"type:varchar(128);default:null"`
	CreatedAt time.Time		  `json:"created_at,omitempty" gorm:"type:timestamp with time zone"`
	UpdatedAt time.Time       `json:"updated_at,omitempty" gorm:"type:timestamp with time zone"`
	SyncedAt  time.Time       `json:"synced_at,omitempty" gorm:"type:timestamp with time zone"`
}

type WalletTotal struct {
	BotId   string          `json:"bot_id,omitempty" gorm:"type:varchar(50);primary_key;not null"`
	AssetId string          `json:"asset_id,omitempty" gorm:"type:varchar(50);primary_key;not null"`
	Total   decimal.Decimal `json:"total,omitempty" gorm:"type:varchar(128);not null;default:null"`
}
