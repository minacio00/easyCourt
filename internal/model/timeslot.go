package model

import "time"

type Timeslot struct {
	ID        int       `gorm:"primaryKey"`
	CourtID   int       `json:"court_id"`
	Court     Court     `gorm:"foreignKey:CourtID"`
	DayID     int       `json:"day_id"` // Foreign key field for Day
	Day       Day       `gorm:"foreignKey:DayID"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	IsActive  bool      `gorm:"default:true"`
}
