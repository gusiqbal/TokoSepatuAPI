package main

import (
	"learnapirest/internal/config"
	"learnapirest/internal/modules/account"
	"learnapirest/internal/modules/cart"
	"learnapirest/internal/modules/order"
	"learnapirest/internal/modules/product"
	"log"
	// import module lain yang punya tabel
)

func main() {
	db := config.InitDB()
	var getAllModels []any

	getAllModels = append(getAllModels, product.GetProduct())
	getAllModels = append(getAllModels, order.GetOrder())
	getAllModels = append(getAllModels, cart.GetCart())
	getAllModels = append(getAllModels, account.GetUser())
	err := db.AutoMigrate(getAllModels...)

	if err != nil {
		log.Fatal("Migration failed: ", err)
	}

	log.Println("Database migration successfully generated!")
}
