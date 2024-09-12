package service

import (
	"sync"

	"github.com/minacio00/easyCourt/internal/model"
	"github.com/minacio00/easyCourt/internal/repository"
)

var timeslotLocks map[int]*sync.Mutex = make(map[int]*sync.Mutex)

func getTimeslotMutex(timeslotID int) *sync.Mutex {
	if _, exists := timeslotLocks[timeslotID]; !exists {
		timeslotLocks[timeslotID] = &sync.Mutex{}
	}
	return timeslotLocks[timeslotID]
}

type BookingService interface {
	CreateBooking(booking *model.Booking) error
	GetBookingByID(id int) (*model.Booking, error)
	GetAllBookings(limit, offeset int) ([]model.Booking, error)
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
	if err := booking.Validate(); err != nil {
		return err
	}
	if err := s.repo.CheckTimeslotAvailability(booking); err != nil {
		return err
	}

	mutex := getTimeslotMutex(booking.TimeslotID)
	if !mutex.TryLock() {
		mutex.Lock()
	}
	defer mutex.Unlock()

	return s.repo.CreateBooking(booking)
}

func (s *bookingService) GetBookingByID(id int) (*model.Booking, error) {
	return s.repo.GetBookingByID(id)
}

func (s *bookingService) GetAllBookings(limit, offeset int) ([]model.Booking, error) {
	return s.repo.GetAllBookings(limit, offeset)
}

func (s *bookingService) UpdateBooking(booking *model.Booking) error {
	// Add any additional logic before updating a booking, if necessary
	return s.repo.UpdateBooking(booking)
}

func (s *bookingService) DeleteBooking(id int) error {
	// Add any additional logic before deleting a booking, if necessary
	return s.repo.DeleteBooking(id)
}
