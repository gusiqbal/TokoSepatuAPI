package main

import (
	"log"
	// Asumsi path import-mu menyesuaikan
	"learnapirest/internal/config"
	"learnapirest/internal/modules/account"
	"learnapirest/internal/modules/product"
	"learnapirest/router"
)

func main() {
	// 1. Load Konfigurasi
	appConfig := config.LoadConfig()

	// 2. Inisiasi Database (Return instance, BUKAN via global config.DB)
	db := config.InitDB()

	// 3. Jalankan AutoMigrate di sini
	err := db.AutoMigrate(
		&account.User{},
		&product.Product{}, // Sesuaikan dengan nama modelmu
	)
	if err != nil {
		log.Fatal("Gagal migrasi database: ", err)
	}

	// 4. Dependency Injection - Modul Sepatu
	productRepo := product.NewProductRepository(db)
	sepatuServ := product.NewProductService(productRepo)

	// 5. Dependency Injection - Modul Account
	// Perhatikan: Tidak ada lagi tanda `*` saat mem-passing repo dan config
	accountRepo := account.NewAccountRepository(db)
	accountServ := account.NewAccountService(accountRepo, appConfig)

	// 6. Setup Router & Jalankan App
	r := router.SetupRouter(sepatuServ, accountServ, appConfig)
	r.Run(":8080")
}
