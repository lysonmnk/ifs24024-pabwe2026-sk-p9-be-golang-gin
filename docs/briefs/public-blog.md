# Design Brief: Public Blog

## Tujuan Layar
Menampilkan galeri artikel yang telah diterbitkan untuk pengunjung publik, memungkinkan pencarian kata kunci secara real-time, filter kategori, dan membaca detail artikel lengkap beserta estimasi waktu baca.

## Alur Navigasi
1. Pengunjung membuka `/` atau `/public`.
2. Pengunjung dapat mengetik di kolom pencarian atau memilih pill kategori.
3. Klik salah satu kartu artikel akan mengarahkan ke `/posts/:slug`.
4. Halaman detail menyajikan judul, kategori, nama penulis, tanggal terbit, estimasi reading time, dan konten Markdown yang dirender.

## Komponen
| Komponen | Jenis | Teks Label |
|---|---|---|
| Navbar | Header | Logo Delcom Blog, Search Input, Tombol Login/Register |
| CategoryFilter | Filter Chips | "Semua", "Teknologi", "Tutorial", "Gaya Hidup", dsb. |
| PostCard | Card | Judul, Excerpt, Reading Time (e.g. "3 menit baca"), Tanggal, Penulis |
| PostDetail | Article View | Judul Utama, Meta Author & Date, Content Markdown Viewer |

## State UI
- **Default**: Menampilkan grid kartu artikel terbaru.
- **Loading**: Skeleton placeholder untuk kartu artikel.
- **Kosong**: Pesan "Belum ada artikel yang diterbitkan" atau "Tidak ada artikel yang cocok dengan pencarian Anda".
- **Error**: Pesan notifikasi gagal memuat artikel dengan tombol "Coba Lagi".
- **Sukses**: Artikel berhasil dimuat.

## Teks UI Final
- "Temukan artikel menarik seputar teknologi dan inovasi."
- "Estimasi waktu baca: X menit"
- "Dipublikasikan oleh: [Nama Penulis]"

## Mobile dan Desktop
- Desktop: Grid 3 kolom untuk kartu artikel.
- Mobile: Grid 1 kolom dengan drawer filter kategori.

## Pertanyaan Terbuka
- Tidak ada.
