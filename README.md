# Delcom Blog

Delcom Blog adalah platform publikasi artikel multi-tier yang dibangun dengan arsitektur Clean Architecture dan Feature-First. Aplikasi ini dilengkapi editor artikel berbasis Markdown serta AI Blog Generator yang terintegrasi dengan Search Engine untuk menghasilkan draf artikel yang relevan dan disertai sumber rujukan.

---

## Daftar Isi

- [Fitur Utama](#fitur-utama)
- [Teknologi](#teknologi)
- [Arsitektur](#arsitektur)
- [Peran Pengguna](#peran-pengguna)
- [Prasyarat](#prasyarat)
- [Instalasi dan Konfigurasi](#instalasi-dan-konfigurasi)
- [Menjalankan Aplikasi](#menjalankan-aplikasi)
- [Pengujian](#pengujian)
- [Dokumentasi API](#dokumentasi-api)
- [Variabel Lingkungan](#variabel-lingkungan)
- [Struktur Proyek](#struktur-proyek)
- [Deployment](#deployment)
- [Catatan Produksi](#catatan-produksi)
- [Lisensi](#lisensi)

---

## Fitur Utama

- Publikasi artikel dengan tiga tingkat akses: Publik, Penulis, dan Admin.
- Editor Markdown dengan pratinjau simultan berbasis Toast UI Editor.
- Pembuatan slug unik secara otomatis untuk setiap artikel.
- Perhitungan estimasi waktu baca dengan asumsi 200 kata per menit.
- Pencarian kata kunci secara real-time dan filter berdasarkan kategori.
- Pengaturan status publikasi (draf dan terbit) secara instan.
- AI Blog Generator yang menggabungkan hasil penelusuran Search Engine dengan OpenAI Chat Completions, lengkap dengan daftar sumber rujukan.
- Autentikasi berbasis JWT dan kontrol akses berdasarkan peran (role).
- Dokumentasi API interaktif melalui Swagger UI.
- Cakupan pengujian 100% pada backend dan frontend.

---

## Teknologi

| Lapisan | Teknologi |
|---|---|
| Bahasa backend | Go 1.27.4 |
| Framework HTTP | Gin v1.12.0 |
| ORM | GORM v1.31.2 (driver PostgreSQL v1.6.3) |
| Basis data | PostgreSQL 17 (SQLite in-memory untuk pengujian) |
| Autentikasi | JWT (golang-jwt/jwt v5.3.1) |
| Dokumentasi API | Swag v1.16.6 dan Gin Swagger v1.6.1 |
| Frontend | React, Vite, Tailwind CSS |
| Editor Markdown | @toast-ui/editor |
| Notifikasi | SweetAlert2 |
| Ikon | @tabler/icons-react |
| Pengujian backend | Go testing dan Testify v1.12.1 |
| Pengujian frontend | Vitest, Testing Library, jsdom |

---

## Arsitektur

Backend mengikuti prinsip Clean Architecture dengan pengelompokan berbasis fitur (Feature-First). Setiap fitur memiliki lapisan yang terpisah untuk DTO, repository, service, dan handler.

Frontend adalah Single Page Application (SPA) yang mengonsumsi REST API backend. Navigasi dikelola melalui sinkronisasi URL peramban (`window.history.pushState` dan event `popstate`), sehingga halaman tetap berada pada rute yang sama saat di-refresh.

Alur AI Blog Generator:

1. Penulis mengirim topik dan kategori ke `POST /api/posts/ai-generate` dengan `IncludeWebSearch=true`.
2. Backend melakukan kueri topik ke Search Engine Endpoint dan memperoleh snippet serta referensi web.
3. Backend mengirim prompt beserta konteks hasil penelusuran ke OpenAI Chat Completions.
4. Model mengembalikan artikel dalam format JSON (judul, kategori, excerpt, konten Markdown, dan sitasi).
5. Backend mengembalikan artikel lengkap beserta daftar sumber penelusuran kepada frontend.
6. Penulis meninjau pratinjau dan menerapkannya langsung ke editor.

---

## Peran Pengguna

| Peran | Rute | Hak Akses |
|---|---|---|
| Publik (Guest) | `/`, `/public`, `/posts/:slug` | Membaca artikel terbit, mencari, memfilter kategori, dan melihat estimasi waktu baca. |
| Penulis (Author) | `/author` | Mengelola artikel sendiri (CRUD), melihat statistik, mengubah status publikasi, dan menggunakan AI Generator. Memerlukan login. |
| Admin | `/admin` | Memantau seluruh pengguna dan seluruh artikel dalam sistem. Memerlukan login dengan peran admin. |

Rute tambahan: `/login` dan `/register`.

---

## Prasyarat

Pastikan perangkat sudah terpasang:

- Go 1.27.4 atau yang lebih baru
- PostgreSQL 17
- Node.js 20 atau yang lebih baru beserta npm
- Git

---

## Instalasi dan Konfigurasi

### 1. Backend (Golang Gin)

Nama repositori mengikuti pola `{username}-pabwe2026-sk-p9-be-golang-gin`.

```bash
git clone https://github.com/{username}/{username}-pabwe2026-sk-p9-be-golang-gin.git
cd {username}-pabwe2026-sk-p9-be-golang-gin
go mod init app
```

Pasang dependensi utama:

```bash
go get github.com/gin-gonic/gin@v1.12.0
go get gorm.io/gorm@v1.31.2
go get gorm.io/driver/postgres@v1.6.3
go get github.com/glebarez/sqlite@v1.11.0
go get github.com/golang-jwt/jwt/v5@v5.3.1
go get golang.org/x/crypto@v0.59.0
go get github.com/joho/godotenv@v1.5.1
go get github.com/swaggo/gin-swagger@v1.6.1
go get github.com/swaggo/files@v1.0.1
go get github.com/swaggo/swag@v1.16.6
go get github.com/stretchr/testify@v1.12.1
```

Buat basis data PostgreSQL:

```bash
psql -U postgres -c "CREATE DATABASE pabwe_p9_blog;"
```

Salin berkas contoh konfigurasi lalu sesuaikan nilainya:

```bash
cp .env.example .env
```

### 2. Frontend (ReactJS Vite)

Nama repositori mengikuti pola `{username}-pabwe2026-sk-p9-fe-reactjs`.

```bash
git clone https://github.com/{username}/{username}-pabwe2026-sk-p9-fe-reactjs.git
cd {username}-pabwe2026-sk-p9-fe-reactjs
npm install
```

---

## Menjalankan Aplikasi

Jalankan backend pada port 8081:

```bash
cd {username}-pabwe2026-sk-p9-be-golang-gin
go run cmd/server/main.go
```

Output yang diharapkan:

```
[GIN-debug] Listening and serving HTTP on :8081
```

Jalankan frontend pada port 5174 di terminal terpisah:

```bash
cd {username}-pabwe2026-sk-p9-fe-reactjs
npm run dev
```

Buka aplikasi di `http://localhost:5174`. Vite meneruskan permintaan `/api` ke backend di `http://localhost:8081`.

| Layanan | Alamat |
|---|---|
| Frontend | http://localhost:5174 |
| Backend API | http://localhost:8081/api |
| Swagger UI | http://localhost:8081/swagger/index.html |

---

## Pengujian

### Backend

```bash
go test -v -cover ./internal/...
```

Seluruh paket harus lulus dengan cakupan 100% statement:

```
ok  app/internal/config          coverage: 100.0% of statements
ok  app/internal/database        coverage: 100.0% of statements
ok  app/internal/middleware      coverage: 100.0% of statements
ok  app/internal/response        coverage: 100.0% of statements
ok  app/internal/features/auth   coverage: 100.0% of statements
ok  app/internal/features/post   coverage: 100.0% of statements
ok  app/internal/features/ai     coverage: 100.0% of statements
ok  app/internal/server          coverage: 100.0% of statements
```

### Frontend

```bash
npm test
```

Seluruh berkas harus lulus dengan 100% Statements, Branches, Functions, dan Lines. Ambang batas tersebut telah dikonfigurasi di `vite.config.js`.

---

## Dokumentasi API

Seluruh endpoint menggunakan prefiks `/api` dan mengembalikan respons JSON dengan format standar:

```json
{
  "status": "success",
  "message": "...",
  "data": {}
}
```

Rute yang tidak dikenal dikembalikan dengan status `404 Not Found` dan pesan `Endpoint not found`.

### Endpoint Publik

| Method | Endpoint | Keterangan |
|---|---|---|
| GET | `/swagger/*any` | Dokumentasi interaktif Swagger UI |
| GET | `/api.json` | Spesifikasi OpenAPI dalam format JSON |
| GET | `/api/public/categories` | Daftar kategori artikel |
| GET | `/api/public/posts` | Daftar artikel terbit, mendukung `?q=` dan `?category=` |
| GET | `/api/public/posts/:slug` | Detail artikel terbit beserta data penulis dan waktu baca |
| POST | `/api/auth/register` | Pendaftaran akun penulis baru |
| POST | `/api/auth/login` | Autentikasi dan penerbitan token JWT |

### Endpoint Penulis (memerlukan autentikasi)

| Method | Endpoint | Keterangan |
|---|---|---|
| POST | `/api/auth/logout` | Mengakhiri sesi pengguna |
| GET | `/api/auth/profile` | Mengambil profil penulis yang sedang login |
| PUT | `/api/auth/profile` | Memperbarui profil penulis |
| GET | `/api/posts` | Daftar seluruh artikel milik penulis (terbit dan draf) |
| POST | `/api/posts` | Membuat artikel baru beserta slug unik dan waktu baca |
| GET | `/api/posts/stats` | Statistik: total artikel, total terbit, dan total draf |
| GET | `/api/posts/:id` | Detail artikel milik penulis untuk keperluan penyuntingan |
| PUT | `/api/posts/:id` | Memperbarui artikel milik penulis |
| PATCH | `/api/posts/:id/toggle-publish` | Mengubah status draf dan terbit beserta `published_at` |
| DELETE | `/api/posts/:id` | Menghapus artikel milik penulis |
| POST | `/api/posts/ai-generate` | Menghasilkan draf artikel dengan AI dan penelusuran Search Engine |

### Endpoint Admin (hanya untuk peran admin)

| Method | Endpoint | Keterangan |
|---|---|---|
| GET | `/api/admin/users` | Daftar seluruh pengguna beserta jumlah artikel |
| GET | `/api/admin/posts` | Arsip seluruh artikel dari semua penulis |
| GET | `/api/admin/users/:id/posts` | Daftar artikel milik pengguna tertentu |

---

## Variabel Lingkungan

Berkas `.env` berada di root direktori backend. Ganti seluruh nilai contoh dengan konfigurasi Anda sendiri.

| Variabel | Contoh Nilai | Keterangan |
|---|---|---|
| `APP_ENV` | `development` | Lingkungan aplikasi |
| `APP_PORT` | `8081` | Port backend |
| `APP_KEY` | `<kunci-rahasia-minimal-32-karakter>` | Kunci penandatanganan JWT |
| `ADMIN_USERNAME` | `admin` | Username akun admin awal |
| `ADMIN_PASSWORD` | `<password-admin>` | Password akun admin awal |
| `DB_CONNECTION` | `pgsql` | Jenis basis data |
| `DB_HOST` | `127.0.0.1` | Host PostgreSQL |
| `DB_PORT` | `5432` | Port PostgreSQL |
| `DB_DATABASE` | `pabwe_p9_blog` | Nama basis data |
| `DB_USERNAME` | `postgres` | Pengguna basis data |
| `DB_PASSWORD` | `<password-database>` | Password basis data |
| `DB_SSLMODE` | `disable` | Mode SSL PostgreSQL |
| `SEARCH_ENGINE_ENDPOINT` | `https://<host>/api/v1/search` | Endpoint Search Engine |
| `OPENAI_ENDPOINT` | `https://<host>/api/v1` | Endpoint layanan kompatibel OpenAI |
| `OPENAI_KEY` | `<api-key>` | Kunci API layanan AI |
| `OPENAI_MODEL` | `<nama-model>` | Nama model yang digunakan |

---

## Struktur Proyek

```
{username}-pabwe2026-sk-p9-be-golang-gin/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point server Gin
├── docs/                           # Spesifikasi OpenAPI dan Swagger
├── internal/
│   ├── config/                     # Pemuatan konfigurasi environment
│   ├── database/                   # Koneksi, migrasi, dan seeding basis data
│   ├── middleware/                 # JWT, CORS, RequireRole, RateLimiter
│   ├── response/                   # Format respons JSON standar
│   ├── features/
│   │   ├── auth/                   # Registrasi, login, dan profil
│   │   ├── post/                   # CRUD artikel, slug, waktu baca, statistik
│   │   └── ai/                     # Klien Search Engine dan OpenAI
│   └── server/                     # Router Gin dan penyedia Swagger UI
├── .env.example
├── .gitignore
├── go.mod
└── README.md
```

```
{username}-pabwe2026-sk-p9-fe-reactjs/
├── src/
│   ├── main.jsx                    # Root React
│   ├── App.jsx                     # Komponen utama dan routing SPA
│   ├── components/                 # MarkdownEditor dan MarkdownViewer
│   ├── helpers/                    # apiHelper dan toolsHelper
│   ├── hooks/                      # useInput dan useTheme
│   └── features/
│       ├── admin/                  # Panel admin
│       ├── auth/                   # Login, registrasi, dan modal autentikasi
│       └── blog/                   # Galeri publik, detail artikel, workspace penulis, editor, dan AI Generator
├── package.json
├── vite.config.js
└── index.html
```

Struktur tersebut menggambarkan pemisahan lapisan pada backend (handler, service, repository) dan pengelompokan komponen berdasarkan fitur pada frontend.

---

## Deployment

Backend dan frontend dapat di-deploy secara terpisah (decoupled deployment).

- Backend: dapat dijalankan sebagai container Docker atau pada VPS maupun layanan cloud. Pastikan variabel lingkungan telah dikonfigurasi dan CORS mengizinkan domain frontend.
- Frontend: jalankan `npm run build` untuk menghasilkan direktori `dist/`, lalu unggah ke layanan hosting statis seperti Vercel, Netlify, Cloudflare Pages, atau AWS S3.

---

## Catatan Produksi

- Ganti `APP_KEY`, `ADMIN_PASSWORD`, dan `DB_PASSWORD` dengan nilai yang kuat dan unik.
- Jangan pernah menyimpan berkas `.env` atau `OPENAI_KEY` di dalam repositori.
- Aktifkan `DB_SSLMODE=require` atau sesuai kebijakan keamanan basis data Anda.
- Batasi origin CORS hanya pada domain frontend resmi.

---

## Integrasi CI/CD

Aturan `.gitignore` untuk biner backend wajib diawali garis miring (`/server`, `/main`, `/app`). Dengan aturan tersebut, runner Jenkins tidak secara keliru mengabaikan direktori source code seperti `internal/server/` saat proses checkout.

---

## Lisensi

Proyek ini dibuat untuk keperluan pembelajaran pada mata kuliah Pengembangan Aplikasi Berbasis Web (PABWE) 2026. Silakan tambahkan lisensi sesuai kebutuhan.

---

## Kontributor
Lyson 
