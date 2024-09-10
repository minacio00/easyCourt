package model

import (
	"fmt"
	"time"
)

type Weekday string

const (
	Domingo      Weekday = "Domingo"
	SegundaFeira Weekday = "Segunda-feira"
	TercaFeira   Weekday = "Terça-feira"
	QuartaFeira  Weekday = "Quarta-feira"
	QuintaFeira  Weekday = "Quinta-feira"
	SextaFeira   Weekday = "Sexta-feira"
	Sabado       Weekday = "Sábado"
)

type Timeslot struct {
	ID        int       `gorm:"primaryKey"`
	CourtID   *int      `json:"court_id"`
	Court     Court     `gorm:"foreignKey:CourtID"`
	Day       Weekday   `json:"week_day" gorm:"type:week_days;not null"`
	StartTime time.Time `json:"start_time" gorm:"type:time"`
	EndTime   time.Time `json:"end_time" gorm:"type:time"`
	IsActive  bool      `gorm:"default:true"`
}

func (t *Timeslot) Validate() error {
	var validationErrors []string
	if t.CourtID == nil {
		validationErrors = append(validationErrors, "court_id is required")
	}
	if t.Day == "" {
		validationErrors = append(validationErrors, "week_day is required")
	}
	if t.StartTime.IsZero() {
		validationErrors = append(validationErrors, "start_time is required")
	}
	if t.EndTime.IsZero() {
		validationErrors = append(validationErrors, "end_time is required")
	}
	if t.EndTime.Before(t.StartTime) {
		validationErrors = append(validationErrors, "end_time cannot be before start_time")
	}
	if len(validationErrors) > 0 {
		return fmt.Errorf("validation errors %v", validationErrors)
	}
	return nil
}

type CreateTimeslot struct {
	CourtID   int       `json:"court_id"`
	Day       Weekday   `json:"week_day"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	IsActive  bool      `json:"isActive"`
}

func (c *CreateTimeslot) ConvertCreateTimeslotToTimeslot() *Timeslot {
	return &Timeslot{
		CourtID:   &c.CourtID,
		Day:       c.Day,
		StartTime: c.StartTime,
		EndTime:   c.EndTime,
		IsActive:  c.IsActive,
	}
}
