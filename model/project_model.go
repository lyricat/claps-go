package model

import (
	"github.com/jinzhu/gorm"
	"github.com/shopspring/decimal"
	"time"
)

func init() {
	RegisterMigrateHandler(func(db *gorm.DB) error {

		if err := db.AutoMigrate(&Project{}).Error; err != nil {
			return err
		}
		return nil
	})
}

type Project struct {
	Id          int64           `json:"id,omitempty" gorm:"type:bigserial;auto_increment;primary_key"`
	Name        string          `json:"name,omitempty" gorm:"type:varchar(50);not null"`
	DisplayName string          `json:"display_name,omitempty" gorm:"type:varchar(50);default:null"`
	Description string          `json:"description,omitempty" gorm:"type:varchar(120);default:null"`
	AvatarUrl   string          `json:"avatar_url,omitempty" gorm:"type:varchar(100);default:null"`
	Donations   int64           `json:"donations,omitempty" gorm:"type:bigint;default:0"`
	Total       decimal.Decimal `json:"total,omitempty" gorm:"type:varchar(128);default:null"`
	CreatedAt   time.Time       `json:"createdAt,omitempty" gorm:"type:timestamp with time zone"`
	UpdatedAt   time.Time       `json:"updatedAt,omitempty" gorm:"type:timestamp with time zone"`
}

type ProjectTotal struct {
	Id        int64           `json:"id,omitempty" gorm:"type:bigserial;auto_increment;primary_key"`
	Donations int64           `json:"donations,omitempty" gorm:"type:bigint;default:0"`
	Total     decimal.Decimal `json:"total,omitempty" gorm:"type:varchar(128);default:null"`
}

type Badge struct {
	Code    string `form:"code" json:"code" binding:"required"`
	Color   string `form:"color" json:"color" binding:"required"`
	BgColor string `form:"bg_color" json:"bg_color" binding:"required"`
	Size    string `form:"size" json:"size" binding:"required"`
}

var (
	PROJECT      *Project
	PROJECTTOTAL *ProjectTotal
)

/**
 * @Description: 获取所有项目
 * @receiver proj
 * @return projects
 * @return err
 */
func (proj *Project) ListProjectsAll() (projects *[]Project, err error) {

	projects = &[]Project{}
	err = db.Debug().Find(projects).Error
	return
}

/**
 * @Description: 获取一共有多少条project记录
 * @receiver proj
 * @return number
 * @return err
 */
func (proj *Project) getProjectsNumbers() (number int, err error) {

	err = db.Debug().Table("project").Count(&number).Error
	return
}

/**
 * @Description: 通过query获取projects
 * @receiver proj
 * @param q
 * @return projects
 * @return number
 * @return err
 */
func (proj *Project) ListProjectsByQuery(q *PaginationQ) (projects *[]Project, number int, err error) {
	projects = &[]Project{}
	number, err = proj.getProjectsNumbers()
	if err != nil {
		return
	}

	tx := db.Debug().Table("project")
	if q.Limit <= 0 {
		q.Limit = 20
	}

	if q.Offset <= 0 {
		q.Offset = 0
	}

	if q.Q != "" {
		tx = tx.Where("name Like ?", "%"+q.Q+"%")
	}

	err = tx.Limit(q.Limit).Offset(q.Offset).Find(projects).Error
	return
}

/**
 * @Description: 通过项目id获取项目
 * @receiver proj
 * @param projectId
 * @return project
 * @return err
 */
func (proj *Project) GetProjectById(projectId int64) (project *Project, err error) {

	project = &Project{}
	err = db.Debug().Where("id=?", projectId).Find(project).Error
	return
}

/**
 * @Description: 根据userId获取所有项目
 * @receiver proj
 * @param userId
 * @return projects
 * @return err
 */
func (proj *Project) ListProjectsByUserId(userId int64) (projects *[]Project, err error) {
	projects = &[]Project{}
	err = db.Debug().Where("id IN(?)",
		db.Debug().Table("member").Select("project_id").Where("user_id=?", userId).SubQuery()).Find(projects).Error
	return
}

/**
 * @Description: 根据botId获取对应project的Total值和收到的捐赠笔数
 * @receiver projTotal
 * @param BotId
 * @return projectTotal
 * @return err
 */
func (projTotal *ProjectTotal) GetProjectTotalByBotId(BotId string) (projectTotal *ProjectTotal, err error) {
	projectTotal = &ProjectTotal{}
	err = db.Debug().Table("project").Select("id,donations,total").Where("id=?",
		db.Debug().Table("bot").Select("project_id").Where("id=?", BotId).SubQuery()).Scan(projectTotal).Error
	return
}

/**
 * @Description: 更新project的捐赠total和收到的捐赠笔数
 * @receiver projTotal
 * @param projectTotal
 * @return err
 */
func (projTotal *ProjectTotal) UpdateProjectTotal(projectTotal *ProjectTotal) (err error) {
	err = db.Debug().Table("project").Save(projectTotal).Error
	return
}

/**
 * @Description: 统计一个用户有获得了多少笔来自不同项目的捐赠捐赠
 * @receiver proj
 * @param userId
 * @return donations
 * @return err
 */
func (proj *Project) SumProjectDonationsByUserId(userId int64) (donations int64, err error) {
	type Result struct {
		Total int64
	}
	var result Result
	err = db.Debug().Table("project").Select("sum(donations) as total").Where("id IN(?)",
		db.Debug().Table("member").Select("project_id").Where("user_id=?", userId).SubQuery()).Scan(&result).Error
	donations = result.Total
	return
}
