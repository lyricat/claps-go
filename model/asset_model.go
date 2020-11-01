package model

import (
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
)

/**
 * @Description:注册自动迁移函数
 */
func init() {
	RegisterMigrateHandler(func(db *gorm.DB) error {

		if err := db.AutoMigrate(&Asset{}).Error; err != nil {
			return err
		}
		return nil
	})
}

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

var (
	ASSET      *Asset
	ASSETTOUSD *AssetIdToPriceUsd
)

/**
 * @Description:插入或定期更新asset表
 * @receiver asset
 * @param assetData
 * @return err
 */
func (asset *Asset) UpdateAsset(assetData *Asset) (err error) {
	err = db.Save(assetData).Error
	return
}

/**
 * @Description:通过assetId从数据库获取对应asset信息
 * @receiver asset
 * @param assetId
 * @return assetData
 * @return err
 */
func (asset *Asset) GetAssetById(assetId string) (assetData *Asset, err error) {
	assetData = &Asset{}
	err = db.Where("asset_id=?", assetId).Find(assetData).Error
	return
}

/**
 * @Description:获取数据库所有的asset信息
 * @receiver asset
 * @return assets
 * @return err
 */
func (asset *Asset) ListAssetsAllByDB() (assets *[]Asset, err error) {
	assets = &[]Asset{}
	err = db.Find(assets).Error
	return
}

/**
 * @Description:获取对应asset对应的美元汇率
 * @receiver assetIdToPriceUsd
 * @param assetId
 * @return priceUsd
 * @return err
 */
func (assetIdToPriceUsd *AssetIdToPriceUsd) GetPriceUsdByAssetId(assetId string) (priceUsd *AssetIdToPriceUsd, err error) {
	priceUsd = &AssetIdToPriceUsd{}
	err = db.Debug().Table("asset").Select("asset_id,price_usd").Where("asset_id=?", assetId).Scan(priceUsd).Error
	return
}
