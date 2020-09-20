package model

import "github.com/shopspring/decimal"

type Asset struct {
	AssetId  string          `json:"asset_id,omitempty" gorm:"type:varchar(36);primary_key;not null;"`
	Symbol   string          `json:"symbol,omitempty" gorm:"type:varchar(512);not null"`
	Name     string          `json:"name,omitempty" gorm:"type:varchar(512);not null"`
	IconUrl  string          `json:"icon_url,omitempty" gorm:"type:varchar(1024);not null"`
	PriceBtc decimal.Decimal `json:"price_btc,omitempty" gorm:"type:varchar(128);not null"`
	PriceUsd decimal.Decimal `json:"price_usd,omitempty" gorm:"type:varchar(128);not null"`
}

type AssetIdToPriceUsd struct {
	AssetId  string          `json:"asset_id,omitempty" gorm:"type:varchar(36);primary_key;not null;"`
	PriceUsd decimal.Decimal `json:"price_usd,omitempty" gorm:"type:varchar(128);not null"`
}
