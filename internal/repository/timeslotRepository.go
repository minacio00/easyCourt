package repository

import (
	"log"

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
	// Find the booking first to get the timeslot_id
	var booking model.Booking
	if err := r.db.First(&booking, bookingId).Error; err != nil {
		return nil, err
	}

	// Then get the timeslot using the timeslot_id from the booking
	return r.GetTimeslotByID(booking.TimeslotID)
}

func (r *timeslotRepository) GetTimeslotsByCourt(courtID int, weekDay string) ([]model.ReadTimeslot, error) {
	// First fetch all timeslots matching the criteria
	var timeslots []model.ReadTimeslot
	query := r.db.Table("timeslots").
		Select(`
	timeslots.*, 
    courts.id AS court_id,
    courts.name,
    courts.created_at,
    courts.updated_at
`).
		Joins("LEFT JOIN courts ON timeslots.court_id = courts.id").
		Where("timeslots.court_id = ?", courtID)

	if weekDay != "" {
		query = query.Where("timeslots.day = ?", weekDay)
	} else {
		query = query.Order("timeslots.day ASC")
	}

	if err := query.Order("timeslots.start_time ASC").Scan(&timeslots).Error; err != nil {
		return nil, err
	}

	// For each timeslot, we need to fetch the associated booking (if any)
	for i := range timeslots {
		var booking model.ReadBooking
		err := r.db.Table("bookings").
			Select("bookings.*, users.*").
			Joins("LEFT JOIN users ON bookings.user_id = users.id").
			Where("bookings.timeslot_id = ?", timeslots[i].ID).
			Scan(&booking).Error

		// Only assign the booking if one was found
		if err == nil && booking.ID != 0 {
			timeslots[i].Booking = &booking
		}
	}

	return timeslots, nil
}

func (r *timeslotRepository) CreateTimeslot(timeslot *model.Timeslot) error {
	return r.db.Select("court_id", "day", "start_time", "end_time", "is_active").Create(timeslot).Error
}

func (r *timeslotRepository) GetTimeslotByID(id int) (*model.ReadTimeslot, error) {
	var timeslot model.ReadTimeslot
	err := r.db.Table("timeslots").
		Select(`
            timeslots.id,
            timeslots.court_id,
            timeslots.day,
            timeslots.start_time,
            timeslots.end_time,
            timeslots.is_active,
            courts.id AS court__id,
            courts.name AS court__name
        `).
		Joins("LEFT JOIN courts ON timeslots.court_id = courts.id").
		Where("timeslots.id = ?", id).
		Scan(&timeslot).Error
	if err != nil {
		log.Println("error in gettimeslotbyid repo ", err)
		return nil, err
	}
	if timeslot.ID == 0 {
		return nil, nil
	}

	// Then get any associated booking
	var booking model.ReadBooking
	err = r.db.Table("bookings").
		Select("bookings.*, users.*").
		Joins("LEFT JOIN users ON bookings.user_id = users.id").
		Where("bookings.timeslot_id = ?", id).
		Scan(&booking).Error

	// Only assign the booking if one was found
	if err == nil && booking.ID != 0 {
		timeslot.Booking = &booking
	}

	return &timeslot, nil
}

func (r *timeslotRepository) GetAllTimeslots() ([]model.ReadTimeslot, error) {
	var timeslots []model.Timeslot
	var readSlots []model.ReadTimeslot

	log.Println("got here")

	err := r.db.Preload("Court").
		Preload("Booking.User").
		Find(&timeslots).Error
	if err != nil {
		log.Println("error after db load ", err)
		return nil, err
	}

	if len(timeslots) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	for _, ts := range timeslots {
		readSlot := model.ReadTimeslot{
			ID:        ts.ID,
			CourtID:   ts.CourtID,
			Court:     ts.Court,
			Day:       ts.Day,
			StartTime: ts.StartTime.Format("15:04:05"),
			EndTime:   ts.EndTime.Format("15:04:05"),
			IsActive:  ts.IsActive,
		}

		if ts.Booking != nil {
			readSlot.Booking = &model.ReadBooking{
				ID:              ts.Booking.ID,
				UserID:          ts.Booking.UserID,
				Opponent:        ts.Booking.Opponent,
				OpponentPartner: ts.Booking.OpponentPartner,
				BookingDate:     ts.Booking.BookingDate,
				IsSinglesGame:   ts.Booking.IsSinglesGame,
				TimeslotID:      ts.Booking.TimeslotID,

				User: model.User{
					ID:    ts.Booking.User.ID,
					Name:  ts.Booking.User.Name,
					Email: ts.Booking.User.Email,
				},
			}
		}

		readSlots = append(readSlots, readSlot)
	}

	return readSlots, nil
}

func (r *timeslotRepository) UpdateTimeslot(timeslot *model.Timeslot) error {
	return r.db.Debug().Model(&model.Timeslot{}).Where("id = ?", timeslot.ID).Updates(map[string]interface{}{
		"court_id":   timeslot.CourtID,
		"day":        timeslot.Day,
		"start_time": timeslot.StartTime,
		"end_time":   timeslot.EndTime,
		"is_active":  timeslot.IsActive,
	}).Error
}

func (r *timeslotRepository) DeleteTimeslot(id int) error {
	return r.db.Delete(&model.Timeslot{}, id).Error
}

func (r *timeslotRepository) GetActiveTimeslots() ([]model.Timeslot, error) {
	// First get the active timeslots
	var simpleTimeslots []model.Timeslot
	if err := r.db.Where("is_active = ?", true).Find(&simpleTimeslots).Error; err != nil {
		return nil, err
	}

	// For each timeslot, fetch courts and any associated bookings
	for i := range simpleTimeslots {
		// Get the court
		if simpleTimeslots[i].CourtID != nil {
			var court model.Court
			r.db.First(&court, *simpleTimeslots[i].CourtID)
			simpleTimeslots[i].Court = court
		}

		// Get any associated booking
		var booking model.Booking
		if err := r.db.Where("timeslot_id = ?", simpleTimeslots[i].ID).First(&booking).Error; err == nil {
			simpleTimeslots[i].Booking = &booking

			// Get the user for the booking if needed
			if booking.UserID != 0 {
				var user model.User
				r.db.First(&user, booking.UserID)
				booking.User = user
			}
		}
	}

	return simpleTimeslots, nil
}
