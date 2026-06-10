package models

type EnvVar struct {
	Project string `gorm:"varchar(255);primaryKey"`
	Name    string `gorm:"varchar(255);primaryKey"`
	Value   string `gorm:"type:text"`
}
