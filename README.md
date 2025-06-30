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

