package entity

import "time"

type AgentAllocationQueue struct {
	ID         uint      `gorm:"primaryKey"`
	RoomID     string    `gorm:"size:50;not null"`
	AgentID    string    `gorm:"size:50"`
	IsResolved bool      `gorm:"default:false"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime:false"`
}
