package model

type Location struct {
	ID   int    `gorm:"primaryKey"`
	Name string `json:"location_name" gorm:"unique;not null"`
}
