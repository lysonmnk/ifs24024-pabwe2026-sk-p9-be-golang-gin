# Design Brief: Admin Dashboard

## Tujuan Layar
Menyediakan panel terpadu bagi administrator untuk memantau seluruh pengguna terdaftar, jumlah artikel per pengguna, serta melihat arsip lengkap seluruh artikel dalam sistem.

## Alur Navigasi
1. Admin login dengan akun yang memiliki role `admin`, lalu diarahkan atau membuka `/admin`.
2. Halaman memuat daftar pengguna beserta ringkasan artikel.
3. Admin dapat mengklik salah satu pengguna untuk melihat daftar artikel spesifik milik pengguna tersebut.
4. Tab "Arsip Artikel" menyajikan seluruh artikel dari semua pengguna tanpa batasan status.

## Komponen
| Komponen | Jenis | Teks Label |
|---|---|---|
| UserListTable | Data Table | ID, Nama, Email, Role, Total Artikel, Aksi ("Lihat Artikel") |
| PostArchiveTable | Data Table | Judul, Penulis, Kategori, Status Publikasi, Tanggal Dibuat |
| UserPostsModal | Modal / Drawer | Daftar artikel dari pengguna terpilih |

## State UI
- **Default**: Menampilkan tabel pengguna dan tab navigasi.
- **Loading**: Indikator loading pada tabel.
- **Kosong**: "Tidak ada pengguna / artikel ditemukan".
- **Error**: Pesan 403 Forbidden bila diakses non-admin.
- **Sukses**: Data termuat dengan rapi.

## Teks UI Final
- "Panel Administrator: Pengawasan Pengguna & Arsip Konten."
- "Hak akses terbatas untuk Administrator."

## Mobile dan Desktop
- Horizontal scroll pada tabel untuk perangkat mobile, tampilan grid penuh pada desktop.

## Pertanyaan Terbuka
- Tidak ada.
