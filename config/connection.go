package config

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"learnapirest/model"
	"log"
)

var DB *gorm.DB

func InitDB() {
	dsn := "root:root@tcp(localhost:3306)/tokosepatu?parseTime=true"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal konek ke database:", err)
	}

	db.AutoMigrate(&model.Sepatu{})
	DB = db
}
