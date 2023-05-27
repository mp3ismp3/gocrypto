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
		models.AddBuyOrder(*order)
	} else {
		trade := engine.LimitSell(*order)
		fmt.Println(trade)
		models.AddBuyOrder(*order)
		// AddSellOrder(symbol string, side int32, price float32, qty float32, ordertype int32)

	}
	return order.ID
}
