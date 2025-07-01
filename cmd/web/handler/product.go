package handler

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/kodekilat/go-ecommerce/cmd/web/view"
	"github.com/kodekilat/go-ecommerce/internal/repository"
)

type ProductHandler struct {
	ProductRepo *repository.ProductRepository
}

func (h *ProductHandler) ShowProductDetail(w http.ResponseWriter, r *http.Request) {
	productIDStr := chi.URLParam(r, "productID")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		http.Error(w, "ID Produk tidak valid", http.StatusBadRequest)
		return
	}

	product, err := h.ProductRepo.GetProductByID(productID)
	if err != nil {
		log.Printf("Produk tidak ditemukan: %v", err)
		http.NotFound(w, r)
		return
	}

	view.Render(w, "product_detail.page.html", product)
}
