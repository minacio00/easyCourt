package service

import (
	"errors"
	"sync"
	"time"

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
	GetBookingByID(id int) (*model.ReadBooking, error)
	GetAllBookings(limit, offset int) (*[]model.ReadBooking, error)
	UpdateBooking(booking *model.Booking) error
	DeleteBooking(id int) error
	ResetBookings() error
}

type bookingService struct {
	repo          repository.BookingRepository
	user_repo     repository.UserRepository
	timeslot_repo repository.TimeslotRepository
}

func NewBookingService(repo repository.BookingRepository, user repository.UserRepository, timeslot repository.TimeslotRepository) BookingService {
	return &bookingService{repo, user, timeslot}
}

func (s *bookingService) CreateBooking(booking *model.Booking) error {
	if err := booking.Validate(); err != nil {
		return err
	}

	if err := s.repo.CheckTimeslotAvailability(booking); err != nil {
		return err
	}

	// Verify user exists
	if _, err := s.user_repo.GetUserByID(uint(booking.UserID)); err != nil {
		return errors.New("user not found")
	}

	// Verify timeslot exists
	if _, err := s.timeslot_repo.GetTimeslotByID(booking.TimeslotID); err != nil {
		return errors.New("timeslot not found")
	}

	// Use mutex to prevent race conditions
	mutex := getTimeslotMutex(booking.TimeslotID)
	if !mutex.TryLock() {
		mutex.Lock()
	}
	defer mutex.Unlock()

	// Set booking date
	booking.BookingDate = time.Now()

	// Create the booking
	return s.repo.CreateBooking(booking)
}

func (s *bookingService) GetBookingByID(id int) (*model.ReadBooking, error) {
	return s.repo.GetBookingByID(id)
}

func (s *bookingService) GetAllBookings(limit, offset int) (*[]model.ReadBooking, error) {
	return s.repo.GetAllBookings(limit, offset)
}

func (s *bookingService) UpdateBooking(booking *model.Booking) error {
	// Update booking date
	booking.BookingDate = time.Now()

	// Verify the new timeslot exists and is available
	// (only if timeslot is changing)
	oldBooking, err := s.GetBookingByID(booking.ID)
	if err != nil {
		return err
	}

	if oldBooking.TimeslotID != booking.TimeslotID {
		// Check if the new timeslot is available
		if err := s.repo.CheckTimeslotAvailability(booking); err != nil {
			return err
		}
	}

	// Update the booking
	return s.repo.UpdateBooking(booking)
}

func (s *bookingService) DeleteBooking(id int) error {
	return s.repo.DeleteBooking(id)
}

func (s *bookingService) ResetBookings() error {
	return s.repo.ResetBookings()
}
