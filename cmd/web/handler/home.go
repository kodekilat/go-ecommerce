package handler

import (
	"net/http"

	"github.com/kodekilat/go-ecommerce/cmd/web/view" // Ganti dengan path modul Anda
)

// Definisikan struktur untuk data produk dummy
type Product struct {
	Name  string
	Price int
}

func ShowHomePage(w http.ResponseWriter, r *http.Request) {
	// Siapkan data yang akan dikirim ke template
	pageData := struct {
		PageTitle      string
		WelcomeMessage string
		Products       []Product
	}{
		PageTitle:      "Selamat Datang di Toko Kami!",
		WelcomeMessage: "Temukan produk terbaik hanya di sini.",
		Products: []Product{
			{Name: "Buku Go Keren", Price: 150000},
			{Name: "Stiker Gopher", Price: 25000},
			{Name: "Mug PostgreSQL", Price: 75000},
		},
	}

	// Panggil fungsi render
	view.Render(w, "home.page.html", pageData)
}
