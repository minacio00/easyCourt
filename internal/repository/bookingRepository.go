package repository

import (
	"github.com/minacio00/easyCourt/internal/model"
	"gorm.io/gorm"
)

type BookingRepository interface {
	CreateBooking(booking *model.Booking) error
	GetBookingByID(id int) (*model.Booking, error)
	GetAllBookings(limit, offset int) ([]model.Booking, error)
	UpdateBooking(booking *model.Booking) error
	DeleteBooking(id int) error
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
