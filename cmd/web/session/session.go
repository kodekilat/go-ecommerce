package session

import (
	"encoding/gob"
	"os"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore

func init() {
	gob.Register(uuid.UUID{})
	// Ambil kunci dari environment variable
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		// Di produksi, ini seharusnya menyebabkan panic
		secret = "kunci-default-hanya-untuk-pengembangan"
	}
	Store = sessions.NewCookieStore([]byte(secret))

	// Konfigurasi opsi cookie
	Store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 hari
		HttpOnly: true,      // Mencegah akses via JavaScript
		// Secure: true,     // Aktifkan di produksi (membutuhkan HTTPS)
	}
}
