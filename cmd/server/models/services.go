package models

type PostgresService struct {
	Project  string `gorm:"varchar(255);primaryKey"`
	Name     string `gorm:"varchar(255);primaryKey"`
	Password string `gorm:"type:varchar(255)"`
}

type RedisService struct {
	Project  string `gorm:"varchar(255);primaryKey"`
	Name     string `gorm:"varchar(255);primaryKey"`
	Password string `gorm:"type:varchar(255)"`
}
