package controllers

import (
	"github.com/mp3ismp3/gocrypto/models"
)

var (
	EngineList map[string]models.Engine
)

func AddEngine(symbol string, openPrice float32) (int32, string) {
	var Engine models.Engine
	Engine.Symbol = symbol
	Engine.Price = openPrice
	err := Engine.AddEngine()
	if err != nil {
		return 404, err.Error()
	}
	return 300, "succesful create"
}

func DeleteEngine(symbol string) (int32, string) {
	err := models.DeleteEngine(symbol)
	if err != nil {
		return 404, err.Error()
	}
	return 300, "succesful delete"
}

func GetEngine(symbol string) *models.Engine {
	engine, err := models.GetEngine(symbol)
	if err != nil {
		return nil
	}
	return engine
}
