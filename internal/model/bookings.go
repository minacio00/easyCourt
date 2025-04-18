package model

import (
	"fmt"
	"time"
)

type Booking struct {
	ID              int       `gorm:"primaryKey"`
	UserID          int       `json:"user_id" gorm:"index"`
	User            User      `json:"user"`
	Opponent        string    `json:"opponent_name"`
	Partner         *string   `json:"partner_name"`
	OpponentPartner *string   `json:"opponent_partner"`
	TimeslotID      int       `json:"timeslot_id"`
	BookingDate     time.Time `json:"booking_date"` // Date of the booking
	IsSinglesGame   bool      `json:"singles_flag"`
	Timeslot        Timeslot  `gorm:"foreignKey:TimeslotID" json:"ommitempty"`
}

func (bk *Booking) Validate() error {
	var errors []string
	if bk.UserID == 0 {
		errors = append(errors, "user_id is required")
	}
	if bk.Opponent == "" {
		errors = append(errors, "opponent is required")
	}
	if bk.Partner == nil && !bk.IsSinglesGame {
		errors = append(errors, "partner_name is required for doubles games")
	}
	if bk.OpponentPartner == nil && !bk.IsSinglesGame {
		errors = append(errors, "opponent_partner is required for doubles games")
	}
	if bk.TimeslotID == 0 {
		errors = append(errors, "timeslot_id is required")
	}

	if len(errors) > 0 {
		return fmt.Errorf("validation errors: %v", errors)
	}
	return nil
}

type CreateBooking struct {
	UserID          int     `json:"user_id"`
	Opponent        string  `json:"opponent_name"`
	Partner         *string `json:"partner_name"`
	OpponentPartner *string `json:"opponent_partner"`
	TimeslotID      int     `json:"timeslot_id"`
	IsSinglesGame   bool    `json:"singles_flag"`
}

func (cb *CreateBooking) ConvertToBooking() *Booking {
	return &Booking{
		UserID:          cb.UserID,
		Opponent:        cb.Opponent,
		Partner:         cb.Partner,
		OpponentPartner: cb.OpponentPartner,
		TimeslotID:      cb.TimeslotID,
		BookingDate:     time.Now(),
		IsSinglesGame:   cb.IsSinglesGame,
	}
}

type ReadBooking struct {
	ID              int          `gorm:"primaryKey" json:"id"`
	UserID          int          `json:"user_id"`
	User            User         `json:"user"`
	Opponent        string       `json:"opponent_name"`
	Partner         *string      `json:"partner_name"`
	OpponentPartner *string      `json:"opponent_partner"`
	TimeslotID      int          `json:"timeslot_id"`
	Timeslot        ReadTimeslot `json:"timeslot"`
	BookingDate     time.Time    `json:"booking_date"`
	IsSinglesGame   bool         `json:"singles_flag"`
}

func (b *ReadBooking) ToBooking() *Booking {
	return &Booking{
		ID:              b.ID,
		User:            b.User,
		Opponent:        b.Opponent,
		Partner:         b.Partner,
		OpponentPartner: b.OpponentPartner,
		BookingDate:     b.BookingDate,
		IsSinglesGame:   b.IsSinglesGame,
		TimeslotID:      b.TimeslotID,
	}
}

func (b *Booking) ToResponse() *ReadBooking {
	if b != nil {
		return &ReadBooking{
			ID:              b.ID,
			User:            b.User,
			Opponent:        b.Opponent,
			Partner:         b.Partner,
			OpponentPartner: b.OpponentPartner,
			BookingDate:     b.BookingDate,
			IsSinglesGame:   b.IsSinglesGame,
			TimeslotID:      b.TimeslotID,
		}
	}
	return nil
}
