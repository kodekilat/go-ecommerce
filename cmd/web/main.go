package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	// Membuat multiplexer (router) baru
	mux := http.NewServeMux()

	// Mendaftarkan handler untuk root URL ("/")
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Selamat Datang di Toko Online Go!")
	})

	// Pesan bahwa server akan dimulai
	log.Println("Memulai server di http://localhost:8080")

	// Menjalankan server di port 8080
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Fatal(err)
	}
}
