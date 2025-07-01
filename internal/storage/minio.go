package storage

import (
	"context"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var MinioClient *minio.Client

func InitMinio() {
	endpoint := "127.0.0.1:9000"
	accessKeyID := "minioadmin"
	secretAccessKey := "miniosecret"
	useSSL := false

	// Inisialisasi Minio client
	var err error
	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalf("Gagal terhubung ke Minio: %v", err)
	}

	log.Println("Berhasil terhubung ke Minio!")

	// Buat bucket jika belum ada
	bucketName := "products"
	err = MinioClient.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{})
	if err != nil {
		exists, errBucketExists := MinioClient.BucketExists(context.Background(), bucketName)
		if errBucketExists == nil && exists {
			log.Printf("Bucket '%s' sudah ada.\n", bucketName)
		} else {
			log.Fatalf("Gagal membuat bucket: %v", err)
		}
	} else {
		log.Printf("Berhasil membuat bucket '%s'.\n", bucketName)
		// Set kebijakan bucket agar bisa diakses publik
		policy := `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject"],"Resource":["arn:aws:s3:::` + bucketName + `/*"]}]}`
		err = MinioClient.SetBucketPolicy(context.Background(), bucketName, policy)
		if err != nil {
			log.Fatalf("Gagal mengatur kebijakan bucket: %v", err)
		}
	}
}
