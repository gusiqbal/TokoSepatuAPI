package main

import (
	"learnapirest/config"
	repo "learnapirest/repository"
	"learnapirest/router"
	service "learnapirest/service"
)

func main() {
	config.InitDB()
	appConfig := config.LoadConfig()

	sepatuRepo := repo.NewSepatuRepo(config.DB)
	sepatuServ := service.NewSepatuService(sepatuRepo)

	accounRepo := repo.NewAccountRepository(config.DB)
	accountServ := service.NewAccountService(*accounRepo, *appConfig)

	r := router.SetupRouter(sepatuServ, accountServ, appConfig)
	r.Run(":8080")
}
