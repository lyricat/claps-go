package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type Project struct {
	Id          int64           `json:"id,omitempty" gorm:"type:bigint;primary_key;not null"`
	Name        string          `json:"name,omitempty" gorm:"type:varchar(50);not null;unique_index:name_UNIQUE"`
	DisplayName string          `json:"display_name,omitempty" gorm:"type:varchar(50);default:null"`
	Description string          `json:"description,omitempty" gorm:"type:varchar(120);default:null"`
	AvatarUrl   string          `json:"avatar_url,omitempty" gorm:"type:varchar(100);default:null"`
	Donations   int64           `json:"donations,omitempty" gorm:"type:bigint;default:0"`
	Total       decimal.Decimal `json:"total,omitempty" gorm:"type:varchar(128);default:null"`
	CreatedAt   time.Time       `json:"createdAt,omitempty" gorm:"type:timestamp with time zone"`
	UpdatedAt   time.Time       `json:"updatedAt,omitempty" gorm:"type:timestamp with time zone"`
}

type ProjectTotal struct {
	Id        int64          `json:"id,omitempty" gorm:"type:bigint;primary_key;not null"`
	Donations int64          `json:"donations,omitempty" gorm:"type:bigint;default:0"`
	Total     decimal.Decimal `json:"total,omitempty" gorm:"type:varchar(128);default:null"`
}
