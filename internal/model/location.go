package model

import "errors"

type Location struct {
	ID       int    `gorm:"primaryKey"`
	Name     string `json:"location_name" gorm:"unique;not null"`
	ImageUrl string `json:"image_url"`
}

func (l *Location) Validate() error {
	if l.Name == "" {
		return errors.New("location name is required")
	}
	return nil
}

type CreateLocation struct {
	Name string `json:"location_name"`
}

func (l *CreateLocation) Validate() error {
	if l.Name == "" {
		return errors.New("location name is required")
	}
	return nil
}
