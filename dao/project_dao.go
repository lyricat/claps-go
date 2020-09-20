package dao

import (
	"claps-test/model"
	"github.com/jinzhu/gorm"
)

func init() {
	RegisterMigrateHandler(func(db *gorm.DB) error {

		if err := db.AutoMigrate(&model.Project{}).Error; err != nil {
			return err
		}
		return nil
	})
}


//获取所有项目
func ListProjectsAll() (projects *[]model.Project, err error) {

	projects = &[]model.Project{}
	err = db.Debug().Find(projects).Error
	return
}

//通过项目名字获取项目
func GetProjectByName(name string) (project *model.Project, err error) {

	project = &model.Project{}
	err = db.Debug().Where("name=?", name).Find(&project).Error
	return
}

//根据userid获取所有项目
func ListProjectsByUserId(userId int64) (projects *[]model.Project, err error) {
	projects = &[]model.Project{}
	err = db.Debug().Where("id IN(?)",
		db.Debug().Table("member").Select("project_id").Where("user_id=?", userId).SubQuery()).Find(projects).Error
	return
}

func GetProjectTotalByBotId(BotId string) (projectTotal *model.ProjectTotal, err error) {
	projectTotal = &model.ProjectTotal{}
	err = db.Debug().Table("project").Select("id,donations,total").Where("id=?",
		db.Debug().Table("bot").Select("project_id").Where("id=?", BotId).SubQuery()).Scan(projectTotal).Error
	return
}

func UpdateProjectTotal(projectTotal *model.ProjectTotal) (err error) {
	err = db.Debug().Table("project").Save(projectTotal).Error
	return
}

func SumProjectDonationsByUserId(userId int64) (donations int64, err error) {
	type Result struct {
		Total int64
	}
	var result Result
	err = db.Debug().Table("project").Select("sum(donations) as total").Where("id IN(?)",
		db.Debug().Table("member").Select("project_id").Where("user_id=?",userId).SubQuery()).Scan(&result).Error
	donations = result.Total
	return
}