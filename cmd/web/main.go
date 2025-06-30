package main

import (
	"log"
	"net/http"

	"github.com/kodekilat/go-ecommerce/cmd/web/router" // Ganti dengan path modul Anda
	"github.com/kodekilat/go-ecommerce/internal/database"
	"github.com/kodekilat/go-ecommerce/internal/repository"
)

func main() {
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Tidak dapat terhubung ke database: %v", err)
	}
	defer db.Close()

	userRepo := &repository.UserRepository{DB: db}

	appRouter := router.New(userRepo)

	log.Println("Memulai server di http://localhost:8080")

	err = http.ListenAndServe(":8080", appRouter)
	if err != nil {
		log.Fatal(err)
	}
}
