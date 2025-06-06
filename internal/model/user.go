package model

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Password  string    `json:"-"`
	IsAdmin   bool      `json:"isAdmin"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
	Bookings  []Booking `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE;" json:"-"`
}

func (u *User) Validate() error {
	// Validate Name
	if strings.TrimSpace(u.Name) == "" {
		return errors.New("nome não deve ser vazio")
	}
	if len(u.Name) < 2 || len(u.Name) > 50 {
		return errors.New("nome deve ter entre 2 e 50 carateres")
	}

	// Validate Phone
	if err := validatePhone(u.Phone); err != nil {
		return err
	}
	if err := validateEmail(u.Email); err != nil {
		return err
	}

	// Validate Password (only if it's being set or updated)
	if u.Password == "" {
		return errors.New("senha não deve ser vazia")
	}
	if u.Password != "" {
		if len(u.Password) < 5 {
			return errors.New("senha deve ter ao menos 6 caracteres")
		}
	}
	return nil
}

func validateEmail(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-/&()]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	if !emailRegex.MatchString(email) {
		return errors.New("email deve ter o formato 'exemplo@dominio.com'")
	}

	return nil
}

func validatePhone(phone string) error {
	phoneRegex := regexp.MustCompile(`^(\d{2})(\d{8,9})$`)

	matches := phoneRegex.FindStringSubmatch(phone)
	if matches == nil {
		return errors.New("telefone deve ter o formato'DDDnúmero', ex., '62995032121'")
	}

	ddd := matches[1]
	number := matches[2]

	// Validate DDD (always 2 digits)
	if len(ddd) != 2 {
		return errors.New("DDD deve conter 2 dígitos")
	}

	if len(number) != 8 && len(number) != 9 {
		return errors.New("número deve conter ou 9 dígitos após o ddd")
	}

	return nil
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	LastName string `json:"lastName"`
	Email    string `json:"email"`
	Password string `json:"-"`
	IsAdmin  bool   `json:"isAdmin"`
}
type CreateUser struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	LastName string `json:"lastName"`
	Email    string `json:"email"`
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
		LastName: u.LastName,
		Email:    u.Email,
	}
}

func (u *User) MapUserToResponse() *UserResponse {
	return &UserResponse{
		ID:       u.ID,
		Name:     u.Name,
		Phone:    u.Phone,
		IsAdmin:  u.IsAdmin,
		LastName: u.LastName,
		Email:    u.Email,
	}
}
