# Contract: Author Posts & AI Generator
Versi: 1.0 | Pemilik: Data Master | Status: Disepakati

## 1. Get Author Posts
### Endpoint
`GET /api/posts` | Auth: ya (Bearer Token)

### Response Sukses (200 OK)
```json
{
  "status": "success",
  "message": "Daftar artikel penulis berhasil dimuat.",
  "data": [
    {
      "id": 1,
      "user_id": 2,
      "title": "Tutorial Gin GORM",
      "slug": "tutorial-gin-gorm",
      "category": "Tutorial",
      "excerpt": "Ringkasan tutorial...",
      "reading_time": 2,
      "is_published": false,
      "published_at": null,
      "created_at": "2026-11-18T08:00:00Z"
    }
  ]
}
```

---

## 2. Get Author Post Stats
### Endpoint
`GET /api/posts/stats` | Auth: ya (Bearer Token)

### Response Sukses (200 OK)
```json
{
  "status": "success",
  "message": "Statistik artikel berhasil dimuat.",
  "data": {
    "total_posts": 10,
    "published_posts": 7,
    "draft_posts": 3
  }
}
```

---

## 3. Create Post
### Endpoint
`POST /api/posts` | Auth: ya (Bearer Token)

### Request
```json
{
  "title": "Belajar Golang dari Nol",
  "category": "Teknologi",
  "excerpt": "Pengantar bahasa Go...",
  "content": "# Konten Markdown...",
  "is_published": false
}
```

### Response Sukses (201 Created)
```json
{
  "status": "success",
  "message": "Artikel berhasil dibuat.",
  "data": {
    "id": 2,
    "user_id": 2,
    "title": "Belajar Golang dari Nol",
    "slug": "belajar-golang-dari-nol",
    "category": "Teknologi",
    "excerpt": "Pengantar bahasa Go...",
    "content": "# Konten Markdown...",
    "reading_time": 4,
    "is_published": false,
    "published_at": null,
    "created_at": "2026-11-18T08:15:00Z"
  }
}
```

---

## 4. Get Post by ID
### Endpoint
`GET /api/posts/:id` | Auth: ya (Bearer Token)

### Response Sukses (200 OK)
```json
{
  "status": "success",
  "message": "Detail artikel berhasil dimuat.",
  "data": {
    "id": 2,
    "user_id": 2,
    "title": "Belajar Golang dari Nol",
    "slug": "belajar-golang-dari-nol",
    "category": "Teknologi",
    "excerpt": "Pengantar bahasa Go...",
    "content": "# Konten Markdown...",
    "reading_time": 4,
    "is_published": false,
    "published_at": null,
    "created_at": "2026-11-18T08:15:00Z"
  }
}
```

---

## 5. Update Post
### Endpoint
`PUT /api/posts/:id` | Auth: ya (Bearer Token)

### Request
```json
{
  "title": "Belajar Golang dari Nol (Update)",
  "category": "Teknologi",
  "excerpt": "Pengantar bahasa Go terupdate...",
  "content": "# Konten Markdown baru...",
  "is_published": true
}
```

### Response Sukses (200 OK)
```json
{
  "status": "success",
  "message": "Artikel berhasil diperbarui.",
  "data": {
    "id": 2,
    "title": "Belajar Golang dari Nol (Update)",
    "slug": "belajar-golang-dari-nol-update",
    "category": "Teknologi",
    "excerpt": "Pengantar bahasa Go terupdate...",
    "reading_time": 4,
    "is_published": true
  }
}
```

---

## 6. Toggle Publish Status
### Endpoint
`PATCH /api/posts/:id/toggle-publish` | Auth: ya (Bearer Token)

### Response Sukses (200 OK)
```json
{
  "status": "success",
  "message": "Status publikasi berhasil diubah.",
  "data": {
    "id": 2,
    "is_published": true,
    "published_at": "2026-11-18T08:20:00Z"
  }
}
```

---

## 7. Delete Post
### Endpoint
`DELETE /api/posts/:id` | Auth: ya (Bearer Token)

### Response Sukses (200 OK)
```json
{
  "status": "success",
  "message": "Artikel berhasil dihapus.",
  "data": null
}
```

---

## 8. AI Generate Post
### Endpoint
`POST /api/posts/ai-generate` | Auth: ya (Bearer Token)

### Request
```json
{
  "topic": "Microservices dengan Go Gin",
  "category": "Teknologi",
  "include_web_search": true
}
```

### Response Sukses (200 OK)
```json
{
  "status": "success",
  "message": "Draf artikel AI berhasil digenerate.",
  "data": {
    "title": "Membangun Microservices Modern Menggunakan Go Gin",
    "category": "Teknologi",
    "excerpt": "Panduan praktis implementasi arsitektur microservices berbasis Go Gin.",
    "content": "# Membangun Microservices Modern\n\n...",
    "sources": [
      {
        "title": "Go Gin Documentation",
        "url": "https://gin-gonic.com",
        "snippet": "Gin is a web framework written in Go..."
      }
    ]
  }
}
```
