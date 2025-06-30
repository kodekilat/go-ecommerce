package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func NewConnection() (*pgxpool.Pool, error) {
	// Muat variabel dari file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: Tidak dapat memuat file .env. Menggunakan variabel lingkungan sistem.")
	}

	// Ambil URL database dari environment variable
	databaseURL := os.Getenv("DATABASE_URL")
	log.Printf("Mencoba terhubung ke: %s", databaseURL)
	if databaseURL == "" {
		log.Fatal("DATABASE_URL tidak diset di environment variable")
	}

	// Buat koneksi pool ke database
	conn, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		return nil, err
	}

	// Uji koneksi
	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	log.Println("Berhasil terhubung ke database!")
	return conn, nil
}
