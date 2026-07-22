# Overview — Glamping Booking API

Backend REST API untuk sistem pemesanan glamping online. Mengelola tipe tenda, unit tenda, fasilitas, booking dengan dynamic pricing, dan pembayaran.

## Tech Stack

| Komponen | Teknologi |
|----------|-----------|
| Bahasa | Go 1.26+ |
| HTTP Framework | Fiber v3 |
| ORM | GORM |
| Database | PostgreSQL |
| Autentikasi | JWT (HS256) — access + refresh token |
| Migrasi | gormigrate (versioned migration) |
| Validasi | go-playground/validator |

## Arsitektur

Menggunakan **Clean Architecture** dengan pembagian yang jelas antara domain, data access, business logic, dan delivery layer.

```
cmd/
  api/main.go                 → Entry point aplikasi
  migrate/main.go             → CLI untuk migrasi & seed
config/                       → Konfigurasi (App, DB, Logger, Env)
internal/
  bootstrap/                  → Wiring: Fiber, routes, DI, migration, seed
  features/                   → Fitur-fitur aplikasi (self-contained modules)
    auth/                     → Autentikasi (register, login, refresh, logout)
    users/                    → Manajemen user (profile + admin)
    amenities/                → CRUD fasilitas
    tent-types/               → Tipe tenda + gambar + rates + amenities
    tents/                    → Unit fisik tenda + availability check
    bookings/                 → Booking + dynamic pricing + overbooking prevention
    payments/                 → Pembayaran + callback handler
  shared/                     → Middleware, response, validator, errors, pagination
pkg/
  database/                   → Koneksi PostgreSQL
  jwt/                        → JWT generate & validate
  hash/                       → Bcrypt hash
  payment_gateway/            → Payment gateway abstraction (pluggable)
```

## Struktur Module

Setiap fitur di `internal/features/` mengikuti pola yang sama:

```
features/{feature}/
  module.go                   → Module struct + RegisterRoutes (DI wiring)
  domain/
    entity.go                 → GORM entity models
    repository.go             → Repository interface
    service.go                → Service (usecase) interface
  dto/
    {name}.go                 → Request/Response DTOs
  repository/
    postgres.go               → Repository implementation (GORM)
  usecase/
    {name}.go                 → Business logic
  delivery/http/
    handler.go                → HTTP handlers
    routes.go                 → Route registration
```

## Dependency Graph

```
auth ──────────> users (userRepo untuk registrasi/login)
bookings ──────> tent-types (tentTypeRepo untuk capacity & rates)
bookings ──────> tents (tentRepo untuk validasi availability)
payments ──────> bookings (bookingRepo untuk status management)
payments ──────> PaymentGateway (external package)
tents ─────────> tent-types (tentTypeRepo untuk rates di availability check)
```

## Environment Variables

| Variable | Default | Deskripsi |
|----------|---------|-----------|
| `APP_NAME` | Glamping API | Nama aplikasi |
| `APP_PORT` | 3000 | Port server |
| `APP_ENV` | development | Environment |
| `DB_HOST` | localhost | Database host |
| `DB_PORT` | 5432 | Database port |
| `DB_USER` | postgres | Database user |
| `DB_PASSWORD` | postgres | Database password |
| `DB_NAME` | glamping_db | Database name |
| `DB_SSLMODE` | disable | SSL mode |
| `DB_TIMEZONE` | Asia/Jakarta | Timezone |
| `LOG_LEVEL` | info | Log level |
| `JWT_SECRET` | secret | JWT access token secret |
| `JWT_EXPIRY` | 15m | Access token expiry |
| `JWT_REFRESH_SECRET` | refresh-secret | Refresh token secret |
| `JWT_REFRESH_EXPIRY` | 168h | Refresh token expiry (7 hari) |

## Setup & Menjalankan

```bash
# 1. Clone & install
git clone <repo-url>
cd nest-glamping-api
go mod tidy

# 2. Konfigurasi
cp .env.example .env
# Edit .env sesuai kebutuhan

# 3. Migrate & Seed
make migrate
make seed

# 4. Jalankan
make dev
```

## API Documentation

Swagger UI tersedia di: `http://localhost:3000/swagger`

Untuk regenerate dokumentasi Swagger setelah perubahan annotasi:

```bash
swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal
```
