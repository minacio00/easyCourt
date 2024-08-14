package model

import (
	"time"

	"gorm.io/gorm"
)

type Tenant struct {
	gorm.Model
	Name       string `json:"name" gorm:"unique;not null"`
	Email      string `json:"email" gorm:"unique; not nul"`
	Password   string `json:"-"`
	SchemaName string `json:"schema_name"`
	FreeTrial  time.Time
}
