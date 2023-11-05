package types

import (
	"errors"

	"gorm.io/gorm"
)

type Clube struct {
	gorm.Model
	DeletedAt interface{} `gorm:"-"` // Ignore DeletedAt
	TenantID  uint        `gorm:"not null" json:"tenant_id"`
	ClubName  string      `gorm:"not null" json:"club_name"`
	Quadras   []Quadra    `json:"quadras"`
	Clientes  []Cliente   `json:"clientes"`
}

func (c *Clube) Validate() error {
	// Check if the tenant ID is zero
	if c.TenantID == 0 {
		return errors.New("tenant ID cannot be zero")
	}
	// Check if the club name is empty
	if c.ClubName == "" {
		return errors.New("club name cannot be empty")
	}
	// No error found, return nil
	return nil
}
