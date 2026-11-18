# Design Brief: Author Workspace & AI Generator

## Tujuan Layar
Menyediakan antarmuka bagi penulis untuk mengelola artikel (CRUD), memantau statistik publikasi (total artikel, total terbit, total draf), beralih status draf/terbit secara instan, serta membuat draf artikel menggunakan bantuan AI Generator.

## Alur Navigasi
1. Penulis login melalui `/login` dan diarahkan ke `/author`.
2. Halaman `/author` menampilkan ringkasan kartu statistik dan tabel artikel milik penulis.
3. Penulis dapat menekan tombol "Tulis Artikel Baru" untuk membuka form editor Markdown.
4. Penulis dapat membuka modal "AI Blog Generator", memasukkan topik dan kategori, serta mencentang opsi penelusuran web.
5. Hasil AI (draft artikel dan sitasi sumber rujukan) dapat ditinjau dan langsung diaplikasikan ke editor.
6. Penulis dapat menekan toggle untuk menerbitkan atau mengembalikan ke status draf.

## Komponen
| Komponen | Jenis | Teks Label |
|---|---|---|
| StatsCards | Stat Widget | "Total Artikel", "Artikel Terbit", "Draf Tersimpan" |
| ArticleTable | Data Table | Judul, Kategori, Status (Badge), Tanggal, Aksi (Edit, Toggle, Hapus) |
| MarkdownEditor | Toast UI Editor | Area penulisan Markdown dengan preview simultan |
| AIGeneratorModal | Modal Dialog | Input Topik, Pilihan Kategori, Checkbox Web Search, Tombol "Generate Draf" |
| CitationList | List | Sumber Rujukan Penelusuran Mesin Pencari |

## State UI
- **Default**: Menampilkan statistik dan daftar artikel penulis.
- **Loading**: Animasi loading saat generate AI atau menyimpan artikel.
- **Kosong**: Pesan "Anda belum menulis artikel apa pun. Mulai tulis artikel pertama Anda!".
- **Error**: SweetAlert notifikasi error validasi input atau otorisasi gagal.
- **Sukses**: Toast notifikasi berhasil disimpan/diterbitkan/dihapus.

## Teks UI Final
- "Kelola artikel dan pantau performa publikasi Anda."
- "Draf berhasil dibuat dengan bantuan AI. Tinjau sebelum diterbitkan."
- "Apakah Anda yakin ingin menghapus artikel ini?"

## Mobile dan Desktop
- Responsive card view pada mobile menggantikan tabel desktop.

## Pertanyaan Terbuka
- Tidak ada.
