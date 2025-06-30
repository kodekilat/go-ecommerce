package middleware

import (
	"context"
	"net/http"

	"github.com/kodekilat/go-ecommerce/cmd/web/session" // Ganti path modul Anda
)

// Definisikan kunci unik untuk konteks
type contextKey string

const userContextKey = contextKey("user")

// Authenticate adalah middleware yang memeriksa sesi pengguna
func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Dapatkan sesi dari request
		sess, _ := session.Store.Get(r, "auth-session")

		// 2. Periksa apakah user_id ada di sesi
		userID, ok := sess.Values["user_id"]
		if !ok || userID == nil {
			// Jika tidak ada, panggil handler berikutnya tanpa melakukan apa-apa
			// (artinya pengguna adalah tamu/guest)
			next.ServeHTTP(w, r)
			return
		}

		// 3. Jika user_id ada, simpan di konteks request
		ctx := context.WithValue(r.Context(), userContextKey, userID)

		// 4. Panggil handler berikutnya dengan request yang sudah dimodifikasi konteksnya
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireAuthentication adalah middleware yang menolak akses jika pengguna belum login
func RequireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ambil user_id dari konteks
		userID := r.Context().Value(userContextKey)

		// Jika tidak ada user_id (pengguna adalah tamu), redirect ke halaman login
		if userID == nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Jika ada, lanjutkan ke handler tujuan
		next.ServeHTTP(w, r)
	})
}
