package models

import "time"

type Project struct {
	ID            string    `gorm:"primaryKey" json:"id"`
	Slug          string    `gorm:"not null" json:"slug"`
	RepositoryUrl string    `gorm:"not null" json:"repositoryUrl"`
	Owner         string    `gorm:"not null" json:"owner"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
