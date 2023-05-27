package mysql

import (
	"log"
	"strings"

	conf "github.com/mp3ismp3/gocrypto/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

var DB *gorm.DB

func InitDB() {
	var err error
	mConfig := conf.Config.MySql["default"]
	pathRead := strings.Join([]string{mConfig.UserName, ":", mConfig.Password, "@tcp(", mConfig.DbHost, ":", mConfig.DbPort, ")/", mConfig.DbName, "?charset=" + mConfig.Charset + "&parseTime=true"}, "")
	pathWrite := strings.Join([]string{mConfig.UserName, ":", mConfig.Password, "@tcp(", mConfig.DbHost, ":", mConfig.DbPort, ")/", mConfig.DbName, "?charset=" + mConfig.Charset + "&parseTime=true"}, "")

	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       pathRead, // DSN data source name
		DefaultStringSize:         256,      // string 类型字段的默认长度
		DisableDatetimePrecision:  true,     // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,     // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,     // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,    // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})
	// db, err := gorm.Open(mysql.Open(pathRead), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	DB.Logger = logger.Default.LogMode(logger.Info)
	// 設置數據庫連接池參數
	sqlDB, _ := DB.DB()
	// 設置數據庫連接池最大數量
	sqlDB.SetMaxOpenConns(100)
	// 連接池最大允许的空閒連接數
	sqlDB.SetMaxIdleConns(20)

	_ = DB.Use(dbresolver.Register(dbresolver.Config{
		// use `db2` as sources, `db3`, `db4` as replicas
		Sources:  []gorm.Dialector{mysql.Open(pathWrite)},
		Replicas: []gorm.Dialector{mysql.Open(pathRead), mysql.Open(pathRead)},
		// sources/replicas load balancing policy
		Policy: dbresolver.RandomPolicy{},
		// print sources/replicas mode in logger
		TraceResolverMode: true,
	}))

	DB = DB.Set("gorm:table_options", "charset=utf8mb4")
	log.Println("Database connection successful.")

}

func GetDB() *gorm.DB {
	return DB
}
