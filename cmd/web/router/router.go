package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kodekilat/go-ecommerce/cmd/web/handler" // Ganti dengan path modul Anda
)

func New() http.Handler {
	r := chi.NewRouter()

	// Middleware dasar
	r.Use(middleware.Logger)    // Mencatat log setiap request
	r.Use(middleware.Recoverer) // Memulihkan dari panic

	// Definisikan rute
	r.Get("/", handler.ShowHomePage)

	return r
}
