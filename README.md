# 🚀 Panduan Setup & Menjalankan Project

## 1. Prasyarat Sistem (*Prerequisites*)

Pastikan komputer/laptop Anda sudah terinstall:
- **Go (Golang)**: Versi 1.20 atau yang lebih baru ([Download Go](https://go.dev/dl/))
- **PostgreSQL Database**: Berjalan di lokal (Port `5432`) atau PostgreSQL Cloud (Railway / Supabase / Neon)
- **Git**: Untuk clone repository ([Download Git](https://git-scm.com/))

---

## 2. Langkah 1: Clone Repository

Buka terminal (CMD, PowerShell, Git Bash, atau Terminal macOS/Linux) dan jalankan perintah:

```bash
# 1. Clone project dari repository
git clone https://github.com/haritsnb/quiz-bootcamp-golang-sanbercode.git

# 2. Masuk ke dalam direktori project
cd quiz-bootcamp-golang-sanbercode
```

---

## 3. Langkah 2: Setup Database PostgreSQL

1. Buka database management tool Anda (pgAdmin, DBeaver, TablePlus, atau CLI `psql`).
2. Buat sebuah database baru bernama **`local_golangquiz`**:
   ```sql
   CREATE DATABASE local_golangquiz;
   ```
*(Tabel dan data seeder **tidak perlu dibuat manual**, aplikasi akan membuatnya secara otomatis saat pertama kali dijalankan).*

---

## 4. Langkah 3: Konfigurasi File `.env`

Cari file bernama **`.env.example`** yang berada tepat di root folder project, lalu ubah menjadi nama file menjadi **`.env`**, lalu salin konfigurasi berikut (sesuaikan username dan password PostgreSQL komputer Anda):

```env
# ==========================================
# Application Configuration
# ==========================================
APP_PORT=8081

# ==========================================
# Database Configuration (PostgreSQL)
# ==========================================
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=root
DB_NAME=local_golangquiz
DB_SSLMODE=disable

# ==========================================
# JWT Secret Key
# ==========================================
JWT_SECRET=disini_text_acak_atau_sesuaikan
```

---

## 5. Langkah 4: Install Dependensi

Jalankan perintah berikut di terminal untuk mengunduh semua package:

```bash
go mod tidy
```

---

## 6. Langkah 5: Jalankan Aplikasi

Jalankan perintah berikut untuk meng-generate dokumentasi Swagger dan menyalakan server backend:

```bash
go run main.go
```

### Output Sukses di Terminal:
```text
Berhasil terhubung ke database PostgreSQL!
Migrasi berhasil dijalankan: 3 migrasi diterapkan.
Seeding user berhasil: haritsnb / secretbanget
[GIN-debug] GET    /swagger/*any             --> ...
[GIN-debug] POST   /api/users/login          --> ...
[GIN-debug] Listening and serving HTTP on :8081
```

---

## 7. Langkah 6: Mengakses dan Menggunakan Swagger UI

### A. Buka di Browser
Buka browser (Google Chrome / Edge / Firefox) dan kunjungi URL:
```url
http://localhost:8081/swagger/index.html
```

---

### B. Cara Login & Otentikasi di Swagger UI:

1. **Login untuk Mendapatkan Token:**
   - Di Swagger UI, cari grup **`Auth`** --> klik **`POST /users/login`**.
   - Klik tombol **`Try it out`**.
   - Pada request body, pastikan isinya:
     ```json
     {
       "username": "haritsnb",
       "password": "secretbanget"
     }
     ```
   - Klik tombol biru **`Execute`**.
   - Pada bagian *Server response*, salin string nilai **`token`** yang muncul (tanpa tanda kutip `"`).

2. **Memasang Token (Authorize):**
   - Gulir ke bagian paling atas halaman Swagger UI.
   - Klik tombol hijau **`Authorize 🔓`** di sebelah kanan atas.
   - Pada kolom input **Value**, ketik `Bearer ` lalu paste token Anda:
     ```text
     Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
     ```
   - Klik tombol **Authorize** --> klik **Close**.
   - Icon gembok akan berubah menjadi terkunci **🔒**.

---

### C. Mencoba Endpoint CRUD & Upload Gambar di Swagger:

1. **Melihat Daftar Seluruh Kategori:**
   - Klik **`GET /categories`** --> klik **`Try it out`** --> **`Execute`**.

2. **Upload Buku Baru Beserta Gambar:**
   - Klik **`POST /books`** --> klik **`Try it out`**.
   - Isi field:
     - `title`: `Pemrograman Golang Expert`
     - `description`: `Buku panduan arsitektur modern`
     - `release_year`: `2023`
     - `price`: `150000`
     - `total_page`: `250` *(karena > 100, otomatis jadi thickness "tebal")*
     - `category_id`: `1`
     - `image`: Klik tombol **Choose File** / pilih gambar dari laptop.
   - Klik **`Execute`**.
   - Response akan mengembalikan data buku lengkap dengan **`image_url` berupa link gambar aktif** yang bisa langsung diklik dan dibuka di tab baru browser!