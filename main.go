package main

import (
	"learnapirest/config"
	sepatuRepository "learnapirest/repository"
	"learnapirest/router"
	sepatuService "learnapirest/service"
)

func main() {
	config.InitDB()
	sepatuRepo := sepatuRepository.NewSepatuRepo(config.DB)
	sepatuServ := sepatuService.NewSepatuService(sepatuRepo)
	r := router.SetupRouter(sepatuServ)
	r.Run(":8080")
}
