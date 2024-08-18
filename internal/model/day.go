package model

type Day struct {
	ID   int    `gorm:"primaryKey"`
	Name string `json:"weekday" gorm:"unique; not null"`
}
