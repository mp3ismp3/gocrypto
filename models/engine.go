package models

import (
	"encoding/json"
	"fmt"

	"context"

	"github.com/mp3ismp3/gocrypto/db/mysql"
	"github.com/mp3ismp3/gocrypto/db/redis"
	"gorm.io/gorm"
)

type Engine struct {
	ID         uint
	Symbol     string `gorm:"unique"`
	Price      float32
	BuyOrders  []OrderNode
	SellOrders []OrderNode
}

var (
	DB    *gorm.DB
	cache = redis.InitRedis()
	ctx   = context.Background()
	// cache = redis.NewRedisClient()
)

func ModelInit() {
	DB = mysql.GetDB()
	err := DB.AutoMigrate(&Engine{}, &OrderNode{})
	if err != nil {
		panic(err)
	} //其作用主要是刷新数据库中的表格，使其保持最新，即让数据库之前存储的记录的表格字段和程序中最新使用的表格字段保持一致（只增不减）
}

func GetEngine(symbol string) (*Engine, error) {
	engine := Engine{}
	//從redis緩存中查詢數據
	result, err := cache.Get(ctx, symbol).Result()
	if err == nil {
		err = json.Unmarshal([]byte(result), &engine)
		if err != nil {
			return nil, err
		}
		return &engine, nil
	}
	//從mysql中獲取數據
	if err := DB.Where("symbol = ?", symbol).First(&engine).Error; err != nil {
		return nil, err
	}
	jsonData, err := json.Marshal(engine)
	if err != nil {
		return nil, err
	}
	//將數據存儲到Redis緩存中
	if err := cache.Set(ctx, symbol, jsonData, 0).Err(); err != nil {
		return nil, err
	}
	return &engine, nil
}

func (engine *Engine) AddEngine() error {
	//新增mysql數據
	err := DB.Create(&engine).Error
	if err != nil {
		return err
	}
	//新增Redis緩存數據
	jsonData, err := json.Marshal(engine)
	if err != nil {
		return err
	}
	fmt.Println("cache:", cache)
	if err := cache.Set(ctx, engine.Symbol, jsonData, 0).Err(); err != nil {
		return err
	}
	// EngineList[engine.Symbol] = *engine
	return nil
}

func DeleteEngine(symbol string) error {

	//新的会话，不会继承之前的查詢條件。
	session := DB.Session(&gorm.Session{})
	//從mysql中刪除symbol，delete(EngineList, symbol)
	if err := session.Where("symbol = ?", symbol).Delete(&Engine{}).Error; err != nil {
		fmt.Println("刪除Mysql數據時出錯:", err)
		return err
	}

	//刪除該緩存數據
	result, err := cache.Del(ctx, symbol).Result()
	if err != nil {
		fmt.Println("刪除Redis數據時出錯:", err)
		return err
	}
	fmt.Println("已成功刪除", result, "緩存數據")
	return nil
}

// func GetAll(db *gorm.DB) ([]engine, error) {
// 	var orders []engine
// 	err := db.Model(&engine{}).Preload("CreditCards").Find(&orders).Error
// 	return orders, err
// }

func (engine *Engine) SetSymbol(symbol string) {
	engine.Symbol = symbol
}

func (engine *Engine) SetPrice(price float32) {
	engine.Price = price
}

func (engine *Engine) addBuyOrder(order OrderNode) {
	l := 0
	r := len(engine.BuyOrders) - 1
	var mid int = (l + r) / 2

	for l <= r {
		if order.Price == engine.BuyOrders[mid].Price { // 會有時間順序問題
			slice2 := engine.BuyOrders[mid:]
			engine.BuyOrders = append(engine.BuyOrders, order)
			copy(engine.BuyOrders[mid+1:], slice2)
			engine.BuyOrders[mid+1] = order
			return
		} else if order.Price < engine.BuyOrders[mid].Price {
			l = mid + 1
		} else {
			r = mid - 1
		}
	}

	if r == mid-1 {
		slice2 := engine.BuyOrders[r+1:]
		engine.BuyOrders = append(engine.BuyOrders, order)
		copy(engine.BuyOrders[r+2:], slice2)
		engine.BuyOrders[r+1] = order
		return
	} else {
		slice2 := engine.BuyOrders[l:]
		engine.BuyOrders = append(engine.BuyOrders, order)
		copy(engine.BuyOrders[l+1:], slice2)
		engine.BuyOrders[l] = order
	}

}

func (engine *Engine) addSellOrder(order OrderNode) {
	l := 0
	r := len(engine.SellOrders) - 1
	var mid int = (l + r) / 2
	for l <= r {
		if engine.SellOrders[mid].Price == order.Price {
			slice2 := engine.SellOrders[mid:]
			engine.SellOrders = append(engine.SellOrders, order)
			copy(engine.SellOrders[mid+1:], slice2)
			engine.SellOrders[mid+1] = order
			return
		} else if engine.SellOrders[mid].Price < order.Price {
			l = mid + 1
		} else {
			r = mid - 1
		}
	}

	if r == mid-1 {
		slice2 := engine.SellOrders[r+1:]
		engine.SellOrders = append(engine.SellOrders, order)
		copy(engine.SellOrders[r+2:], slice2)
		engine.SellOrders[r+1] = order
		return
	} else {
		slice2 := engine.SellOrders[l:]
		engine.SellOrders = append(engine.SellOrders, order)
		copy(engine.SellOrders[l+1:], slice2)
		engine.SellOrders[l] = order
	}

}

