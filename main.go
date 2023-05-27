package main

import (
	"fmt"

	conf "github.com/mp3ismp3/gocrypto/config"
	mysql "github.com/mp3ismp3/gocrypto/db/mysql"
)

func main() {
	fmt.Println("InitConfig")
	conf.InitConfig()
	mysql.InitDB()
	fmt.Println("InitDB")
	// models.InitEngineList()
	// fmt.Println("InitInitEngineList")
	// server.Init()
	// fmt.Println("Initserver")
}
