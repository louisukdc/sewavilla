package model

import (
	"time"
)

type Reservation struct {
	ID              uint      `json:"id" gorm:"primaryKey"`
	UserID          *uint     `json:"user_id"` // Mengizinkan NULL
	RoomID          *uint     `json:"room_id"` // Mengizinkan NULL
	BlogID          *uint     `json:"blog_id"` // Mengizinkan NULL
	ReservationDate time.Time `json:"reservation_date" gorm:"not null"`
	StartTime       time.Time `json:"start_time" gorm:"not null"`
	EndTime         time.Time `json:"end_time" gorm:"not null"`
	Status          string    `json:"status" gorm:"type:enum('pending','approved','completed','cancelled');default:'pending';not null"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`

	// Associations
	User *User `json:"user" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Room *Room `json:"room" gorm:"foreignKey:RoomID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Blog *Blog `json:"blog" gorm:"foreignKey:BlogID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
