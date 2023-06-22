package main

import (
	"fmt"

	conf "github.com/mp3ismp3/gocrypto/config"
	"github.com/mp3ismp3/gocrypto/db/mysql"
	"github.com/mp3ismp3/gocrypto/models"
	"github.com/mp3ismp3/gocrypto/servers/grpcserver"
)

func main() {
	fmt.Println("InitConfig")
	conf.InitConfig()
	mysql.InitDB()
	fmt.Println("InitDB")
	models.ModelInit()
	models.InitEngineList()
	fmt.Println("InitInitEngineList")
	grpcserver.Init()
	fmt.Println("Initserver")
}
