package models

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type OrderNode struct {
	gorm.Model
	Symbol    string
	Price     float32
	Side      int32
	Quantity  float32
	Value     float32
	OrderType int32
	EngineID  uint
}

func SetOrder(symbol string, side int32, price float32, qty float32, ordertype int32) *OrderNode {
	node := &OrderNode{}
	node.SetCreateTime()
	node.SetSymboyType(symbol)
	node.SetPrice(price)
	node.SetSide(side)
	node.SetQuantity(qty)
	node.SetValue(qty, price)
	node.SetOrderType(ordertype)

	return node

}

func GetOrder(id uint64) (OrderNode, *gorm.DB) {
	var order OrderNode
	db := DB.Where("ID=?", id).Find(&order)
	return order, db
}

func (node *OrderNode) AddOrder() OrderNode {
	result := DB.Create(&node)
	if result.Error != nil {
		fmt.Printf("AddOrder Faild, error: %s", result.Error)
	}
	return *node
}

func DeleteOrder(id uint64) OrderNode {
	var order OrderNode
	result := DB.Where("ID = ?", id).Delete(&order)
	if result.Error != nil {
		fmt.Printf("DeleteOrder Faild, error: %s", result.Error)
	}
	return order
}

func (node *OrderNode) SetCreateTime() {
	node.CreatedAt = time.Now()
}

func (node *OrderNode) SetSymboyType(symbol string) {
	node.Symbol = symbol
}

func (node *OrderNode) SetPrice(price float32) {
	node.Price = price
}

func (node *OrderNode) SetSide(side int32) {
	node.Side = side
}

func (node *OrderNode) SetQuantity(qty float32) {
	node.Quantity = qty
}

func (node *OrderNode) SetValue(qty float32, price float32) {
	node.Value = qty * price
}

func (node *OrderNode) SetOrderType(ordertype int32) {
	node.OrderType = int32(ordertype)
}
