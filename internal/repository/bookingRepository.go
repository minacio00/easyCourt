package repository

import (
	"errors"
	"fmt"

	"github.com/minacio00/easyCourt/internal/model"
	"gorm.io/gorm"
)

type BookingRepository interface {
	CreateBooking(booking *model.Booking) error
	GetBookingByID(id int) (*model.Booking, error)
	GetAllBookings(limit, offset int) (*[]model.ReadBooking, error)
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
	var readTimeSlot model.ReadTimeslot
	if err := r.db.Model(&model.Timeslot{}).Where("id = ?", booking.TimeslotID).First(&readTimeSlot).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("timeslot with ID %d not found", booking.TimeslotID)
		}
		return fmt.Errorf("error fetching timeslot: %w", err)
	}

	timeslot, err := readTimeSlot.ToTimeslot()
	if err != nil {
		return fmt.Errorf("error converting ReadTimeslot to Timeslot: %w", err)
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(booking).Error; err != nil {
			return fmt.Errorf("error creating booking: %w", err)
		}

		if err := tx.Model(&timeslot).Update("booking_id", booking.ID).Error; err != nil {
			return fmt.Errorf("error updating timeslot with booking ID: %w", err)
		}

		return nil
	})
}

func (r *bookingRepository) CheckTimeslotAvailability(booking *model.Booking) error {
	var book model.Booking
	result := r.db.Where(&model.Booking{TimeslotID: booking.TimeslotID}).First(&book)
	if result.Error == nil {
		// A booking was found, so the timeslot is not available
		return fmt.Errorf("timeslot already booked %d", booking.TimeslotID)
	}

	if result.Error == gorm.ErrRecordNotFound {
		// No booking was found, so the timeslot is available
		return nil
	}

	return fmt.Errorf("error checking timeslot availability: %w", result.Error)
}

func (r *bookingRepository) GetBookingByID(id int) (*model.Booking, error) {
	var booking model.Booking
	if err := r.db.Preload("User").Preload("Timeslot").First(&booking, id).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) GetAllBookings(limit, offset int) (*[]model.ReadBooking, error) {
	var bookings []model.ReadBooking
	// var bookings []map[string]interface{}
	query := r.db.
		Preload("User").
		Preload("Timeslot", func(db *gorm.DB) *gorm.DB {
			return db.Table("timeslots")
		}).
		Table("bookings")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	if err := query.Find(&bookings).Error; err != nil {
		return nil, err
	}
	return &bookings, nil
}

func (r *bookingRepository) UpdateBooking(booking *model.Booking) error {
	return r.db.Save(booking).Error
}

func (r *bookingRepository) DeleteBooking(id int) error {
	return r.db.Delete(&model.Booking{}, id).Error
}
