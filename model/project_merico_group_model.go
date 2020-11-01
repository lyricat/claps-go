package model

import "github.com/jinzhu/gorm"

/**
 * @Description:注册自动迁移函数
 */
func init() {
	RegisterMigrateHandler(func(db *gorm.DB) error {

		if err := db.AutoMigrate(&ProjectToMericoGroup{}).Error; err != nil {
			return err
		}
		return nil
	})
}

type ProjectToMericoGroup struct {
	ProjectId     int64  `gorm:"type:bigint;primary_key;not null"`
	MericoGroupId string `gorm:"type:varchar(50);primary_key;not null"`
}

var (
	PROJECTTOMERICOGROUP *ProjectToMericoGroup
)

func (ptmg *ProjectToMericoGroup) GetByProjectId(projectId int64) (ptm *ProjectToMericoGroup, err error) {
	ptm = &ProjectToMericoGroup{}
	err = db.Debug().Where("project_id=?", projectId).Find(ptm).Error
	return
}
