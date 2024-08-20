package service

import (
	"github.com/minacio00/easyCourt/internal/model"
	"github.com/minacio00/easyCourt/internal/repository"
)

type BookingService interface {
	CreateBooking(booking *model.Booking) error
	GetBookingByID(id int) (*model.Booking, error)
	GetAllBookings() ([]model.Booking, error)
	UpdateBooking(booking *model.Booking) error
	DeleteBooking(id int) error
}

type bookingService struct {
	repo repository.BookingRepository
}

func NewBookingService(repo repository.BookingRepository) BookingService {
	return &bookingService{repo}
}

func (s *bookingService) CreateBooking(booking *model.Booking) error {
	// Add any additional logic before creating a booking, if necessary
	return s.repo.CreateBooking(booking)
}

func (s *bookingService) GetBookingByID(id int) (*model.Booking, error) {
	return s.repo.GetBookingByID(id)
}

func (s *bookingService) GetAllBookings() ([]model.Booking, error) {
	return s.repo.GetAllBookings()
}

func (s *bookingService) UpdateBooking(booking *model.Booking) error {
	// Add any additional logic before updating a booking, if necessary
	return s.repo.UpdateBooking(booking)
}

func (s *bookingService) DeleteBooking(id int) error {
	// Add any additional logic before deleting a booking, if necessary
	return s.repo.DeleteBooking(id)
}
