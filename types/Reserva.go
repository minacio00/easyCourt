package types

import (
	"time"

	"gorm.io/gorm"
)

type Reserva struct {
	gorm.Model
	ClienteID     uint      `gorm:"null" json:"client_id"`
	QuadraID      uint      `gorm:"not null" json:"court_id"`
	InicioHorario time.Time `gorm:"type:timestamp;default:current_timestamp" json:"inicio_horario"`
	FimHorario    time.Time `gorm:"type:timestamp;default:current_timestamp" json:"fim_horario"`
}
