package model

// Room represents a room in the system
type Room struct {
	ID          uint   `json:"id" gorm:"primaryKey"`
	RoomName    string `json:"room_name" gorm:"not null;unique;size:255"`
	Capacity    int    `json:"capacity" gorm:"not null"`
	Location    string `json:"location" gorm:"size:255"`
	Description string `json:"description" gorm:"size:255"`
	Status      string `json:"status" gorm:"size:50"` // Example: "available", "booked", "under maintenance"
}
