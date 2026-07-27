package models

import "time"

type Permission struct {
	ID        uint      `gorm:"primaryKey"`
	Key       string    `gorm:"uniqueIndex;not null"`
	Name      string    `gorm:"not null"`
	Group     string    `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
