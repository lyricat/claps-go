package model

type Member struct {
	ProjectId int64 `gorm:"type:bigint;not null;primary_key"`
	UserId    int64 `gorm:"type:bigint;not null;primary_key"`
}
