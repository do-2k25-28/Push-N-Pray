package models

import "time"

type ManagedServiceType string

const (
	PostgresService ManagedServiceType = "postgres"
)

type ManagedServiceStatus string

const (
	ServiceProvisioning ManagedServiceStatus = "provisioning"
	ServiceRunning      ManagedServiceStatus = "running"
	ServiceError        ManagedServiceStatus = "error"
	ServiceDeleted      ManagedServiceStatus = "deleted"
)

type ManagedService struct {
	ID            string               `gorm:"primaryKey" json:"id"`
	ProjectID     string               `gorm:"index;not null" json:"projectId"`
	Name          string               `gorm:"not null" json:"name"`
	Type          ManagedServiceType   `gorm:"not null" json:"type"`
	Status        ManagedServiceStatus `gorm:"not null" json:"status"`
	ContainerName string               `json:"containerName,omitempty"`
	VolumeName    string               `json:"volumeName,omitempty"`
	Message       string               `json:"message,omitempty"`
	CreatedAt     time.Time            `json:"createdAt"`
	UpdatedAt     time.Time            `json:"updatedAt"`
}
