package service

import (
	"github.com/minacio00/easyCourt/internal/model"
	"github.com/minacio00/easyCourt/internal/repository"
)

type LocationService interface {
	CreateLocation(location *model.Location) error
	GetAllLocations() ([]model.Location, error)
	UpdateLocation(location *model.Location) error
	DeleteLocation(id uint) error
}

type locationService struct {
	repo repository.LocationRepository
}

func NewLocationService(repo repository.LocationRepository) LocationService {
	return &locationService{repo}
}

func (s *locationService) CreateLocation(location *model.Location) error {
	// Add any additional logic before creating a location, if necessary
	return s.repo.CreateLocation(location)
}

func (s *locationService) GetAllLocations() ([]model.Location, error) {
	return s.repo.GetAllLocations()
}

func (s *locationService) UpdateLocation(location *model.Location) error {
	// Add any additional logic before updating a location, if necessary
	return s.repo.UpdateLocation(location)
}

func (s *locationService) DeleteLocation(id uint) error {
	// Add any additional logic before deleting a location, if necessary
	return s.repo.DeleteLocation(id)
}
