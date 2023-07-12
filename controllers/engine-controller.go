package controllers

import (
	"fmt"
	"log"

	"github.com/mp3ismp3/gocrypto/models"
)

var (
	EngineList map[string]models.Engine
)

func AddEngine(symbol string, openPrice float32) (int32, string) {
	EngineList = models.GetEngineList()
	var Engine models.Engine
	Engine.Symbol = symbol
	Engine.Price = openPrice
	_, ok := EngineList[symbol]
	if ok {
		log.Println("Engine is already exist.")
		return 300, "Engine is already exist."
	}
	// engine := models.GetEngine(symbol)
	// if engine != nil {
	// 	return 300, "Engine is already exist."
	// }
	Engine.AddEngine()
	return 300, "succesful create"
}

func DeleteEngine(symbol string) (int32, string) {
	EngineList = models.GetEngineList()
	fmt.Println(EngineList)
	_, ok := EngineList[symbol]
	if !ok {
		log.Println("Engine is not exist.")
		return 404, "Engine is not exist."
	}
	models.DeleteEngine(symbol)
	return 300, "succesful create"
}

func GetEngine(symbol string) models.Engine {
	EngineList = models.GetEngineList()
	engine := EngineList[symbol]
	return engine
}
