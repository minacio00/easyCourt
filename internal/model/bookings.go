package model

import (
	"fmt"
	"time"
)

type Booking struct {
	ID              int       `gorm:"primaryKey" json:"-"`
	UserID          int       `json:"user_id"`           // Foreign key to User
	User            User      `gorm:"foreignKey:UserID"` // Relationship with User
	Opponent        string    `json:"oponnent_name"`
	Partner         *string   `json:"partner_name"`
	OpponentPartner *string   `json:"opponent_partner"`
	TimeslotID      int       `json:"timeslot_id"`                    // Foreign key to Timeslot
	Timeslot        Timeslot  `gorm:"foreignKey:TimeslotID" json:"-"` // Relationship with Timeslot
	BookingDate     time.Time `json:"booking_date"`                   // Date of the booking
	IsSinglesGame   bool      `json:"singles_flag"`
}

func (bk *Booking) Validate() error {
	var errors []string
	if bk.UserID == 0 {
		errors = append(errors, "user_id is required")
	}
	if bk.Opponent == "" {
		errors = append(errors, "oponnent is required")
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
