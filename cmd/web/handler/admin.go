package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/google/uuid"
	"github.com/kodekilat/go-ecommerce/cmd/web/view"
	"github.com/kodekilat/go-ecommerce/internal/models"
	"github.com/kodekilat/go-ecommerce/internal/repository"
	"github.com/kodekilat/go-ecommerce/internal/storage"
	"github.com/minio/minio-go/v7"
)

type AdminHandler struct {
	ProductRepo *repository.ProductRepository
}

// ShowAdminProducts menampilkan halaman admin produk
func (h *AdminHandler) ShowAdminProducts(w http.ResponseWriter, r *http.Request) {
	// (Nanti di sini kita akan ambil daftar produk dan menampilkannya)
	view.Render(w, "admin_products.page.html", nil)
}

// HandleAddProduct memproses penambahan produk baru
func (h *AdminHandler) HandleAddProduct(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20) // 10 MB

	// 1. Ambil data dari form dan konversi tipe datanya
	name := r.PostFormValue("name")
	description := r.PostFormValue("description")
	priceStr := r.PostFormValue("price")
	stockStr := r.PostFormValue("stock")

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		http.Error(w, "Harga tidak valid", http.StatusBadRequest)
		return
	}
	stock, err := strconv.Atoi(stockStr)
	if err != nil {
		http.Error(w, "Stok tidak valid", http.StatusBadRequest)
		return
	}

	// 2. Proses upload file
	var imageURL string
	file, header, err := r.FormFile("image")
	if err == nil {
		defer file.Close()
		bucketName := "products"
		objectName := uuid.New().String() + filepath.Ext(header.Filename)

		_, err = storage.MinioClient.PutObject(context.Background(), bucketName, objectName, file, header.Size, minio.PutObjectOptions{
			ContentType: header.Header.Get("Content-Type"),
		})
		if err != nil {
			log.Printf("Gagal mengunggah gambar: %v", err)
			http.Error(w, "Gagal mengunggah gambar", http.StatusInternalServerError)
			return
		}
		imageURL = fmt.Sprintf("http://127.0.0.1:9000/%s/%s", bucketName, objectName)
	}

	// 3. Buat model produk dan simpan ke database
	newProduct := &models.Product{
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		ImageURL:    imageURL,
	}

	err = h.ProductRepo.CreateProduct(newProduct)
	if err != nil {
		log.Printf("Gagal menyimpan produk: %v", err)
		http.Error(w, "Gagal menyimpan produk", http.StatusInternalServerError)
		return
	}

	log.Printf("Produk baru berhasil disimpan dengan ID: %s", newProduct.ID)
	http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
}
