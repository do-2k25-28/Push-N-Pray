package models

import "time"

const (
	ServiceProvisioning = "provisioning"
	ServiceRunning      = "running"
	ServiceError        = "error"
	ServiceDeleted      = "deleted"
)

type ManagedService struct {
	ID            string    `gorm:"primaryKey" json:"id"`
	ProjectID     string    `gorm:"index;not null" json:"projectId"`
	Name          string    `gorm:"not null" json:"name"`
	Type          string    `gorm:"not null" json:"type"`
	Status        string    `gorm:"not null" json:"status"`
	ContainerName string    `json:"containerName,omitempty"`
	VolumeName    string    `json:"volumeName,omitempty"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
