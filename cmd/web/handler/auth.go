package handler

import (
	"log"
	"net/http"

	"github.com/kodekilat/go-ecommerce/cmd/web/view" // Ganti dengan path modul Anda
	"github.com/kodekilat/go-ecommerce/internal/models"
	"github.com/kodekilat/go-ecommerce/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UserRepo *repository.UserRepository
}

// ShowRegistrationForm menampilkan halaman registrasi
func (h *AuthHandler) ShowRegistrationForm(w http.ResponseWriter, r *http.Request) {
	view.Render(w, "register.page.html", nil)
}

// HandleRegistration memproses data dari form registrasi
func (h *AuthHandler) HandleRegistration(w http.ResponseWriter, r *http.Request) {
	// 1. Parse form data
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Gagal memproses form", http.StatusBadRequest)
		return
	}

	fullName := r.PostForm.Get("full_name")
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	// 2. Validasi sederhana
	if fullName == "" || email == "" || password == "" {
		http.Error(w, "Semua field harus diisi", http.StatusBadRequest)
		return
	}

	// 3. Hash password dengan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Gagal hash password: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Buat model user baru
	newUser := &models.User{
		FullName:     fullName,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	err = h.UserRepo.CreateUser(newUser)
	if err != nil {
		// (Tambahkan pengecekan error duplikat email di sini nanti)
		log.Printf("Gagal menyimpan user: %v", err)
		http.Error(w, "Gagal mendaftarkan pengguna", http.StatusInternalServerError)
		return
	}

	log.Printf("User baru berhasil disimpan dengan ID: %s", newUser.ID)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
