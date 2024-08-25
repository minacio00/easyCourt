package service

import (
	"github.com/minacio00/easyCourt/internal/model"
	"github.com/minacio00/easyCourt/internal/repository"
)

type CourtService interface {
	CreateCourt(court *model.Court) error
	GetAllCourts() ([]model.Court, error)
	GetCourtByID(id int) (*model.Court, error)
	UpdateCourt(court *model.Court) error
	DeleteCourt(id int) error
}

type courtService struct {
	repo repository.CourtRepository
}

// NewCourtService creates a new instance of CourtService
func NewCourtService(repo repository.CourtRepository) CourtService {
	return &courtService{repo}
}

// CreateCourt calls the repository to create a new court
func (s *courtService) CreateCourt(court *model.Court) error {
	return s.repo.CreateCourt(court)
}

// GetAllCourts retrieves all courts by calling the repository
func (s *courtService) GetAllCourts() ([]model.Court, error) {
	return s.repo.GetAllCourts()
}

// GetCourtByID retrieves a court by its ID by calling the repository
func (s *courtService) GetCourtByID(id int) (*model.Court, error) {
	return s.repo.GetCourtByID(id)
}

// UpdateCourt updates an existing court by calling the repository
func (s *courtService) UpdateCourt(court *model.Court) error {
	return s.repo.UpdateCourt(court)
}

// DeleteCourt deletes a court by its ID by calling the repository
func (s *courtService) DeleteCourt(id int) error {
	return s.repo.DeleteCourt(id)
}
