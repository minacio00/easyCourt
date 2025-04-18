package model

import (
	"fmt"
	"log"
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
	ID        int       `gorm:"primaryKey" json:"id"`
	CourtID   *int      `json:"court_id"`
	Court     Court     `gorm:"foreignKey:CourtID"`
	Day       Weekday   `json:"week_day" gorm:"type:week_days;not null"`
	StartTime time.Time `json:"start_time" gorm:"type:time"`
	EndTime   time.Time `json:"end_time" gorm:"type:time"`
	IsActive  bool      `json:"is_active"`
	Booking   *Booking  `json:"booking,omitempty"`
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

func (t *Timeslot) ToReadTimeslot() {
}

type CreateTimeslot struct {
	CourtID   *int    `json:"court_id"`
	Day       Weekday `json:"week_day"`
	StartTime string  `json:"start_time"`
	EndTime   string  `json:"end_time"`
	IsActive  bool    `json:"is_active"`
}

func (c *CreateTimeslot) ConvertCreateTimeslotToTimeslot() (*Timeslot, error) {
	fmt.Println(c.StartTime)
	startTime, err := time.Parse(time.TimeOnly, c.StartTime)
	fmt.Println(startTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start_time format: %v", err)
	}
	fmt.Printf("endtime: %v \n", c.EndTime)
	endTime, err := time.Parse("15:04:05", c.EndTime)
	fmt.Printf("endtime var: %v \n", endTime)

	fmt.Println(err)

	if err != nil {
		fmt.Println("errored during conversion")
		return nil, fmt.Errorf("invalid end_time format: %v", err)
	}
	return &Timeslot{
		CourtID:   c.CourtID,
		Day:       c.Day,
		StartTime: startTime,
		EndTime:   endTime,
		IsActive:  c.IsActive,
	}, nil
}

type ReadTimeslot struct {
	ID        int          `json:"id" gorm:"timeslot_id"`
	CourtID   *int         `json:"court_id" gorm:"column:court_id"`
	Court     Court        `json:"court"`
	Day       Weekday      `json:"week_day"`
	StartTime string       `json:"start_time" gorm:"column:start_time"`
	EndTime   string       `json:"end_time" gorm:"column:end_time"`
	IsActive  bool         `json:"is_active"`
	Booking   *ReadBooking `json:"booking,omitempty" gorm:"-"`
}

func (rt *ReadTimeslot) ToTimeslot() (*Timeslot, error) {
	parserdStart, err := time.Parse(time.RFC3339, rt.StartTime)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("Error Parsing  startTime:", rt.StartTime)
		return nil, fmt.Errorf("invalid start_time format: %v", err)
	}
	parserdEnd, err := time.Parse(time.RFC3339, rt.EndTime)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("Error Parsing  endTime:", rt.StartTime)
		return nil, fmt.Errorf("invalid end_time format: %v", err)
	}

	ts := &Timeslot{
		ID:        rt.ID,
		CourtID:   rt.CourtID,
		Court:     rt.Court,
		Day:       rt.Day,
		StartTime: parserdStart,
		EndTime:   parserdEnd,
		IsActive:  rt.IsActive,
	}

	if rt.Booking != nil {
		booking := &Booking{
			ID:         rt.Booking.ID,
			UserID:     rt.Booking.UserID,
			TimeslotID: rt.ID,
		}
		ts.Booking = booking
	}

	return ts, nil
}
