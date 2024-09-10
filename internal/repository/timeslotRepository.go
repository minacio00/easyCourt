package repository

import (
	"github.com/minacio00/easyCourt/internal/model"
	"gorm.io/gorm"
)

type TimeslotRepository interface {
	CreateTimeslot(timeslot *model.Timeslot) error
	GetTimeslotByID(id int) (*model.Timeslot, error)
	GetAllTimeslots() ([]model.ReadTimeslot, error)
	UpdateTimeslot(timeslot *model.Timeslot) error
	DeleteTimeslot(id int) error
	GetActiveTimeslots() ([]model.Timeslot, error)
}

type timeslotRepository struct {
	db *gorm.DB
}

func NewTimeslotRepository(db *gorm.DB) TimeslotRepository {
	return &timeslotRepository{db}
}

func (r *timeslotRepository) CreateTimeslot(timeslot *model.Timeslot) error {
	println(&timeslot.CourtID)
	return r.db.Create(timeslot).Error
}

func (r *timeslotRepository) GetTimeslotByID(id int) (*model.Timeslot, error) {
	var timeslot model.Timeslot
	if err := r.db.Preload("Court").First(&timeslot, id).Error; err != nil {
		return nil, err
	}
	return &timeslot, nil
}

func (r *timeslotRepository) GetAllTimeslots() ([]model.ReadTimeslot, error) {
	var timeslots []model.ReadTimeslot
	if err := r.db.Preload("Court").Model(&model.Timeslot{}).Find(&timeslots).Error; err != nil {
		return nil, err
	}
	return timeslots, nil
}

func (r *timeslotRepository) UpdateTimeslot(timeslot *model.Timeslot) error {
	return r.db.Save(timeslot).Error
}

func (r *timeslotRepository) DeleteTimeslot(id int) error {
	return r.db.Delete(&model.Timeslot{}, id).Error
}

func (r *timeslotRepository) GetActiveTimeslots() ([]model.Timeslot, error) {
	var timeslots []model.Timeslot
	if err := r.db.Preload("Court").Where("is_active = ?", true).Find(&timeslots).Error; err != nil {
		return nil, err
	}
	return timeslots, nil
}
