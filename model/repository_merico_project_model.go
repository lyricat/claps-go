package model

import "github.com/jinzhu/gorm"

/**
 * @Description:注册自动迁移函数
 */
func init() {
	RegisterMigrateHandler(func(db *gorm.DB) error {

		if err := db.AutoMigrate(&RepositoryToMericoProject{}).Error; err != nil {
			return err
		}
		return nil
	})
}

type RepositoryToMericoProject struct {
	RepositoryId    int64  `gorm:"type:bigint;primary_key;not null"`
	MericoProjectId string `gorm:"type:varchar(50);primary_key;not null"`
}
