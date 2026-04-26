package main

import (
	"learnapirest/internal/config"
	"learnapirest/internal/modules/account"
	"learnapirest/internal/modules/cart"
	"learnapirest/internal/modules/order"
	"learnapirest/internal/modules/product"
	"learnapirest/internal/modules/transaction"
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

	paymentRepo := transaction.NewPaymentRepository(db)
	paymentServ := transaction.NewPaymentService(paymentRepo, orderRepo, appConfig.StripeSecretKey, appConfig.StripeWebhookSecret)

	r := router.SetupRouter(sepatuServ, accountServ, cartServ, orderServ, paymentServ, appConfig)
	r.Run(":8080")
}
