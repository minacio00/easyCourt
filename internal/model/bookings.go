package model

import "time"

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

type CreateBooking struct {
	UserID          int     `json:"user_id"`
	Opponent        string  `json:"oponnent_name"`
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
