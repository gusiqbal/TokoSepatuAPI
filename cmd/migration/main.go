package main

import (
	"learnapirest/internal/config"
	"learnapirest/internal/modules/account"
	"learnapirest/internal/modules/cart"
	"learnapirest/internal/modules/order"
	"learnapirest/internal/modules/product"
	"learnapirest/internal/modules/transaction"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Could not load .env file:", err)
	}

	db := config.InitDB()
	var getAllModels []any

	getAllModels = append(getAllModels, product.GetProduct()...)
	getAllModels = append(getAllModels, order.GetOrder()...)
	getAllModels = append(getAllModels, cart.GetCart()...)
	getAllModels = append(getAllModels, account.GetUser()...)
	getAllModels = append(getAllModels, transaction.GetPayment()...)
	err = db.AutoMigrate(getAllModels...)

	if err != nil {
		log.Fatal("Migration failed: ", err)
	}

	log.Println("Database migration successfully generated!")
}
