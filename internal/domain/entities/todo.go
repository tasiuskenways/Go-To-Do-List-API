package entities

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	ID          string `gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Title       string `gorm:"not null"`
	Description string `gorm:"not null"`
	UserID      string `gorm:"not null;type:uuid;index"`
	User        User   `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}