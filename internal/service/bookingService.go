package service

import (
	"errors"
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
	GetBookingByID(id int) (*model.ReadBooking, error)
	GetAllBookings(limit, offeset int) (*[]model.ReadBooking, error)
	UpdateBooking(booking *model.Booking) error
	DeleteBooking(id int) error
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
	if _, err := s.user_repo.GetUserByID(uint(booking.UserID)); err != nil {
		return errors.New("user not found")
	}

	mutex := getTimeslotMutex(booking.TimeslotID)
	if !mutex.TryLock() {
		mutex.Lock()
	}
	defer mutex.Unlock()

	return s.repo.CreateBooking(booking)
}

func (s *bookingService) GetBookingByID(id int) (*model.ReadBooking, error) {
	return s.repo.GetBookingByID(id)
}

func (s *bookingService) GetAllBookings(limit, offeset int) (*[]model.ReadBooking, error) {
	return s.repo.GetAllBookings(limit, offeset)
}

func (s *bookingService) UpdateBooking(booking *model.Booking) error {
	// Update the booking
	if err := s.repo.UpdateBooking(booking); err != nil {
		return err
	}

	// Clear the old timeslot association
	oldTs, err := s.timeslot_repo.GetTimeslotByBookingId(uint(booking.ID))
	if err != nil {
		return err
	}
	if oldTs != nil {
		oldTs.Booking_id = nil
		oldTimeslot, err := oldTs.ToTimeslot()
		if err != nil {
			return err
		}
		if err = s.timeslot_repo.UpdateTimeslot(oldTimeslot); err != nil {
			return err
		}
	}

	// Set the new timeslot association
	newTs, err := s.timeslot_repo.GetTimeslotByID(booking.TimeslotID)
	if err != nil {
		return err
	}
	if newTs != nil && newTs.Booking_id != nil {
		*newTs.Booking_id = booking.ID
		newTimeslot, err := newTs.ToTimeslot()
		if err != nil {
			return err
		}
		if err = s.timeslot_repo.UpdateTimeslot(newTimeslot); err != nil {
			return err
		}
	}

	return nil
}

func (s *bookingService) DeleteBooking(id int) error {
	err := s.repo.DeleteBooking(id)
	if err != nil {
		return err
	}

	return nil
}
