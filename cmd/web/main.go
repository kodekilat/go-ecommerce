package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kodekilat/go-ecommerce/internal/database" // Ganti dengan path modul Anda
)

func main() {
	// Membuat koneksi database
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Tidak dapat terhubung ke database: %v", err)
	}
	defer db.Close() // Pastikan koneksi ditutup saat aplikasi berhenti

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Selamat Datang di Toko Online Go! Koneksi database berhasil.")
	})

	log.Println("Memulai server di http://localhost:8080")

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
