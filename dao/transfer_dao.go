package dao

import (
	"claps-test/model"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

func init() {
	RegisterMigrateHandler(func(db *gorm.DB) error {

		if err := db.AutoMigrate(&model.Transfer{}).Error; err != nil {
			return err
		}
		return nil
	})
}


func InsertTransfer(transfer *model.Transfer) (err error) {
	err = db.Create(transfer).Error
	return
}

func UpdateTransferTraceId(transferMap *map[string]interface{}, traceId string) (err error) {
	err = db.Debug().Model(model.Transfer{}).Where("trace_id=?", traceId).Updates(*transferMap).Error
	return
}

func GetTransferByMixinId(mixinId string) (transfers *[]model.Transfer, err error) {
	transfers = &[]model.Transfer{}
	err = db.Debug().Where("mixin_id = ?", mixinId).Find(transfers).Error
	return
}

//status only '0' or '1' or '2'
func ListTransfersByStatus(status string) (transfer *[]model.Transfer, err error) {
	transfer = &[]model.Transfer{}
	err = db.Where("status=?", status).Find(transfer).Error
	return
}

func UpdateTransferStatusByUserIdAndAssetId(mixinId string, assetId string, status string) (err error) {
	err = db.Debug().Table("transfer").Where("mixin_id = ? AND asset_id= ?", mixinId, assetId).Update("status", status).Error
	return
}

func CountUnfinishedTransfer(mixinId string) (count int, err error) {
	err = db.Debug().Table("transfer").Where("mixin_id = ? AND status = ?", mixinId, model.UNFINISHED).Count(&count).Error
	log.Debug(count)
	return
}
