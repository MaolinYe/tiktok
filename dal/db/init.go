package db

import (
	"fmt"
	"log"
	"tiktok/pkg/viper"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

var (
	db     *gorm.DB
	config = viper.InitConfig()
)

func getDSN(role string) string {
	username := config.GetString(fmt.Sprintf("mysql.%s.username", role))
	password := config.GetString(fmt.Sprintf("mysql.%s.password", role))
	host := config.GetString(fmt.Sprintf("mysql.%s.host", role))
	port := config.GetInt(fmt.Sprintf("mysql.%s.port", role))
	dbname := config.GetString(fmt.Sprintf("mysql.%s.dbname", role))
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname)
}

func init() {
	// 连接数据库
	dsn_source := getDSN("source")
	dsn_replica := getDSN("replica")
	var err error // go的:=会初始化局部变量
	db, err = gorm.Open(mysql.Open(dsn_source), &gorm.Config{})
	if err != nil {
		log.Print(err)
		panic("failed to connect database")
	}
	db.Use(dbresolver.Register(dbresolver.Config{
		Sources:           []gorm.Dialector{mysql.Open(dsn_source)},
		Replicas:          []gorm.Dialector{mysql.Open(dsn_replica)},
		Policy:            dbresolver.RandomPolicy{},
		TraceResolverMode: false, // 不记录数据库切换信息
	}))
	// 自动迁移
	if err = db.AutoMigrate(&User{}, &Video{}, &Favor{}); err != nil {
		log.Println(err)
	}
	log.Println("database connected")
}
