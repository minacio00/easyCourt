package model

import "time"

type Booking struct {
	ID          int       `gorm:"primaryKey"`
	UserID      int       `json:"user_id"`           // Foreign key to User
	User        User      `gorm:"foreignKey:UserID"` // Relationship with User
	Opponent    string    `json:"oponnent_name"`
	TimeslotID  int       `json:"timeslot_id"`           // Foreign key to Timeslot
	Timeslot    Timeslot  `gorm:"foreignKey:TimeslotID"` // Relationship with Timeslot
	BookingDate time.Time `json:"booking_date"`          // Date of the booking
}
