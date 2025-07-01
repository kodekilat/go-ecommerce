package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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
	products, err := h.ProductRepo.GetAllProducts()
	if err != nil {
		log.Printf("Gagal mengambil produk untuk admin: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// 2. Buat struct untuk menampung data yang akan dikirim ke template
	pageData := struct {
		Products []models.Product
	}{
		Products: products,
	}

	view.Render(w, "admin_products.page.html", pageData)
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

// ShowEditProductForm menampilkan form untuk mengedit produk yang ada
func (h *AdminHandler) ShowEditProductForm(w http.ResponseWriter, r *http.Request) {
	// 1. Ambil ID dari URL
	productIDStr := chi.URLParam(r, "productID")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		http.Error(w, "ID Produk tidak valid", http.StatusBadRequest)
		return
	}

	// 2. Ambil data produk dari database
	product, err := h.ProductRepo.GetProductByID(productID)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.NotFound(w, r)
		} else {
			log.Printf("Gagal mengambil produk: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// 3. Render template dengan data produk
	view.Render(w, "edit_product.page.html", product)
}

// HandleUpdateProduct memproses update produk
func (h *AdminHandler) HandleUpdateProduct(w http.ResponseWriter, r *http.Request) {
	// 1. Ambil ID dari URL
	productIDStr := chi.URLParam(r, "productID")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		http.Error(w, "ID Produk tidak valid", http.StatusBadRequest)
		return
	}

	// Ambil data produk yang ada untuk mendapatkan URL gambar lama
	existingProduct, err := h.ProductRepo.GetProductByID(productID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	r.ParseMultipartForm(10 << 20) // 10 MB

	// 2. Ambil data dari form
	name := r.PostFormValue("name")
	description := r.PostFormValue("description")
	price, _ := strconv.ParseFloat(r.PostFormValue("price"), 64)
	stock, _ := strconv.Atoi(r.PostFormValue("stock"))

	imageURL := existingProduct.ImageURL // Default ke URL gambar yang lama

	// 3. Cek apakah ada gambar baru yang diunggah
	file, header, err := r.FormFile("image")
	if err == nil { // Jika ada file baru
		defer file.Close()

		// Hapus gambar lama dari Minio (opsional, praktik yang baik)
		// ...

		// Unggah gambar baru
		bucketName := "products"
		objectName := uuid.New().String() + filepath.Ext(header.Filename)
		_, err = storage.MinioClient.PutObject(context.Background(), bucketName, objectName, file, header.Size, minio.PutObjectOptions{ContentType: header.Header.Get("Content-Type")})
		if err != nil {
			http.Error(w, "Gagal mengunggah gambar baru", http.StatusInternalServerError)
			return
		}
		imageURL = fmt.Sprintf("http://127.0.0.1:9000/%s/%s", bucketName, objectName)
	}

	// 4. Buat model produk yang diperbarui
	updatedProduct := &models.Product{
		ID:          productID,
		Name:        name,
		Description: description,
		Price:       price,
		Stock:       stock,
		ImageURL:    imageURL,
	}

	// 5. Panggil repository untuk meng-update database
	err = h.ProductRepo.UpdateProduct(updatedProduct)
	if err != nil {
		log.Printf("Gagal mengupdate produk: %v", err)
		http.Error(w, "Gagal menyimpan perubahan", http.StatusInternalServerError)
		return
	}

	// 6. Redirect kembali ke halaman admin
	http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
}

func (h *AdminHandler) HandleDeleteProduct(w http.ResponseWriter, r *http.Request) {
	// 1. Ambil ID dari URL
	productIDStr := chi.URLParam(r, "productID")
	productID, err := uuid.Parse(productIDStr)
	if err != nil {
		http.Error(w, "ID Produk tidak valid", http.StatusBadRequest)
		return
	}

	// 2. Panggil repository untuk menghapus produk
	// (Opsional: hapus juga gambar dari Minio di sini)
	err = h.ProductRepo.DeleteProduct(productID)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.NotFound(w, r)
		} else {
			log.Printf("Gagal menghapus produk: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	// 3. Redirect kembali ke halaman admin setelah berhasil
	http.Redirect(w, r, "/admin/products", http.StatusSeeOther)
}
