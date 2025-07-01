Source Code [Membangun Autentikasi Modern dengan Supabase di Next.js](https://liataja.id//tutorial/membuat-ecommerce-dengan-go-part1)

## Menjalankan aplikasi

Unduh atau clone repository ini

Buat file .env di root folder aplikasi
```
DATABASE_URL=""
SESSION_SECRET="kunci-rahasia-yang-sangat-panjang-dan-sulit-ditebak"
```

jangan lupa ditambahkan di file .gitignore
```
# .gitignore
.env
```

lalu jalankan server developement:

```bash
go run ./cmd/web/main.go
```

Halaman yang di buat
home page http://localhost:8080
account page http://localhost:8080/account
produk crud http://localhost:8080/admin/products

login  http://localhost:8080/login
register http://localhost:8080/register 

