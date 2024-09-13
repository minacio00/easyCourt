package repository

import (
	"github.com/minacio00/easyCourt/internal/model"
	"gorm.io/gorm"
)

type LocationRepository interface {
	CreateLocation(location *model.Location) error
	GetAllLocations() ([]model.Location, error)
	UpdateLocation(location *model.Location) error
	GetAllLocationCourts(location_id uint) ([]model.Court, error)
	DeleteLocation(id uint) error
}

type locationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) LocationRepository {
	return &locationRepository{db}
}

func (l *locationRepository) CreateLocation(location *model.Location) error {
	return l.db.Create(location).Error
}

func (l *locationRepository) GetAllLocations() ([]model.Location, error) {
	var location []model.Location
	if err := l.db.Find(&location).Error; err != nil {
		return nil, err
	}
	return location, nil
}

func (l *locationRepository) GetAllLocationCourts(location_id uint) ([]model.Court, error) {
	var courts []model.Court
	if err := l.db.Where("location_id = ?", location_id).Find(&courts).Error; err != nil {
		return nil, err
	}
	return courts, nil
}

func (l *locationRepository) UpdateLocation(location *model.Location) error {
	existingLocation := &model.Location{}
	if err := l.db.First(existingLocation, location.ID).Error; err != nil {
		return err // return error if location is not found
	}
	location.ID = existingLocation.ID

	// Update the location with the new data
	return l.db.Save(location).Error
}

func (l *locationRepository) DeleteLocation(id uint) error {
	return l.db.Delete(&model.Location{}, id).Error
}
