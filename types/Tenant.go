package types

import (
	"errors"

	"gorm.io/gorm"
)

type Tenant struct {
	gorm.Model
	Name        string `gorm:"not null" json:"name"`
	Email       string `gorm:"not null" json:"email"`
	Password    string `gorm:"not null" json:"password"`
	TrialPeriod bool   `gorm:"not null" json:"periodo de teste"`
	Clube       Clube
}

func (t *Tenant) Validate() error {
	// Check if the name is empty
	if t.Name == "" {
		return errors.New("name cannot be empty")
	}
	if t.Email == "" {
		return errors.New("name cannot be empty")
	}
	if t.Password == "" {
		return errors.New("name cannot be empty")
	}
	// No error found, return nil
	return nil
}
