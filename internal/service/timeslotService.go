package service

import (
	"fmt"

	"github.com/minacio00/easyCourt/internal/model"
	"github.com/minacio00/easyCourt/internal/model/utils"
	"github.com/minacio00/easyCourt/internal/repository"
)

type TimeslotService interface {
	CreateTimeslot(timeslot *model.Timeslot) error
	GetTimeslotByID(id int) (*model.Timeslot, error)
	GetAllTimeslots() ([]model.ReadTimeslot, error)
	UpdateTimeslot(timeslot *model.Timeslot) error
	DeleteTimeslot(id int) error
	GetActiveTimeslots() ([]model.Timeslot, error)
}

type timeslotService struct {
	repo      repository.TimeslotRepository
	courtRepo repository.CourtRepository
}

func NewTimeslotService(repo repository.TimeslotRepository, court_repo repository.CourtRepository) TimeslotService {
	return &timeslotService{repo, court_repo}
}

func (s *timeslotService) CreateTimeslot(timeslot *model.Timeslot) error {
	if err := timeslot.Validate(); err != nil {
		return err
	}
	var err error = nil
	timeslot.Day, err = utils.MapWeekDay(string(timeslot.Day))
	if err != nil {
		return err
	}

	//check if the court exist in the db
	_, err = s.courtRepo.GetCourtByID(*timeslot.CourtID)
	if err != nil {
		return fmt.Errorf("court with id %d does not exist: %w", *timeslot.CourtID, err)
	}
	return s.repo.CreateTimeslot(timeslot)
}

func (s *timeslotService) GetTimeslotByID(id int) (*model.Timeslot, error) {
	return s.repo.GetTimeslotByID(id)
}

func (s *timeslotService) GetAllTimeslots() ([]model.ReadTimeslot, error) {
	readSlots, err := s.repo.GetAllTimeslots()
	if err != nil {
		return nil, err
	}

	// var timeslots []model.Timeslot
	// for _, rt := range readSlots {
	// 	ts, err := rt.ToTimeslot()
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	timeslots = append(timeslots, *ts)
	// }
	return readSlots, nil
}

func (s *timeslotService) UpdateTimeslot(timeslot *model.Timeslot) error {
	// Add any additional logic before updating a timeslot, if necessary
	return s.repo.UpdateTimeslot(timeslot)
}

func (s *timeslotService) DeleteTimeslot(id int) error {
	// Add any additional logic before deleting a timeslot, if necessary
	return s.repo.DeleteTimeslot(id)
}

func (s *timeslotService) GetActiveTimeslots() ([]model.Timeslot, error) {
	return s.repo.GetActiveTimeslots()
}
