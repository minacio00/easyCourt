package repository

import (
	"github.com/minacio00/easyCourt/internal/model"
	"gorm.io/gorm"
)

type TimeslotRepository interface {
	CreateTimeslot(timeslot *model.Timeslot) error
	GetTimeslotByID(id int) (*model.ReadTimeslot, error)
	GetAllTimeslots() ([]model.ReadTimeslot, error)
	UpdateTimeslot(timeslot *model.Timeslot) error
	DeleteTimeslot(id int) error
	GetActiveTimeslots() ([]model.Timeslot, error)
	GetTimeslotsByCourt(courtID int, weekDay string) ([]model.ReadTimeslot, error)
	GetTimeslotByBookingId(bookingId uint) (*model.ReadTimeslot, error)
}

type timeslotRepository struct {
	db *gorm.DB
}

func NewTimeslotRepository(db *gorm.DB) TimeslotRepository {
	return &timeslotRepository{db}
}

func (r *timeslotRepository) GetTimeslotByBookingId(bookingId uint) (*model.ReadTimeslot, error) {
	var timeslot *model.ReadTimeslot
	err := r.db.Model(&model.Timeslot{}).Preload("Court").Where("booking_id = ?", bookingId).First(&timeslot).Error
	if err != nil {
		return nil, err
	}
	return timeslot, nil
}

func (r *timeslotRepository) GetTimeslotsByCourt(courtID int, weekDay string) ([]model.ReadTimeslot, error) {
	var timeslots []model.ReadTimeslot

	query := r.db.Model(&model.Timeslot{}).
		Preload("Booking.User").
		Where("court_id = ?", courtID)

	if weekDay != "" {
		query = query.Where("day = ?", weekDay)
	} else {
		query = query.Order("day ASC")
	}

	if err := query.Order("start_time ASC").Find(&timeslots).Error; err != nil {
		return nil, err
	}
	return timeslots, nil

}
func (r *timeslotRepository) CreateTimeslot(timeslot *model.Timeslot) error {
	return r.db.Select("court_id", "day", "start_time", "end_time", "is_active").Create(timeslot).Error
}

func (r *timeslotRepository) GetTimeslotByID(id int) (*model.ReadTimeslot, error) {
	var timeslot *model.ReadTimeslot
	if err := r.db.Model(&model.Timeslot{}).Preload("Court").Preload("Booking.User").First(&timeslot, id).Error; err != nil {
		return nil, err
	}
	// s, err := timeslot.ToTimeslot()
	// if err != nil {
	// 	return nil, err
	// }
	return timeslot, nil
}
func (r *timeslotRepository) GetAllTimeslots() ([]model.ReadTimeslot, error) {
	var timeslots []model.ReadTimeslot
	if err := r.db.Preload("Court").Preload("Booking.User").Model(&model.Timeslot{}).Find(&timeslots).Error; err != nil {
		return nil, err
	}
	return timeslots, nil
}

func (r *timeslotRepository) UpdateTimeslot(timeslot *model.Timeslot) error {
	return r.db.Model(&model.Timeslot{}).Where("id = ?", timeslot.ID).Updates(map[string]interface{}{
		"court_id":   timeslot.CourtID,
		"day":        timeslot.Day,
		"start_time": timeslot.StartTime,
		"end_time":   timeslot.EndTime,
		"is_active":  timeslot.IsActive,
		"booking_id": timeslot.Booking_id,
	}).Error
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
