package dao

import (
	"claps-test/model"
	"github.com/jinzhu/gorm"
)

func init() {
	RegisterMigrateHandler(func(db *gorm.DB) error {

		if err := db.AutoMigrate(&model.Repository{}).Error; err != nil {
			return err
		}
		return nil
	})
}


//根据project获取所有的仓库信息
func ListRepositoriesByProjectId(projectId int64) (repositories *[]model.Repository, err error) {
	repositories = &[]model.Repository{}
	err = db.Debug().Table("repository").Where("project_id=?", projectId).Scan(repositories).Error
	return
}
