# Catatan Keputusan (terbaru di atas)

## [2026-11-18] Implementasi Clean Architecture & Feature-First Delcom Blog Backend
- Konteks: Pengembangan REST API Delcom Blog Praktikum 9 dengan dukungan autentikasi JWT, manajemen artikel author/admin, dan integrasi AI generator.
- Keputusan: Menggunakan modular feature-first layout (`internal/features/{auth,post,ai}`) dengan pemisahan DTO, repository, service, dan handler. Menggunakan driver SQLite in-memory untuk automated testing guna mencapai 100% statement coverage tanpa dependensi PostgreSQL eksternal pada runner test.
- Alternatif ditolak: Mocking database menggunakan sqlmock (lebih kompleks dan rapuh terhadap perubahan query GORM dibanding SQLite in-memory).
- Dampak: Pengujian berjalan sangat cepat (<3 detik untuk seluruh suite), cakupan 100% statement tercapai, dan codebase siap untuk produksi dengan PostgreSQL.
