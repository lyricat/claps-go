package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type MemberWallet struct {
	ProjectId int64 `json:"project_id,omitempty" gorm:"type:bigint;primary_key;not null"`
	//user表的Id
	UserId    int64 `json:"user_id,omitempty" gorm:"type:bigint;primary_key;not null"`
	BotId     string `json:"bot_id,omitempty" gorm:"type:varchar(50);primary_key;not null"`
	AssetId   string `json:"asset_id,omitempty" gorm:"type:varchar(50);primary_key;not null"`
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"type:timestamp with time zone"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"type:timestamp with time zone"`
	Total     decimal.Decimal `json:"total,omitempty" gorm:"type:varchar(128);default:null"`
	Balance   decimal.Decimal `json:"balance,omitempty" gorm:"type:varchar(128);default:null"`
}

type MemberWalletDto struct {
	ProjectId int64          `json:"project_id,omitempty" gorm:"type:bigint;primary_key;not null"`
	UserId    int64          `json:"user_id,omitempty" gorm:"type:bigint;primary_key;not null"`
	BotId     string          `json:"bot_id,omitempty" gorm:"type:varchar(50);primary_key;not null"`
	AssetId   string          `json:"asset_id,omitempty" gorm:"type:varchar(50);primary_key;not null"`
	Total     decimal.Decimal `json:"total,omitempty" gorm:"type:varchar(128);default:null"`
	Balance   decimal.Decimal `json:"balance,omitempty" gorm:"type:varchar(128);default:null"`
}
