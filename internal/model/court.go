package model

import "gorm.io/gorm"

type Court struct {
	gorm.Model
	Name       string   `json:"name" gorm:"not null"`
	LocationID int      `json:"location_id"`
	Location   Location `gorm:"foreignKey:LocationID"`
}
