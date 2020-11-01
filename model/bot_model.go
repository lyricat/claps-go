package model

import "github.com/jinzhu/gorm"

/**
 * @Description:注册自动迁移函数
 */
func init() {
	RegisterMigrateHandler(func(db *gorm.DB) error {

		if err := db.AutoMigrate(&Bot{}).Error; err != nil {
			return err
		}
		return nil
	})
}

type Bot struct {
	Id           string `gorm:"type:varchar(50);primary_key;not null"`
	ProjectId    int64  `gorm:"type:bigint;primary_key;not null"`
	Distribution string `gorm:"type:char;primary_key;not null"`
	SessionId    string `gorm:"type:varchar(50);not null;unique_index:session_id_UNIQUE"`
	Pin          string `gorm:"type:varchar(6);not null"`
	PinToken     string `gorm:"type:varchar(200);not null;unique_index:pin_token_UNIQUE"`
	PrivateKey   string `gorm:"type:text;not null"`
}

type BotDto struct {
	Id           string `json:"id,omitempty" gorm:"type:varchar(50);primary_key;not null"`
	Distribution string `json:"distribution,omitempty" gorm:"type:char;primary_key;not null"`
}

const (
	MericoAlgorithm = "0" //Merico算法
	Commits         = "1" //commit数量
	ChangedLines    = "2" //代码行数
	IdenticalAmount = "3" //平均分配
)

var (
	BOT    *Bot
	BOTDTO *BotDto
)

/**
 * @Description: 通过botId获取bot信息
 * @receiver bot
 * @param botId
 * @return botData
 * @return err
 */
func (bot *Bot) GetBotById(botId string) (botData *Bot, err error) {
	botData = &Bot{}
	err = db.Debug().Where("id=?", botId).Find(botData).Error
	return
}

/**
 * @Description: 根据projectId获取所有的BotId
 * @receiver bot
 * @param projectId
 * @return botDto
 * @return err
 */
func (bot *BotDto) ListBotDtosByProjectId(projectId int64) (botDto *[]BotDto, err error) {
	botDto = &[]BotDto{}
	err = db.Debug().Table("bot").Where("project_id=?", projectId).Scan(botDto).Error
	return
}

/**
 * @Description: 根据BotId获取所有的BotId
 * @receiver bot
 * @param botId
 * @return botDto
 * @return err
 */
func (bot *BotDto) GetBotDtoById(botId string) (botDto *BotDto, err error) {
	botDto = &BotDto{}
	err = db.Table("bot").Where("id=?", botId).Scan(botDto).Error
	return
}
