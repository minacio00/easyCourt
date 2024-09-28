package service

import (
	"errors"

	"github.com/minacio00/easyCourt/internal/model"
	"github.com/minacio00/easyCourt/internal/repository"
)

type LocationService interface {
	CreateLocation(location *model.Location) error
	GetAllLocations() ([]model.Location, error)
	UpdateLocation(location *model.Location) error
	DeleteLocation(id uint) error
	GetLocationById(id uint) (*model.Location, error)
}

type locationService struct {
	repo repository.LocationRepository
}

func NewLocationService(repo repository.LocationRepository) LocationService {
	return &locationService{repo}
}

func (s *locationService) GetLocationById(id uint) (*model.Location, error) {
	return s.repo.GetLocationById(id)
}
func (s *locationService) CreateLocation(location *model.Location) error {
	if location.Name != "" {
		println(location.Name)
		return s.repo.CreateLocation(location)
	}
	return errors.New("location_name cannot be empty")
}

func (s *locationService) GetAllLocations() ([]model.Location, error) {
	return s.repo.GetAllLocations()
}

func (s *locationService) UpdateLocation(location *model.Location) error {
	return s.repo.UpdateLocation(location)
}

func (s *locationService) DeleteLocation(id uint) error {
	return s.repo.DeleteLocation(id)
}
