package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"-" gorm:"unique; not null"`
	IsAdmin  bool   `json:"-"`
}
