package controllers

import (
	"log"

	"github.com/mp3ismp3/gocrypto/models"
)

var (
	EngineList map[string]models.Engine
)

func AddEngine(symbol string, openPrice float32) (int32, string) {
	EngineList = models.GetEngineList()
	var Engine models.Engine
	_, ok := EngineList[symbol]
	if ok {
		log.Println("Engine is already exist.")
		return 404, "Engine is already exist."
	}
	Engine.AddEngine(symbol, openPrice)
	return 300, "succesful create"
}

func DeleteEngine(symbol string) (int32, string) {
	EngineList = models.GetEngineList()
	_, ok := EngineList[symbol]
	if ok {
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
