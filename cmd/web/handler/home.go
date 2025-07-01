package handler

import (
	"log"
	"net/http"

	"github.com/kodekilat/go-ecommerce/cmd/web/view"
	"github.com/kodekilat/go-ecommerce/internal/models"
	"github.com/kodekilat/go-ecommerce/internal/repository"
)

type HomeHandler struct {
	ProductRepo *repository.ProductRepository
}

func (h *HomeHandler) ShowHomePage(w http.ResponseWriter, r *http.Request) {
	// Ambil semua produk dari repository
	products, err := h.ProductRepo.GetAllProducts()
	if err != nil {
		log.Printf("Gagal mengambil produk: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	pageData := struct {
		Products []models.Product
	}{
		Products: products,
	}

	view.Render(w, "home.page.html", pageData)
}
