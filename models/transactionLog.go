package models

import (
	"gorm.io/gorm"
)

type TransactionLog struct {
	gorm.Model
	TakerOrderId uint
	MakerOrderId uint
	Quantity     float32
	Price        float32
}
