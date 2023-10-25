package types

import "gorm.io/gorm"

type Quadra struct {
	gorm.Model
	ClubeID    uint   `gorm:"not null" json:"club_id"`
	TipoQuadra string `gorm:"not null" json:"tipo_quadra"`
	Clube      Clube  `json:"-"`
	Reserva    Reserva
}
