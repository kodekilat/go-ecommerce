package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kodekilat/go-ecommerce/cmd/web/handler" // Ganti dengan path modul Anda
	authMiddleware "github.com/kodekilat/go-ecommerce/cmd/web/middleware"
	"github.com/kodekilat/go-ecommerce/internal/repository"
)

func New(userRepo *repository.UserRepository, productRepo *repository.ProductRepository) http.Handler {
	r := chi.NewRouter()

	// Middleware dasar
	r.Use(middleware.Logger)    // Mencatat log setiap request
	r.Use(middleware.Recoverer) // Memulihkan dari panic
	r.Use(authMiddleware.Authenticate)

	authHandler := &handler.AuthHandler{UserRepo: userRepo}
	adminHandler := &handler.AdminHandler{ProductRepo: productRepo}
	homeHandler := &handler.HomeHandler{ProductRepo: productRepo}
	productHandler := &handler.ProductHandler{ProductRepo: productRepo}

	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))))

	// Definisikan rute
	r.Get("/", homeHandler.ShowHomePage)
	r.Get("/products/{productID}", productHandler.ShowProductDetail)

	r.Get("/register", authHandler.ShowRegistrationForm)
	r.Post("/register", authHandler.HandleRegistration)

	r.Get("/login", authHandler.ShowLoginForm)
	r.Post("/login", authHandler.HandleLogin)

	r.Post("/logout", authHandler.HandleLogout)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.RequireAuthentication)

		// Daftarkan rute yang dilindungi di sini
		r.Get("/account", handler.ShowAccountPage)

		r.Get("/admin/products", adminHandler.ShowAdminProducts)
		r.Post("/admin/products", adminHandler.HandleAddProduct)

		r.Get("/admin/products/edit/{productID}", adminHandler.ShowEditProductForm)
		r.Post("/admin/products/edit/{productID}", adminHandler.HandleUpdateProduct)

		r.Post("/admin/products/delete/{productID}", adminHandler.HandleDeleteProduct)
	})

	return r
}
