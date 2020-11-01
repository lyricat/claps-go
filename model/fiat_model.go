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

		if err := db.AutoMigrate(&Fiat{}).Error; err != nil {
			return err
		}
		return nil
	})
}

type Fiat struct {
	Code string          `json:"code,omitempty" gorm:"type:varchar(25);primary_key;not null"`
	Rate decimal.Decimal `json:"rate,omitempty" gorm:"type:varchar(128);not null"`
}

var FIAT *Fiat

/**
 * @Description: 更新汇率信息
 * @receiver fiat
 * @param fiatData
 * @return err
 */
func (fiat *Fiat) UpdateFiat(fiatData *Fiat) (err error) {
	err = db.Debug().Save(fiatData).Error
	return
}

/**
 * @Description: 从数据库中读取所有汇率信息
 * @receiver fiat
 * @return fiats
 * @return err
 */
func (fiat *Fiat) ListFiatsAllByDB() (fiats *[]Fiat, err error) {
	fiats = &[]Fiat{}
	err = db.Find(fiats).Error
	return
}

/**
 * @Description: 通过对应Code获取汇率信息
 * @receiver fiat
 * @param code
 * @return fiatData
 * @return err
 */
func (fiat *Fiat) GetFiatByCode(code string) (fiatData *Fiat, err error) {
	fiatData = &Fiat{}
	err = db.Select("rate").Where("code = ? ", code).Find(&fiat).Error
	return
}
