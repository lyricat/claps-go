package dao

import (
	"claps-test/model"
	"github.com/jinzhu/gorm"
)

func init() {
	RegisterMigrateHandler(func(db *gorm.DB) error {

		if err := db.AutoMigrate(&model.Asset{}).Error; err != nil {
			return err
		}
		return nil
	})
}

func UpdateAsset(asset *model.Asset) (err error) {
	err = db.Save(asset).Error
	return
}

func GetAssetById(assetId string) (asset *model.Asset, err error) {
	asset = &model.Asset{}
	err = db.Where("asset_id=?", assetId).Find(asset).Error
	return
}

func ListAssetsAllByDB() (assets *[]model.Asset, err error) {
	assets = &[]model.Asset{}
	err = db.Find(assets).Error
	return
}

func GetPriceUsdByAssetId(assetId string) (priceUsd *model.AssetIdToPriceUsd, err error) {
	priceUsd = &model.AssetIdToPriceUsd{}
	err = db.Debug().Table("asset").Select("asset_id,price_usd").Where("asset_id=?", assetId).Scan(priceUsd).Error
	return
}
