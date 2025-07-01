package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kodekilat/go-ecommerce/internal/models" // Ganti dengan path modul Anda
)

type ProductRepository struct {
	DB *pgxpool.Pool
}

// CreateProduct menyimpan produk baru ke database
func (r *ProductRepository) CreateProduct(product *models.Product) error {
	query := `
		INSERT INTO products (name, description, price, stock, image_url)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	err := r.DB.QueryRow(context.Background(), query,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.ImageURL,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	return err
}
