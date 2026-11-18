# Contract: Auth (Register, Login, Profile, Logout)
Versi: 1.0 | Pemilik: Data Master | Status: Disepakati

## 1. Register
### Endpoint
`POST /api/auth/register` | Auth: tidak | Rate limit: 10/menit per IP

### Request
```json
{
  "name": "string (wajib, min 2, maks 100)",
  "username": "string (opsional, min 3, maks 50)",
  "email": "string (wajib, format email, unik)",
  "password": "string (wajib, min 8)"
}
```

### Validasi
| field | aturan | pesan error |
|---|---|---|
| name | wajib, min 2 | Nama wajib diisi minimal 2 karakter. |
| email | wajib, email valid | Email wajib diisi dengan format valid. |
| password | wajib, min 8 | Password minimal 8 karakter. |

### Response Sukses (201 Created)
```json
{
  "status": "success",
  "message": "Registrasi berhasil.",
  "data": {
    "user": {
      "id": 1,
      "name": "Budi",
      "username": "budi",
      "email": "budi@contoh.com",
      "role": "author"
    }
  }
}
```

### Response Error
- **400 Bad Request / 422 Unprocessable Entity**:
  ```json
  { "status": "error", "message": "Email sudah terdaftar atau input tidak valid.", "data": null }
  ```

---

## 2. Login
### Endpoint
`POST /api/auth/login` | Auth: tidak | Rate limit: 5/menit per IP+email

### Request
```json
{
  "email": "user@contoh.com",
  "password": "rahasia123"
}
```

### Response Sukses (200 OK)
```json
{
  "status": "success",
  "message": "Login berhasil.",
  "data": {
    "token": "eyJhbGciOi...",
    "token_type": "Bearer",
    "user": {
      "id": 1,
      "name": "Budi",
      "username": "budi",
      "email": "user@contoh.com",
      "role": "author"
    }
  }
}
```

### Response Error
- **401 Unauthorized**:
  ```json
  { "status": "error", "message": "Email atau password salah.", "data": null }
  ```

---

## 3. Get Profile
### Endpoint
`GET /api/auth/profile` | Auth: ya (Bearer Token)

### Response Sukses (200 OK)
```json
{
  "status": "success",
  "message": "Profil berhasil dimuat.",
  "data": {
    "id": 1,
    "name": "Budi",
    "username": "budi",
    "email": "user@contoh.com",
    "role": "author"
  }
}
```

---

## 4. Update Profile
### Endpoint
`PUT /api/auth/profile` | Auth: ya (Bearer Token)

### Request
```json
{
  "name": "Budi Santoso",
  "username": "budisantoso"
}
```

### Response Sukses (200 OK)
```json
{
  "status": "success",
  "message": "Profil berhasil diperbarui.",
  "data": {
    "id": 1,
    "name": "Budi Santoso",
    "username": "budisantoso",
    "email": "user@contoh.com",
    "role": "author"
  }
}
```

---

## 5. Logout
### Endpoint
`POST /api/auth/logout` | Auth: ya (Bearer Token)

### Response Sukses (200 OK)
```json
{
  "status": "success",
  "message": "Logout berhasil.",
  "data": null
}
```
