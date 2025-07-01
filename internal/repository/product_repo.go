package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (r *ProductRepository) GetAllProducts() ([]models.Product, error) {
	query := `SELECT id, name, price, image_url, stock FROM products ORDER BY created_at DESC`

	rows, err := r.DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.ImageURL, &product.Stock)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *ProductRepository) GetProductByID(id uuid.UUID) (*models.Product, error) {
	query := `SELECT id, name, description, price, stock, image_url FROM products WHERE id = $1`
	var product models.Product
	err := r.DB.QueryRow(context.Background(), query, id).Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.ImageURL)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// UpdateProduct memperbarui data produk yang ada di database.
func (r *ProductRepository) UpdateProduct(product *models.Product) error {
	query := `
		UPDATE products
		SET name = $1, description = $2, price = $3, stock = $4, image_url = $5, updated_at = NOW()
		WHERE id = $6
		RETURNING updated_at
	`
	// Kita hanya perlu mengambil updated_at untuk memperbarui model kita, jika diperlukan.
	err := r.DB.QueryRow(context.Background(), query,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.ImageURL,
		product.ID,
	).Scan(&product.UpdatedAt)

	return err
}

// DeleteProduct menghapus produk dari database berdasarkan ID.
func (r *ProductRepository) DeleteProduct(id uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`

	commandTag, err := r.DB.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	// Opsional: Periksa apakah ada baris yang benar-benar dihapus.
	if commandTag.RowsAffected() == 0 {
		return pgx.ErrNoRows // Mengindikasikan produk dengan ID tersebut tidak ada.
	}

	return nil
}
