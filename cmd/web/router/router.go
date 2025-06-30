package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kodekilat/go-ecommerce/cmd/web/handler" // Ganti dengan path modul Anda
	"github.com/kodekilat/go-ecommerce/internal/repository"
)

func New(userRepo *repository.UserRepository) http.Handler {
	r := chi.NewRouter()

	// Middleware dasar
	r.Use(middleware.Logger)    // Mencatat log setiap request
	r.Use(middleware.Recoverer) // Memulihkan dari panic

	authHandler := &handler.AuthHandler{UserRepo: userRepo}

	// Definisikan rute
	r.Get("/", handler.ShowHomePage)
	r.Get("/register", authHandler.ShowRegistrationForm)
	r.Post("/register", authHandler.HandleRegistration)

	return r
}
