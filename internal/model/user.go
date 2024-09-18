package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password" gorm:"unique; not null"`
	IsAdmin  bool   `json:"-"`
}
type UserResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"-"`
	IsAdmin  bool   `json:"isAdmin"`
}

func (u *User) MapUserToResponse() *UserResponse {
	return &UserResponse{
		ID:      u.ID,
		Name:    u.Name,
		Phone:   u.Phone,
		IsAdmin: u.IsAdmin,
	}
}
