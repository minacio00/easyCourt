package repository

import (
	"fmt"

	"github.com/minacio00/easyCourt/internal/model"
	"gorm.io/gorm"
)

type BookingRepository interface {
	CreateBooking(booking *model.Booking) error
	GetBookingByID(id int) (*model.Booking, error)
	GetAllBookings(limit, offset int) ([]model.Booking, error)
	UpdateBooking(booking *model.Booking) error
	DeleteBooking(id int) error
	CheckTimeslotAvailability(booking *model.Booking) error
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db}
}

func (r *bookingRepository) CreateBooking(booking *model.Booking) error {
	return r.db.Create(booking).Error
}

func (r *bookingRepository) CheckTimeslotAvailability(booking *model.Booking) error {
	// var book model.CreateBooking
	// result := r.db.Preload("Timeslot").Where("timeslot_id = ?", booking.TimeslotID).First(&book)
	// if result.RowsAffected > 0 {
	// 	return fmt.Errorf("horário não disponível para a quadra %d", &booking.TimeslotID)
	// }
	query := `
		SELECT b.id, b.timeslot_id, t.start_time, t.end_time
		from bookings b
		JOIN timeslots t ON b.timeslot_id = t.id
		WHERE b.timeslot_id = ?
		LIMIT 1
	`
	var result map[string]interface{}
	err := r.db.Raw(query, booking.TimeslotID).Scan(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Timeslot is available
			return nil
		}
		// Other database error occurred
		return fmt.Errorf("error checking timeslot availability: %w", err)
	}
	return fmt.Errorf("horário não disponível para a quadra %d", booking.TimeslotID)
}

func (r *bookingRepository) GetBookingByID(id int) (*model.Booking, error) {
	var booking model.Booking
	if err := r.db.Preload("User").Preload("Timeslot").First(&booking, id).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) GetAllBookings(limit, offset int) ([]model.Booking, error) {
	var bookings []model.Booking
	if err := r.db.Preload("User").Preload("Timeslot").
		Limit(limit).
		Offset(offset).
		Find(&bookings).
		Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepository) UpdateBooking(booking *model.Booking) error {
	return r.db.Save(booking).Error
}

func (r *bookingRepository) DeleteBooking(id int) error {
	return r.db.Delete(&model.Booking{}, id).Error
}
