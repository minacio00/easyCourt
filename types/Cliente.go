package types

import "gorm.io/gorm"

type Cliente struct {
	gorm.Model
	DeletedAt       interface{} `gorm:"-"` // Ignore DeletedAt
	ClubeID         uint        `json:"-"`
	Nome            string      `gorm:"not null" json:"nome"`
	Email           string      `gorm:"not null" json:"email"`
	Senha           string      `gorm:"not null" json:"senha"`
	Telefone        string      `json:"telefone" gorm:"not null"`
	PagamentoStatus bool        `json:"pagamento" gorm:"null"`
	Reserva         Reserva
}
