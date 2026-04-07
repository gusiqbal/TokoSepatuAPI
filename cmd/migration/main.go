package main

import (
	"learnapirest/internal/config"
	"learnapirest/internal/modules/account"
	"learnapirest/internal/modules/product"
	"log"
	// import module lain yang punya tabel
)

func main() {
	db := config.InitDB()

	err := db.AutoMigrate(
		&account.User{},
		&product.Product{},
	)

	if err != nil {
		log.Fatal("Migration failed: ", err)
	}

	log.Println("Database migration successfully generated!")
}
