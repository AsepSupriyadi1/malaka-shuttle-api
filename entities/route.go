package entities

import (
	"gorm.io/gorm"
)

type Route struct {
	gorm.Model
	OriginCity      string     `gorm:"size:100;not null"`
	DestinationCity string     `gorm:"size:100;not null"`
	Schedules       []Schedule `gorm:"foreignKey:RouteID"`
}
