package main

import (

	// Asumsi path import-mu menyesuaikan
	"learnapirest/internal/config"
	"learnapirest/internal/modules/account"
	"learnapirest/internal/modules/product"
	"learnapirest/router"
)

func main() {
	appConfig := config.LoadConfig()

	db := config.InitDB()

	productRepo := product.NewProductRepository(db)
	sepatuServ := product.NewProductService(productRepo)

	accountRepo := account.NewAccountRepository(db)
	accountServ := account.NewAccountService(accountRepo, appConfig)

	r := router.SetupRouter(sepatuServ, accountServ, appConfig)
	r.Run(":8080")
}
