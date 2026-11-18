# Contract: Admin Endpoints
Versi: 1.0 | Pemilik: Data Master | Status: Disepakati

## 1. Get All Users
### Endpoint
`GET /api/admin/users` | Auth: ya (Bearer Token, Role: admin)

### Response Sukses (200 OK)
```json
{
  "status": "success",
  "message": "Daftar pengguna berhasil dimuat.",
  "data": [
    {
      "id": 1,
      "name": "Administrator",
      "username": "admin",
      "email": "admin@delcom.org",
      "role": "admin",
      "total_posts": 0,
      "created_at": "2026-11-18T00:00:00Z"
    },
    {
      "id": 2,
      "name": "Budi",
      "username": "budi",
      "email": "budi@contoh.com",
      "role": "author",
      "total_posts": 5,
      "created_at": "2026-11-18T08:00:00Z"
    }
  ]
}
```

### Response Error
- **403 Forbidden**:
  ```json
  { "status": "error", "message": "Akses ditolak. Memerlukan hak akses administrator.", "data": null }
  ```

---

## 2. Get All Posts Archive
### Endpoint
`GET /api/admin/posts` | Auth: ya (Bearer Token, Role: admin)

### Response Sukses (200 OK)
```json
{
  "status": "success",
  "message": "Arsip seluruh artikel berhasil dimuat.",
  "data": [
    {
      "id": 1,
      "user_id": 2,
      "title": "Tutorial Gin GORM",
      "slug": "tutorial-gin-gorm",
      "category": "Tutorial",
      "is_published": true,
      "author": {
        "id": 2,
        "name": "Budi",
        "username": "budi"
      },
      "created_at": "2026-11-18T08:00:00Z"
    }
  ]
}
```

---

## 3. Get User Posts by User ID
### Endpoint
`GET /api/admin/users/:id/posts` | Auth: ya (Bearer Token, Role: admin)

### Response Sukses (200 OK)
```json
{
  "status": "success",
  "message": "Daftar artikel pengguna berhasil dimuat.",
  "data": [
    {
      "id": 1,
      "user_id": 2,
      "title": "Tutorial Gin GORM",
      "slug": "tutorial-gin-gorm",
      "category": "Tutorial",
      "is_published": true,
      "created_at": "2026-11-18T08:00:00Z"
    }
  ]
}
```
