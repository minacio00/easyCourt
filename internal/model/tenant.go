package model

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Tenant struct {
	gorm.Model
	Name      string `json:"name" gorm:"unique;not null"`
	Email     string `json:"email" gorm:"unique; not nul"`
	Password  string `json:"-"`
	FreeTrial time.Time
}
