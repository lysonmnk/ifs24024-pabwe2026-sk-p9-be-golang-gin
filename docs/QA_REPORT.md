# QA Report
Fitur: Delcom Blog Backend API Suite | Tanggal: 2026-11-18

| ID | Skenario | Langkah | Diharapkan | Aktual | Status | Severity |
|---|---|---|---|---|---|---|
| TC-01 | Registrasi Akun Penulis | POST /api/auth/register dengan kredensial valid | Status 201 Created dan user data | Status 201 Created dan user data | Lulus | None |
| TC-02 | Login Akun Penulis | POST /api/auth/login dengan email dan password benar | Status 200 OK dan JWT Bearer token | Status 200 OK dan JWT Bearer token | Lulus | None |
| TC-03 | Login Password Salah | POST /api/auth/login dengan password keliru | Status 401 Unauthorized | Status 401 Unauthorized | Lulus | None |
| TC-04 | Akses Profil Terautentikasi | GET /api/auth/profile dengan Bearer token | Status 200 OK dengan profil author | Status 200 OK dengan profil author | Lulus | None |
| TC-05 | Pembuatan Artikel dengan Slug Unik | POST /api/posts dengan judul duplikat | Slug otomatis diberi suffix numerik unik (-1, dsb.) | Slug unik otomatis terbentuk | Lulus | None |
| TC-06 | Estimasi Waktu Baca Artikel | POST /api/posts dengan konten teks | Kalkulasi reading time berbasis 200 wpm | Waktu baca terhitung akurat | Lulus | None |
| TC-07 | Toggle Status Publikasi | PATCH /api/posts/:id/toggle-publish | Status beralih antara draf dan terbit beserta timestamp | Status beralih dan published_at diperbarui | Lulus | None |
| TC-08 | Proteksi IDOR Antar Penulis | PUT /api/posts/:id milik penulis lain | Status 400/403 dengan pesan akses ditolak | Ditolak dengan pesan otorisasi | Lulus | None |
| TC-09 | Akses Panel Admin oleh Non-Admin | GET /api/admin/users dengan role author | Status 403 Forbidden | Status 403 Forbidden | Lulus | None |
| TC-10 | Akses Panel Admin oleh Admin | GET /api/admin/users dengan role admin | Status 200 OK dan daftar pengguna lengkap | Status 200 OK dan daftar pengguna | Lulus | None |
| TC-11 | AI Blog Generator dengan Search Engine | POST /api/posts/ai-generate dengan IncludeWebSearch=true | Draf artikel terstruktur dengan sitasi sumber | Draf artikel dan daftar rujukan dikembalikan | Lulus | None |
| TC-12 | Penanganan Route Tidak Dikenal | GET /api/unknown-endpoint | Status 404 Not Found dengan pesan "Endpoint not found" | Status 404 Not Found | Lulus | None |

Verdict: SIAP RILIS
