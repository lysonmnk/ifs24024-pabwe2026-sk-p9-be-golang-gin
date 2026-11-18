# Contract: Public Posts & Categories
Versi: 1.0 | Pemilik: Data Master | Status: Disepakati

## 1. Get Categories
### Endpoint
`GET /api/public/categories` | Auth: tidak

### Response Sukses (200 OK)
```json
{
  "status": "success",
  "message": "Daftar kategori berhasil dimuat.",
  "data": [
    "Teknologi",
    "Tutorial",
    "Pengembangan Web",
    "Kecerdasan Buatan"
  ]
}
```

---

## 2. Get Published Posts
### Endpoint
`GET /api/public/posts?q={keyword}&category={category}` | Auth: tidak

### Query Parameters
- `q` (opsional): kata kunci pencarian pada judul atau konten.
- `category` (opsional): filter nama kategori.

### Response Sukses (200 OK)
```json
{
  "status": "success",
  "message": "Daftar artikel publik berhasil dimuat.",
  "data": [
    {
      "id": 1,
      "user_id": 2,
      "title": "Memulai Belajar Golang Gin",
      "slug": "memulai-belajar-golang-gin",
      "category": "Teknologi",
      "excerpt": "Panduan ringkas arsitektur REST API dengan Go Gin.",
      "reading_time": 3,
      "is_published": true,
      "published_at": "2026-11-18T08:00:00Z",
      "created_at": "2026-11-18T07:30:00Z",
      "author": {
        "id": 2,
        "name": "Budi",
        "username": "budi"
      }
    }
  ]
}
```

---

## 3. Get Post Detail by Slug
### Endpoint
`GET /api/public/posts/:slug` | Auth: tidak

### Response Sukses (200 OK)
```json
{
  "status": "success",
  "message": "Detail artikel berhasil dimuat.",
  "data": {
    "id": 1,
    "user_id": 2,
    "title": "Memulai Belajar Golang Gin",
    "slug": "memulai-belajar-golang-gin",
    "category": "Teknologi",
    "excerpt": "Panduan ringkas arsitektur REST API dengan Go Gin.",
    "content": "# Konten Markdown Lengkap...",
    "reading_time": 3,
    "is_published": true,
    "published_at": "2026-11-18T08:00:00Z",
    "created_at": "2026-11-18T07:30:00Z",
    "author": {
      "id": 2,
      "name": "Budi",
      "username": "budi"
    }
  }
}
```

### Response Error
- **404 Not Found**:
  ```json
  { "status": "error", "message": "Artikel tidak ditemukan.", "data": null }
  ```
