package models

type RefreshToken struct {
	Owner string `gorm:"type:varchar(255);primaryKey"`
	Token string `gorm:"type:varchar(255);primaryKey"`
}
