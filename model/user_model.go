package model

import "github.com/jinzhu/gorm"

/**
 * @Description:注册自动迁移函数
 */
func init() {
	RegisterMigrateHandler(func(db *gorm.DB) error {

		if err := db.AutoMigrate(&User{}).Error; err != nil {
			return err
		}
		return nil
	})
}

type User struct {
	Id            int64  `json:"id,omitempty" gorm:"type:bigint;primary_key;not null;"`
	Name          string `json:"name,omitempty" gorm:"type:varchar(50);unique_index:name_UNIQUE;not null"`
	DisplayName   string `json:"display_name,omitempty" gorm:"type:varchar(50);default:null"`
	Email         string `json:"email,omitempty" gorm:"type:varchar(50);unique_index:email_UNIQUE;not null"`
	AvatarUrl     string `json:"avatar_url,omitempty" gorm:"type:varchar(100);default:null"`
	MixinId       string `json:"mixin_id,omitempty" gorm:"type:varchar(50);default:null"`
	WithdrawalWay string `json:"withdrawal_way,omitempty" gorm:"type:char;default:0"`
}

type UserMixinId struct {
	Id      int64  `json:"id,omitempty" gorm:"type:bigint;primary_key;not null;"`
	MixinId string `json:"mixin_id,omitempty" gorm:"type:varchar(50);default:null"`
}

const (
	WithdrawByUser  = "0" //用户手动点提现
	WithdrawByClaps = "1" //捐赠到了自动提现到绑定的mixin账号
)

var (
	USER        *User
	USERMIXINID *UserMixinId
)

/**
 * @Description: 根据id获取user完整信息
 * @receiver u
 * @param id
 * @return user
 * @return err
 */
func (u *User) GetUserById(id int64) (user *User, err error) {
	user = &User{}
	err = db.Debug().Where("id=?", id).Find(user).Error
	return
}

/**
 * @Description: 从数据库中通过ID获取user信息,存储在user中,引用传值
 * @receiver user
 * @param id
 * @return userData
 * @return err
 */
func (user *UserMixinId) GetMixinIdById(id int64) (userData *UserMixinId, err error) {
	userData = &UserMixinId{}
	err = db.Debug().Table("user").Where("id = ?", id).Scan(userData).Error
	return
}

/**
 * @Description: 不管记录是否找到，都将参数赋值给 struct 并保存至数据库
 * @receiver user
 * @param userData
 * @return err
 */
func (user *User) InsertOrUpdateUser(userData *User) (err error) {

	var cnt int64
	db.Debug().Table("user").Where("id = ?", userData.Id).Count(&cnt)
	if cnt == 0 {
		err = db.Debug().Create(userData).Error
		return
	} else {
		db.Debug().Model(&userData).Omit("mixin_id").Updates(userData)
	}
	return
}

/**
 * @Description: 通过projectId获取一个项目的所有成员信息
 * @receiver user
 * @param projectId
 * @return users
 * @return err
 */
func (user *User) ListMembersByProjectId(projectId int64) (users *[]User, err error) {

	users = &[]User{}
	err = db.Debug().Where("id IN (?)",
		db.Debug().Table("member").Select("user_id").Where("project_id=?", projectId).SubQuery()).Find(users).Error
	return
}

/**
 * @Description: 根据user_id更新表中的mixin_id信息
 * @receiver user
 * @param userId
 * @param mixinId
 * @return err
 */
func (user *User) UpdateUserMixinId(userId int64, mixinId string) (err error) {

	err = db.Debug().Model(&User{}).Where("id = ?", userId).Update("mixin_id", mixinId).Error
	return
}

/**
 * @Description: 通过userId更新WithdrawalWay
 * @receiver user
 * @param userId
 * @param withdrawalWay
 * @return err
 */
func (user *User) UpdateUserWithdrawalWay(userId int64, withdrawalWay string) (err error) {
	err = db.Debug().Model(&User{}).Where("id = ?", userId).Update("withdrawal_way", withdrawalWay).Error
	return
}
