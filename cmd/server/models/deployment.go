package models

import "time"

type Deployment struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	ProjectID string    `gorm:"index;not null" json:"projectId"`
	Status    string    `gorm:"not null" json:"status"`
	Message   string    `json:"message"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
