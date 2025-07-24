package main

import (
	"learnapirest/config"
	"learnapirest/router"
)

func main() {
	config.InitDB()
	r := router.SetupRouter()
	r.Run(":8080")
}
