# Changelog

Semua perubahan penting pada proyek ini dicatat dalam berkas ini.

## [1.0.0] - 2026-11-18

### Ditambahkan
- Inisialisasi modul Go `app` dan instalasi seluruh dependensi Gin, GORM, JWT, Crypto, Swag, Testify.
- Konfigurasi environment di `internal/config` dengan dukungan fallback dan `.env.example`.
- Inisialisasi basis data PostgreSQL & SQLite di `internal/database` dengan auto-migration dan initial admin seeder.
- Standard format response JSON di `internal/response`.
- Middleware: JWT authentication, CORS header handling, role-based access control, dan in-memory rate limiter di `internal/middleware`.
- Fitur Autentikasi (`internal/features/auth`): registrasi author, login JWT, profil, dan logout.
- Fitur Artikel (`internal/features/post`): kategori, artikel publik dengan filter kata kunci & kategori, detail artikel slug, CRUD artikel author, generator slug otomatis, estimasi waktu baca (200 wpm), toggle publish instan, dan statistik artikel.
- Fitur Administrasi (`internal/features/post`): pemantauan seluruh pengguna, arsip artikel semua penulis, dan artikel spesifik per pengguna.
- Fitur AI Generator (`internal/features/ai`): integrasi Search Engine endpoint untuk referensi web dan OpenAI Chat Completions untuk draf artikel berformat Markdown beserta sitasi.
- Dokumentasi API interaktif Swagger UI pada `/swagger/*any` dan spesifikasi OpenAPI pada `/api.json`.
- Pengujian unit dan fitur otomatis di seluruh paket `internal/...` dengan pencapaian **100% statement coverage**.
- Dokumen proyek: `docs/PRD.md`, `docs/ARCHITECTURE.md`, `docs/BACKLOG.md`, `docs/briefs/`, `contracts/`, `docs/SECURITY_REPORT.md`, `docs/QA_REPORT.md`.
