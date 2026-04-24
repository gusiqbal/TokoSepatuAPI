package main

import (
	"learnapirest/internal/config"
	"learnapirest/internal/modules/account"
	"learnapirest/internal/modules/cart"
	"learnapirest/internal/modules/order"
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

	cartRepo := cart.NewCartRepository(db)
	cartServ := cart.NewCartService(cartRepo)

	orderRepo := order.NewOrderRepository(db)
	orderServ := order.NewOrderService(orderRepo, cartRepo)

	r := router.SetupRouter(sepatuServ, accountServ, cartServ, orderServ, appConfig)
	r.Run(":8080")
}
