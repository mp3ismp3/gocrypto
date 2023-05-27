package models

import (
	"gorm.io/gorm"
)

type OrderLog struct {
	gorm.Model
	OrderID string
}
