package model

//用来做重要数据的持久化,保存key-value
type Property struct {
	Key   string `gorm:"type:varchar(50);not null;default:0;primary_key;"`
	Value string `gorm:"type:varchar(50);not null;default:0"`
}
