package repository

import (
	"github.com/minacio00/easyCourt/internal/model"
	"gorm.io/gorm"
)

type CourtRepository interface {
	CreateCourt(court *model.Court) error
	GetAllCourts() ([]model.Court, error)
	GetCourtByID(id int) (*model.Court, error)
	UpdateCourt(court *model.Court) error
	DeleteCourt(id int) error
}

type courtRepository struct {
	db *gorm.DB
}

// NewCourtRepository creates a new instance of CourtRepository
func NewCourtRepository(db *gorm.DB) CourtRepository {
	return &courtRepository{db}
}

// CreateCourt adds a new court to the database
func (r *courtRepository) CreateCourt(court *model.Court) error {
	return r.db.Create(court).Error
}

// GetAllCourts retrieves all courts from the database
func (r *courtRepository) GetAllCourts() ([]model.Court, error) {
	var courts []model.Court
	err := r.db.Find(&courts).Error
	return courts, err
}

// GetCourtByID retrieves a court by its ID from the database
func (r *courtRepository) GetCourtByID(id int) (*model.Court, error) {
	var court model.Court
	err := r.db.First(&court, id).Error
	if err != nil {
		return nil, err
	}
	return &court, nil
}

// UpdateCourt updates an existing court in the database
func (r *courtRepository) UpdateCourt(court *model.Court) error {
	// Find the court by ID to ensure it exists
	existingCourt := &model.Court{}
	if err := r.db.First(existingCourt, court.ID).Error; err != nil {
		return err // return error if court is not found
	}

	// Assign the found court's ID to the court object to ensure the correct record is updated
	court.ID = existingCourt.ID

	// Update the court with the new data
	return r.db.Save(court).Error
}

// DeleteCourt deletes a court by its ID from the database
func (r *courtRepository) DeleteCourt(id int) error {
	return r.db.Delete(&model.Court{}, id).Error
}
