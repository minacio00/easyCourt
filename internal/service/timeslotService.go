package service

import (
	"fmt"
	"log"

	"github.com/minacio00/easyCourt/internal/model"
	"github.com/minacio00/easyCourt/internal/model/utils"
	"github.com/minacio00/easyCourt/internal/repository"
)

type TimeslotService interface {
	CreateTimeslot(timeslot *model.Timeslot) error
	GetTimeslotByID(id int) (*model.ReadTimeslot, error)
	GetAllTimeslots() ([]model.ReadTimeslot, error)
	UpdateTimeslot(timeslot *model.Timeslot) error
	DeleteTimeslot(id int) error
	GetActiveTimeslots() ([]model.Timeslot, error)
	GetTimeslotsByCourt(courtID int, weekDay string) ([]model.ReadTimeslot, error)
}

type timeslotService struct {
	repo      repository.TimeslotRepository
	courtRepo repository.CourtRepository
}

func NewTimeslotService(repo repository.TimeslotRepository, court_repo repository.CourtRepository) TimeslotService {
	return &timeslotService{repo, court_repo}
}

func (s *timeslotService) GetTimeslotsByCourt(courtID int, weekDay string) ([]model.ReadTimeslot, error) {
	// First, check if the court exists
	_, err := s.courtRepo.GetCourtByID(courtID)
	if err != nil {
		return nil, err
	}
	// todo: parse weekDay string
	// If the court exists, get its timeslots
	return s.repo.GetTimeslotsByCourt(courtID, weekDay)
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

	// check if the court exist in the db
	_, err = s.courtRepo.GetCourtByID(*timeslot.CourtID)
	if err != nil {
		return fmt.Errorf("court with id %d does not exist: %w", *timeslot.CourtID, err)
	}
	return s.repo.CreateTimeslot(timeslot)
}

func (s *timeslotService) GetTimeslotByID(id int) (*model.ReadTimeslot, error) {
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
	// Validate the timeslot
	if err := timeslot.Validate(); err != nil {
		return err
	}

	// Check if the timeslot exists
	existingTimeslot, err := s.repo.GetTimeslotByID(timeslot.ID)
	if err != nil {
		return fmt.Errorf("timeslot with id %d returned error: %w", timeslot.ID, err)
	}
	ts, err := existingTimeslot.ToTimeslot()
	if err != nil {
		return err
	}

	// Update only the fields that are allowed to be updated
	ts.CourtID = timeslot.CourtID
	ts.Day = timeslot.Day
	ts.StartTime = timeslot.StartTime
	ts.EndTime = timeslot.EndTime
	ts.IsActive = timeslot.IsActive
	err = s.repo.UpdateTimeslot(ts)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Println("error during update timeslot: ", err)
		return err
	}

	return nil
}

func (s *timeslotService) DeleteTimeslot(id int) error {
	// Add any additional logic before deleting a timeslot, if necessary
	return s.repo.DeleteTimeslot(id)
}

func (s *timeslotService) GetActiveTimeslots() ([]model.Timeslot, error) {
	return s.repo.GetActiveTimeslots()
}
