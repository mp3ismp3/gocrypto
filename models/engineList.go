package models

var (
	EngineList map[string]Engine
)

func InitEngineList() {

	EngineList = make(map[string]Engine)
	// ETH := Engine{
	// 	Symbol: "ETHUSDT",
	// 	Price:  1800,
	// }
	// BTC := Engine{
	// 	Symbol: "BTCUSDT",
	// 	Price:  26000,
	// }

	// ETH.AddEngine()
	// BTC.AddEngine()

	//redis
	// for k, v := range engineList {
	// 	err := cache.HSet(ctx, "EngineList", k, v).Err()
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }
}

func GetEngineList() map[string]Engine {
	return EngineList
}
