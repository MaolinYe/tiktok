package db

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
)

func init() {
	// 连接数据库
	username := "root"
	password := ""
	host := "localhost"
	port := 3306
	dbname := "tiktok"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Print(err)
		panic("failed to connect database")
	}
	// 自动迁移
	db.AutoMigrate(&User{}, &Video{})
	log.Println("database connected")
}
