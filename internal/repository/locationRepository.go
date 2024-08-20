package repository

import (
	"github.com/minacio00/easyCourt/internal/model"
	"gorm.io/gorm"
)

type LocationRepository interface {
	CreateLocation(location *model.Location) error
	GetAllLocations() ([]model.Location, error)
	UpdateLocation(location *model.Location) error
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

func (l *locationRepository) UpdateLocation(location *model.Location) error {
	return l.db.Save(location).Error
}

func (l *locationRepository) DeleteLocation(id uint) error {
	return l.db.Delete(&model.Location{}, id).Error
}
