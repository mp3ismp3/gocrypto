package controllers

import (
	"fmt"

	"github.com/mp3ismp3/gocrypto/models"
)

func AddOrder(symbol string, side int32, price float32, qty float32, ordertype int32) uint {
	engine := GetEngine(symbol)
	order := models.SetOrder(symbol, side, price, qty, ordertype)
	if side == 0 {
		trade := engine.LimitBuy(*order)
		fmt.Println(trade)
	} else {
		trade := engine.LimitSell(*order)
		fmt.Println(trade)
	}
	orderNode := order.AddOrder()
	return orderNode.ID
}

func DeleteOrder(symbol string, OrderId uint64) uint {
	order := models.DeleteOrder(OrderId)
	return order.ID
}

func GetOrder(OrderId uint64) models.OrderNode {
	order, _ := models.GetOrder(uint64(OrderId))
	return order
}

func UpdateOrder(id uint64, price float32, qty float32) uint {
	order := GetOrder(id)
	order.Price = price
	order.Quantity = qty
	return order.ID
}
