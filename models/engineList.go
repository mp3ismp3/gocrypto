package models

var EngineList map[string]Engine

func InitEngineList() {
	engineList := map[string]Engine{
		"ETH": {
			"ds",
			"ETHUSDT",
			1800,
			make([]OrderNode, 0, 100),
			make([]OrderNode, 0, 100),
		},
		"BTC": {
			"ddw",
			"BTCUSDT",
			26000,
			make([]OrderNode, 0, 100),
			make([]OrderNode, 0, 100),
		},
	}
	//redis
	for k, v := range engineList {
		err := cache.HSet(ctx, "EngineList", k, v).Err()
		if err != nil {
			panic(err)
		}
	}

	EngineList = engineList
}

func GetEngineList() map[string]Engine {
	return EngineList
}
