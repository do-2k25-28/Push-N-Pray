package models

import "time"

type PersonalAccessToken struct {
	ID        string     `gorm:"primaryKey" json:"id"`
	Name      string     `gorm:"not null" json:"name"`
	Hash      string     `gorm:"not null;uniqueIndex" json:"-"`
	Owner     string     `gorm:"index;not null" json:"owner"`
	ExpiresAt *time.Time `json:"expiresAt"`
	CreatedAt time.Time  `json:"createdAt"`
}