func (engine *Engine) removeSellOrder(index int) {
	engine.SellOrders = append(engine.SellOrders[:index], engine.SellOrders[index+1:]...)
}

func (engine *Engine) removeBuyOrder(index int) {
	engine.BuyOrders = append(engine.BuyOrders[:index], engine.BuyOrders[index+1:]...)
}

const DEFAULT_EXPIRATION = 3600

func (engine *Engine) BuyOrSell(order OrderNode) []TransactionLog {

	if order.Side == 0 {
		// redis
		trade := engine.LimitBuy(order)
		// a, _ := json.Marshal(engine)
		// _, error := cache.Set(ctx, engine.Symbol, string(a), DEFAULT_EXPIRATION*time.Second).Result()
		// fmt.Println(error)
		return trade
	}
	// redis
	trade := engine.LimitSell(order)
	// a, _ := json.Marshal(engine)
	// _, error := cache.Set(ctx, engine.Symbol, string(a), DEFAULT_EXPIRATION*time.Second).Result()
	// fmt.Println(error)
	return trade
}

// 連續競價

func (engine *Engine) BuyContinuousAuction(last_Order *OrderNode, Order *OrderNode) {
	if Order.Price > last_Order.Price {
		engine.Price = last_Order.Price
	} else {
		engine.Price = Order.Price
	}
	fmt.Println("當前價格:", engine.Price)
}

func (engine *Engine) SellContinuousAuction(last_Order *OrderNode, Order *OrderNode) {
	if Order.Price < last_Order.Price {
		engine.Price = last_Order.Price
	} else {
		engine.Price = Order.Price
	}
	fmt.Println("當前價格:", engine.Price)
}

func (engine *Engine) LimitBuy(order OrderNode) []TransactionLog {
	trades := make([]TransactionLog, 0, 1)
	n := len(engine.SellOrders)
	if n != 0 && engine.SellOrders[0].Price <= order.Price {
		//買入的價格>=最低賣出的價格
		for i := 0; i < n; i++ {
			sellOrder := engine.SellOrders[i]
			if sellOrder.Price > order.Price {
				break
			}
			//匹配完整的數量
			if sellOrder.Quantity >= order.Quantity {
				sellOrder.Quantity -= order.Quantity
				sellOrder.Value -= (order.Quantity * order.Price)

				if sellOrder.Quantity == 0 {
					engine.removeSellOrder(i)
				}
				engine.BuyContinuousAuction(&sellOrder, &order)
				trades = append(trades, TransactionLog{gorm.Model{}, order.Model.ID, sellOrder.Model.ID, order.Quantity, engine.Price})
				return trades
			}
			//匹配部分數量
			if sellOrder.Quantity < order.Quantity {
				order.Quantity -= sellOrder.Quantity
				order.Value -= (sellOrder.Quantity * order.Price)

				engine.removeSellOrder(i)
				engine.BuyContinuousAuction(&sellOrder, &order)
				trades = append(trades, TransactionLog{gorm.Model{}, order.Model.ID, sellOrder.Model.ID, order.Quantity, engine.Price})
				continue
			}

		}

	}
	// 買入的價格 < 最低賣出的價格
	// 剩下的訂單加入BuyOrder
	engine.addBuyOrder(order)
	return trades
}

func (engine *Engine) LimitSell(order OrderNode) []TransactionLog {
	trades := make([]TransactionLog, 0, 1)
	n := len(engine.BuyOrders)
	if n != 0 && engine.BuyOrders[0].Price >= order.Price {
		for i := 0; i < n; i++ {
			buyOrder := engine.BuyOrders[i]
			if buyOrder.Price < order.Price {
				break
			}
			// 完全匹配
			if buyOrder.Quantity >= order.Quantity {
				buyOrder.Quantity -= order.Quantity
				buyOrder.Value -= (order.Price * order.Quantity)

				if buyOrder.Quantity == 0 {
					engine.removeBuyOrder(i)
				}
				engine.SellContinuousAuction(&buyOrder, &order)
				trades = append(trades, TransactionLog{gorm.Model{}, order.Model.ID, buyOrder.Model.ID, order.Quantity, engine.Price})
				return trades
			}
			// 部分匹配
			if buyOrder.Quantity < order.Quantity {
				order.Quantity -= buyOrder.Quantity
				order.Value -= (order.Price * buyOrder.Quantity)

				engine.removeBuyOrder(i)

				engine.SellContinuousAuction(&buyOrder, &order)
				trades = append(trades, TransactionLog{gorm.Model{}, order.Model.ID, buyOrder.Model.ID, order.Quantity, engine.Price})
				continue
			}
		}
	}
	engine.addSellOrder(order)
	return trades
}
