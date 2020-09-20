package dao

import (
	"claps-test/model"
	"github.com/jinzhu/gorm"
)

func init() {
	RegisterMigrateHandler(func(db *gorm.DB) error {

		if err := db.AutoMigrate(&model.Wallet{}).Error; err != nil {
			return err
		}
		return nil
	})
}


func GetWalletTotalByBotIdAndAssetId(botId string, assetId string) (total *model.WalletTotal, err error) {
	total = &model.WalletTotal{}
	err = db.Debug().Table("wallet").Where("bot_id=? AND asset_id=?", botId, assetId).Find(total).Error
	return
}

func UpdateWalletTotal(walletTotal *model.WalletTotal) (err error) {
	err = db.Table("wallet").Save(walletTotal).Error
	return
}

//通过项目Id获取项目的Total?
func GetWalletTotalByProjectId(projectId uint32) (total *[]model.WalletTotal, err error) {
	total = &[]model.WalletTotal{}
	err = db.Debug().Table("wallet").Select("project_id,asset_id,total").Where("project_id=?", projectId).Scan(total).Error
	return
}
