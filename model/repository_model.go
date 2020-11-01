package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

/**
 * @Description:注册自动迁移函数
 */
func init() {
	RegisterMigrateHandler(func(db *gorm.DB) error {

		if err := db.AutoMigrate(&Repository{}).Error; err != nil {
			return err
		}
		return nil
	})
}

type Repository struct {
	Id          int64     `json:"id,omitempty" gorm:"type:bigserial;primary_key;auto_increment"`
	ProjectId   int64     `json:"project_id,omitempty" gorm:"type:bigint;index:repository_project_INDEX;not null"`
	Type        string    `json:"type,omitempty" gorm:"varchar(10);not null"`
	Name        string    `json:"name,omitempty" gorm:"type:varchar(50);not null"`
	Slug        string    `json:"slug,omitempty" gorm:"type:varchar(100);not null"`
	Description string    `json:"description,omitempty" gorm:"type:varchar(120);default:null"`
	CreatedAt   time.Time `json:"created_at,omitempty" gorm:"type:timestamp with time zone"`
	UpdatedAt   time.Time `json:"updated_at,omitempty" gorm:"type:timestamp with time zone"`
}

var REPOSITORY *Repository

/**
 * @Description: 根据project获取所有的仓库信息
 * @receiver repository
 * @param projectId
 * @return repositories
 * @return err
 */
func (repository *Repository) ListRepositoriesByProjectId(projectId int64) (repositories *[]Repository, err error) {
	repositories = &[]Repository{}
	err = db.Debug().Table("repository").Where("project_id=?", projectId).Scan(repositories).Error
	return
}
