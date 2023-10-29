package types

import (
	"errors"
	"net/mail"
	"strings"

	"gorm.io/gorm"
)

type Tenant struct {
	gorm.Model
	DeletedAt   interface{} `gorm:"-"` // Ignore DeletedAt
	Name        string      `gorm:"not null" json:"name"`
	Email       string      `gorm:"not null" json:"email"`
	Password    string      `gorm:"not null" json:"password"`
	TrialPeriod bool        `gorm:"not null" json:"periodo_teste"`
	Clube       Clube       `gorm:"null constraint:OnUpdate:CASCADE OnDelete:CASCADE;"`
}

func (t *Tenant) Validate() error {
	var validationErrors []string

	// Check if the name is empty
	if t.Name == "" {
		validationErrors = append(validationErrors, "Nome não pode ser vazio=")
	}

	// Check if the email is empty or invalid
	if t.Email == "" {
		validationErrors = append(validationErrors, "Email não pode ser vazio")
	} else if _, err := mail.ParseAddress(t.Email); err != nil {
		validationErrors = append(validationErrors, "Email inválido")
	}

	// Check if the password is empty
	if t.Password == "" {
		validationErrors = append(validationErrors, "Senha não pode ser vazio")
	}

	// Check if there are any validation errors
	if len(validationErrors) > 0 {
		return errors.New(strings.Join(validationErrors, ", "))
	}

	// No error found, return nil
	return nil
}
