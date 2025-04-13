package model

import (
	"errors"

	"gorm.io/gorm"
)

type Court struct {
	gorm.Model
	Name       string    `json:"name" gorm:"not null"`
	LocationID int       `json:"location_id"`
	Location   *Location `json:"location,omitempty" gorm:"foreignKey:LocationID;constraint:OnDelete:CASCADE"`
}

func (c *Court) Validate() error {
	if c.Name == "" {
		return errors.New("name is required")
	}
	if c.LocationID <= 0 {
		return errors.New("location_id must be a positive integer")
	}
	return nil
}
