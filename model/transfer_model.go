package model

import (
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	log "github.com/sirupsen/logrus"
	"time"
)

/**
 * @Description:注册自动迁移函数
 */
func init() {
	RegisterMigrateHandler(func(db *gorm.DB) error {

		if err := db.AutoMigrate(&Transfer{}).Error; err != nil {
			return err
		}
		return nil
	})
}

type Transfer struct {
	//机器人ID
	BotId      string          `json:"bot_id,omitempty" gorm:"type:varchar(50);not null"`
	SnapshotId string          `json:"snapshot_id,omitempty" gorm:"type:varchar(50);default null"`
	MixinId    string          `json:"mixin_id,omitempty" gorm:"type:varchar(50);not null"`
	TraceId    string          `json:"trace_id,omitempty" gorm:"type:varchar(50);not null;primary_key"`
	AssetId    string          `json:"asset_id,omitempty" gorm:"type:varchar(50);not null"`
	Amount     decimal.Decimal `json:"amount,omitempty" gorm:"type:varchar(128);not null"`
	Memo       string          `json:"memo,omitempty" gorm:"type:varchar(120);not null"`
	Status     string          `json:"status,omitempty" gorm:"type:char;not null;index:status_INDEX"`
	//0是用户点击提现后(提现操作未完成) 1机器人完成提现
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"type:timestamp with time zone"`
}

const (
	UNFINISHED = "0"
	FINISHED   = "1"
)

var TRANSFER *Transfer

/**
 * @Description: 首次插入或者更新transfer的status值
 * @receiver transfer
 * @param transferData
 * @return err
 */
func (transfer *Transfer) InsertOrUpdateTransfer(transferData *Transfer) (err error) {
	err = db.Debug().Save(transferData).Error
	return
}

/**
 * @Description: 通过mixinId获取对应transfers的值,暂时废弃
 * @receiver transfer
 * @param mixinId
 * @return transfers
 * @return err
 */
func (transfer *Transfer) ListTransferByMixinId(mixinId string) (transfers *[]Transfer, err error) {
	transfers = &[]Transfer{}
	err = db.Debug().Where("mixin_id = ?", mixinId).Find(transfers).Error
	return
}

/**
 * @Description: 通过mixinId获取对应transfers的值,暂时废弃
 * @receiver transfer
 * @param mixinId
 * @return transfers
 * @return err
 */
func (transfer *Transfer) getTransfersNumbersByMixinId(mixinId string) (number int, err error) {

	err = db.Debug().Table("transfer").Where("mixin_id = ?", mixinId).Count(&number).Error
	return
}

/**
 * @Description: 通过mixinId和query值获取对应transfers的值
 * @receiver transfer
 * @param mixinId
 * @param q
 * @return transfers
 * @return number
 * @return err
 */
func (transfer *Transfer) ListTransfersByMixinIdAndQuery(mixinId string, q *PaginationQ) (transfers *[]Transfer, number int, err error) {

	transfers = &[]Transfer{}
	number, err = transfer.getTransfersNumbersByMixinId(mixinId)
	if err != nil {
		return
	}

	tx := db.Debug().Table("transfer")
	if q.Limit <= 0 {
		q.Limit = 20
	}

	if q.Offset <= 0 {
		q.Offset = 0
	}
	err = tx.Where("mixin_id = ?",
		mixinId).Order("created_at desc").Limit(q.Limit).Offset(q.Offset).Find(transfers).Error

	return
}

/**
 * @Description: status only '0' or '1'
 * @receiver transfer
 * @param status
 * @return transfers
 * @return err
 */
func (transfer *Transfer) ListTransfersByStatus(status string) (transfers *[]Transfer, err error) {
	transfers = &[]Transfer{}
	err = db.Where("status=?", status).Find(transfers).Error
	return
}

/**
 * @Description: 统计该mixinId提现记录处于未完成状态的数量
 * @receiver transfer
 * @param mixinId
 * @return count
 * @return err
 */
func (transfer *Transfer) CountUnfinishedTransfer(mixinId string) (count int, err error) {
	err = db.Debug().Table("transfer").Where("status = ? AND mixin_id = ?", UNFINISHED, mixinId).Count(&count).Error
	log.Debug(count)
	return
}
