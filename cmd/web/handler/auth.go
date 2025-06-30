package handler

import (
	"log"
	"net/http"

	"github.com/kodekilat/go-ecommerce/cmd/web/session"
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

// ShowLoginForm menampilkan halaman login
func (h *AuthHandler) ShowLoginForm(w http.ResponseWriter, r *http.Request) {
	view.Render(w, "login.page.html", nil)
}

// HandleLogin memproses data dari form login
func (h *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Gagal memproses form", http.StatusBadRequest)
		return
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	// 1. Cari pengguna berdasarkan email
	user, err := h.UserRepo.GetUserByEmail(email)
	if err != nil {
		// Jika tidak ada baris yang ditemukan, email tidak terdaftar.
		// Kita berikan pesan error yang umum untuk keamanan.
		log.Printf("Email tidak ditemukan: %s, error: %v", email, err)
		http.Error(w, "Email atau password salah", http.StatusUnauthorized)
		return
	}

	// 2. Verifikasi password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		// Jika password tidak cocok, bcrypt akan mengembalikan error.
		log.Printf("Password salah untuk email: %s", email)
		http.Error(w, "Email atau password salah", http.StatusUnauthorized)
		return
	}

	// Jika berhasil sampai sini, email dan password valid!
	log.Printf("Pengguna berhasil login: %s (ID: %s)", user.Email, user.ID)

	// 1. Dapatkan sesi atau buat yang baru
	sess, _ := session.Store.Get(r, "auth-session")

	// 2. Set nilai di dalam sesi
	sess.Values["user_id"] = user.ID
	sess.Values["user_email"] = user.Email

	// 3. Simpan sesi (ini akan mengirim cookie ke browser)
	err = sess.Save(r, w)
	if err != nil {
		log.Printf("Gagal menyimpan sesi: %v", err)
		http.Error(w, "Gagal login", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *AuthHandler) HandleLogout(w http.ResponseWriter, r *http.Request) {
	sess, _ := session.Store.Get(r, "auth-session")

	// Hapus sesi dengan mengatur MaxAge menjadi -1
	sess.Options.MaxAge = -1
	sess.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
