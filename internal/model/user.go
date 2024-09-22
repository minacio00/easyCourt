package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Phone     string    `json:"phone"`
	Password  string    `json:"-"`
	IsAdmin   bool      `json:"isAdmin"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}
type UserResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"-"`
	IsAdmin  bool   `json:"isAdmin"`
}
type CreateUser struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"-"`
}

func (u *CreateUser) MapCreateToUser() *User {
	return &User{
		ID:       u.ID,
		Name:     u.Name,
		Phone:    u.Phone,
		IsAdmin:  u.IsAdmin,
		Password: u.Password,
	}
}

func (u *User) MapUserToResponse() *UserResponse {
	return &UserResponse{
		ID:      u.ID,
		Name:    u.Name,
		Phone:   u.Phone,
		IsAdmin: u.IsAdmin,
	}
}
