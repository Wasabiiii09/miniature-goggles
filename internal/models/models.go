package models

import (
	"time"
)

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}

type Comment struct {
	ID        uint `gorm:"primaryKey"`
	Text      string `gorm:"not null"`
	UserID    uint `gorm:"not null"`
	User      User `gorm:"foreignKey:UserID"`
	CreatedAt time.Time
}