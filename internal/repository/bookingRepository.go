package repository

import (
	"errors"
	"fmt"

	"github.com/minacio00/easyCourt/internal/model"
	"gorm.io/gorm"
)

type BookingRepository interface {
	CreateBooking(booking *model.Booking) error
	GetBookingByID(id int) (*model.ReadBooking, error)
	GetAllBookings(limit, offset int) (*[]model.ReadBooking, error)
	UpdateBooking(booking *model.Booking) error
	DeleteBooking(id int) error
	CheckTimeslotAvailability(booking *model.Booking) error
	ResetBookings() error
	GetUserBookings(userId int, limit, offset int) (*[]model.ReadBooking, error)
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db}
}

func (r *bookingRepository) CreateBooking(booking *model.Booking) error {
	// Verify timeslot exists
	var timeslot model.Timeslot
	if err := r.db.Where("id = ?", booking.TimeslotID).First(&timeslot).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("timeslot with ID %d not found", booking.TimeslotID)
		}
		return fmt.Errorf("error fetching timeslot: %w", err)
	}

	// Create the booking in a single transaction
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(booking).Error; err != nil {
			return fmt.Errorf("error creating booking: %w", err)
		}
		return nil
	})
}

func (r *bookingRepository) CheckTimeslotAvailability(booking *model.Booking) error {
	// Skip check if updating an existing booking for the same timeslot
	if booking.ID != 0 {
		var existingBooking model.Booking
		if err := r.db.First(&existingBooking, booking.ID).Error; err == nil {
			if existingBooking.TimeslotID == booking.TimeslotID {
				// Same timeslot, so it's available for this booking
				return nil
			}
		}
	}

	// Check if any booking exists for this timeslot
	var count int64
	if err := r.db.Model(&model.Booking{}).
		Where("timeslot_id = ?", booking.TimeslotID).
		Not("id = ?", booking.ID). // Exclude current booking if updating
		Count(&count).Error; err != nil {
		return fmt.Errorf("error checking timeslot availability: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("timeslot already booked %d", booking.TimeslotID)
	}

	return nil
}

func (r *bookingRepository) GetBookingByID(id int) (*model.ReadBooking, error) {
	var booking model.ReadBooking
	if err := r.db.Preload("User").
		Preload("Timeslot", func(db *gorm.DB) *gorm.DB {
			return db.Table("timeslots")
		}).
		Table("bookings").
		First(&booking, id).Error; err != nil {
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) GetUserBookings(userId int, limit, offset int) (*[]model.ReadBooking, error) {
	var bookings []model.ReadBooking

	query := r.db.
		Preload("User").
		Preload("Timeslot").
		Where("user_id = ?", userId)
	if limit > 0 {
		query = query.Limit(limit)
	}

	if offset > 0 {
		query = query.Offset(offset)
	}
	result := query.Find(&bookings)

	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return &[]model.ReadBooking{}, nil
	}
	return &bookings, nil
}

func (r *bookingRepository) GetAllBookings(limit, offset int) (*[]model.ReadBooking, error) {
	var bookings []model.ReadBooking
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

func (r *bookingRepository) ResetBookings() error {
	return r.db.Delete(&model.Booking{}, "1 = 1").Error
}
