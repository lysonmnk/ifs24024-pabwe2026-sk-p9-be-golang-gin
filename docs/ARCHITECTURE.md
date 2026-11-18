# Architecture: Delcom Blog Backend

## System Overview
Client (Frontend ReactJS / Swagger / Postman) -> HTTP Gin Router -> Middleware (CORS, JWT, Role, RateLimiter) -> Feature Handler -> Feature Service -> Feature Repository -> GORM ORM -> PostgreSQL (Production) / SQLite (Testing)

## Tech Stack (TERKUNCI)
- Bahasa Backend: Go 1.27.4
- Framework HTTP: Gin v1.12.0
- ORM: GORM v1.31.2 (Driver PostgreSQL v1.6.3 & SQLite in-memory v1.11.0 untuk pengujian)
- Autentikasi: JWT (golang-jwt/jwt v5.3.1) & Hashing Password (golang.org/x/crypto/bcrypt)
- Konfigurasi: godotenv v1.5.1
- Dokumentasi API: Swag v1.16.6 & Gin Swagger v1.6.1
- Pengujian: Go testing & Testify v1.12.1

## Struktur Folder (Clean Architecture & Feature-First)
```
cmd/
└── server/
    └── main.go                 # Entry point server Gin
docs/                           # Swagger & OpenAPI spec
internal/
├── config/                     # Load environment config
├── database/                   # GORM connection, auto-migration, admin seeder
├── middleware/                 # AuthJWT, CORS, RequireRole, RateLimiter
├── response/                   # JSON response helper standar
├── features/
│   ├── auth/                   # Register, Login, Profile (DTO, Repo, Service, Handler)
│   ├── post/                   # Post CRUD, Slug, Reading Time, Stats (DTO, Repo, Service, Handler)
│   └── ai/                     # Search Engine Client & OpenAI Client (DTO, Repo, Service, Handler)
└── server/                     # Router setup, routes registration, Swagger handler, 404 handler
```

## Data Flow
HTTP Request -> Gin Engine -> Middleware -> Feature Handler (binding & validation) -> Feature Service (business logic) -> Feature Repository (database access via GORM) -> Standard JSON Response
