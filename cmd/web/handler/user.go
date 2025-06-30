package handler

import (
	"net/http"

	"github.com/kodekilat/go-ecommerce/cmd/web/session"
	"github.com/kodekilat/go-ecommerce/cmd/web/view"
)

// Definisikan kunci konteks yang sama
type contextKey string

const userContextKey = contextKey("user")

func ShowAccountPage(w http.ResponseWriter, r *http.Request) {
	// Ambil data dari sesi (cara alternatif selain dari konteks)
	sess, _ := session.Store.Get(r, "auth-session")
	userID := sess.Values["user_id"]
	email := sess.Values["user_email"]

	pageData := struct {
		UserID interface{}
		Email  interface{}
	}{
		UserID: userID,
		Email:  email,
	}

	view.Render(w, "account.page.html", pageData)
}
