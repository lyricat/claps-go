package dao

import (
	"claps-test/model"
	"github.com/jinzhu/gorm"
)

func init() {
	RegisterMigrateHandler(func(db *gorm.DB) error {

		if err := db.AutoMigrate(&model.Bot{}).Error; err != nil {
			return err
		}
		return nil
	})
}


func GetBotById(botId string) (bot *model.Bot, err error) {
	bot = &model.Bot{}
	err = db.Debug().Where("id=?", botId).Find(bot).Error
	return
}

//根据projectId获取所有的机器人Id
func ListBotDtosByProjectId(projectId int64) (botDto *[]model.BotDto, err error) {
	botDto = &[]model.BotDto{}
	err = db.Debug().Table("bot").Where("project_id=?", projectId).Scan(botDto).Error
	return
}

func GetBotDtoById(botId string) (botDto *model.BotDto, err error) {
	botDto = &model.BotDto{}
	err = db.Table("bot").Where("id=?", botId).Scan(botDto).Error
	return
}
