package model

type User struct {
	Id          int64  `json:"id,omitempty" gorm:"type:bigint;primary_key;not null;"`
	Name        string `json:"name,omitempty" gorm:"type:varchar(50);unique_index:name_UNIQUE;not null"`
	DisplayName string `json:"display_name,omitempty" gorm:"type:varchar(50);default:null"`
	Email       string `json:"email,omitempty" gorm:"type:varchar(50);unique_index:email_UNIQUE;not null"`
	AvatarUrl   string `json:"avatar_url,omitempty" gorm:"type:varchar(100);default:null"`
	//UserId string `json:"user_id" gorm:"type:varchar(50);default:null"`
	MixinId string `json:"mixin_id,omitempty" gorm:"type:varchar(50);default:null"`
}

type UserMixinId struct {
	Id      int64  `json:"id,omitempty" gorm:"type:bigint;primary_key;not null;"`
	MixinId string `json:"mixin_id,omitempty" gorm:"type:varchar(50);default:null"`
}
