# Backlog

| ID | Fitur | Prioritas | Acceptance Criteria | Agen | Status |
|---|---|---|---|---|---|
| T-01 | Config & Database Init | P0 | Given file `.env` ada, When aplikasi di-boot, Then konfigurasi termuat, koneksi DB terbentuk, dan admin awal di-seed jika belum ada. | Backend | Selesai |
| T-02 | Response & Error Middleware | P0 | Given request masuk, When endpoint tidak ditemukan atau handler mengembalikan response, Then format response JSON konsisten `{status, message, data}` dan 404 menghasilkan `{status: "error", message: "Endpoint not found"}`. | Backend | Selesai |
| T-03 | Auth Feature (Register & Login) | P0 | Given kredensial user, When POST `/api/auth/register` atau `/api/auth/login`, Then validasi data berjalan, password di-hash bcrypt, token JWT diterbitkan, dan error 401/422 ditangani. | Backend | Selesai |
| T-04 | Auth Feature (Profile & Logout) | P0 | Given user terautentikasi dengan Bearer token, When GET atau PUT `/api/auth/profile` atau POST `/api/auth/logout`, Then data profil dikembalikan/diperbarui atau sesi diakhiri. | Backend | Selesai |
| T-05 | Public Posts & Categories | P0 | Given artikel terbit, When GET `/api/public/posts` (dengan filter kata kunci `q` atau `category`), GET `/api/public/categories`, atau GET `/api/public/posts/:slug`, Then data dikembalikan dengan informasi penulis dan reading time. | Backend | Selesai |
| T-06 | Author Posts Management | P0 | Given author terautentikasi, When CRUD `/api/posts`, toggle publish, atau get stats, Then artikel tersimpan dengan auto-slug unik dan kalkulasi reading time (200 wpm), dan author hanya dapat mengubah artikel miliknya. | Backend | Selesai |
| T-07 | Admin Panel Endpoints | P1 | Given akun dengan role admin, When akses `/api/admin/users`, `/api/admin/posts`, atau `/api/admin/users/:id/posts`, Then seluruh data sistem dapat diakses; jika non-admin akses maka return 403 Forbidden. | Backend | Selesai |
| T-08 | AI Blog Generator | P1 | Given topik dan opsi web search, When POST `/api/posts/ai-generate`, Then sistem mencari konteks via Search Engine (jika enabled), memanggil OpenAI Chat Completions, dan mengembalikan draf artikel terstruktur dengan sitasi sumber rujukan. | Backend | Selesai |
| T-09 | Swagger Documentation | P1 | Given server berjalan, When GET `/swagger/*any` atau `/api.json`, Then dokumentasi interaktif Swagger UI dan OpenAPI spec JSON dapat diakses dengan benar. | Backend | Selesai |
| T-10 | SQA Automated Tests | P0 | Given seluruh paket di `internal/...`, When `go test -v -cover ./internal/...` dijalankan, Then seluruh tes lulus dengan 100% statement coverage. | SQA | Selesai |
