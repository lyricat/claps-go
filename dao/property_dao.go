package dao

import (
	"claps-test/model"
	"github.com/jinzhu/gorm"
)

func init() {
	RegisterMigrateHandler(func(db *gorm.DB) error {

		if err := db.AutoMigrate(&model.Property{}).Error; err != nil {
			return err
		}
		return nil
	})
}


func GetPropertyByKey(Key string) (property *model.Property, err error) {
	property = &model.Property{
		Key: Key,
	}
	err = db.First(property).Error
	return
}

func UpdateProperty(property *model.Property) (err error) {
	err = db.Save(property).Error
	return
}

func InsertProperty(property *model.Property) (err error) {
	err = db.Create(property).Error
	return
}
