# Glamping Booking API

Backend API untuk sistem pemesanan glamping online. Dibangun dengan Go, Fiber v3, GORM, PostgreSQL.

## Fitur

- **Auth**: Register, Login, Refresh Token, Logout (JWT)
- **Tent Types**: CRUD tipe tenda + gambar + amenities + rates (dynamic pricing)
- **Amenities**: CRUD fasilitas
- **Tents**: CRUD unit fisik tenda
- **Availability Check**: Cek ketersediaan + harga per malam
- **Bookings**: Buat booking, list, detail, cancel + overbooking prevention
- **Payments**: Proses pembayaran, callback handler
- **RBAC**: Role-Based Access Control (admin, customer) + permission-based middleware
- **Admin**: Manajemen user, booking, konfirmasi manual

## Tech Stack

- Go 1.21+
- Fiber v3
- GORM + PostgreSQL
- JWT (access + refresh token)
- gormigrate (versioned migration)

## Setup

### 1. Clone & Install

```bash
git clone <repo-url>
cd nest-glamping-api
go mod tidy
```

### 2. Konfigurasi Database

Copy `.env.example` ke `.env` dan sesuaikan:

```bash
cp .env .env.local
```

Edit `.env`:

```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=glamping_db
JWT_SECRET=your-secret-key
JWT_REFRESH_SECRET=your-refresh-secret
```

### 3. Migrate & Seed

```bash
make migrate        # Jalankan migrasi
make seed           # Seed dummy data
make seed-fresh     # Fresh: rollback → migrate → seed
```

### 4. Jalankan

```bash
make dev            # go run cmd/api/main.go
make build          # Build binary
make run            # Build & run
```

## API Endpoints

### Auth

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/auth/register` | Register customer |
| POST | `/api/v1/auth/login` | Login |
| POST | `/api/v1/auth/refresh` | Refresh token |
| POST | `/api/v1/auth/logout` | Logout |

### Users

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/users/profile` | JWT | Get profile |
| PUT | `/api/v1/users/profile` | JWT | Update profile |
| GET | `/api/v1/admin/users` | Admin | List users |
| GET | `/api/v1/admin/users/:id` | Admin | Detail user |
| PUT | `/api/v1/admin/users/:id` | Admin | Update user |
| DELETE | `/api/v1/admin/users/:id` | Admin | Delete user |

### Tent Types

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/tent-types` | - | List tipe tenda |
| GET | `/api/v1/tent-types/:id` | - | Detail tipe (images, amenities, rates) |
| POST | `/api/v1/tent-types` | Admin | Buat tipe |
| PUT | `/api/v1/tent-types/:id` | Admin | Update tipe |
| DELETE | `/api/v1/tent-types/:id` | Admin | Hapus tipe |
| POST | `/api/v1/tent-types/:id/images` | Admin | Tambah gambar |
| DELETE | `/api/v1/tent-types/:id/images/:imageId` | Admin | Hapus gambar |
| POST | `/api/v1/tent-types/:id/rates` | Admin | Tambah rate |
| PUT | `/api/v1/tent-types/:id/rates/:rateId` | Admin | Update rate |
| DELETE | `/api/v1/tent-types/:id/rates/:rateId` | Admin | Hapus rate |

### Amenities

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/amenities` | - | List amenities |
| POST | `/api/v1/amenities` | Admin | Buat amenity |
| PUT | `/api/v1/amenities/:id` | Admin | Update amenity |
| DELETE | `/api/v1/amenities/:id` | Admin | Hapus amenity |

### Tents

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/api/v1/tent-types/:id/availability` | - | Cek ketersediaan |
| GET | `/api/v1/tents/units` | Admin | List unit |
| POST | `/api/v1/tents/units` | Admin | Tambah unit |
| PUT | `/api/v1/tents/units/:id` | Admin | Update unit |
| DELETE | `/api/v1/tents/units/:id` | Admin | Hapus unit |

### Bookings

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/v1/bookings` | Customer | Buat booking |
| GET | `/api/v1/bookings` | Customer | List booking saya |
| GET | `/api/v1/bookings/:id` | Customer | Detail booking |
| PATCH | `/api/v1/bookings/:id/cancel` | Customer | Cancel booking |
| GET | `/api/v1/admin/bookings` | Admin | List semua booking |
| GET | `/api/v1/admin/bookings/:id` | Admin | Detail booking |
| PATCH | `/api/v1/admin/bookings/:id/confirm` | Admin | Konfirmasi booking |

### Payments

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| POST | `/api/v1/bookings/:id/pay` | Customer | Bayar booking |
| GET | `/api/v1/bookings/:id/payment` | Customer | Status pembayaran |
| POST | `/api/v1/payments/:id/callback` | - | Webhook callback |

## Default Accounts

| Role | Email | Password |
|------|-------|----------|
| Admin | admin@glamping.com | admin123 |
| Customer | budi@example.com | password123 |
| Customer | siti@example.com | password123 |
| Customer | andi@example.com | password123 |

## Struktur Project

```
cmd/
  api/main.go           → Entry point
  migrate/main.go       → CLI migration
config/                 → App, DB, Logger config
internal/
  bootstrap/            → Fiber, routes, dependency injection, migration, seed
  features/
    auth/               → Register, Login, Refresh, Logout
    users/              → Profile + Admin user management
    amenities/          → CRUD amenities
    tent-types/         → CRUD tipe tenda + images + rates
    tents/              → CRUD unit fisik + availability check
    bookings/           → Booking CRUD + overbooking prevention
    payments/           → Payment processing + callback
  shared/               → Middleware, response, validator, errors, pagination, auth
pkg/
  database/             → GORM connection
  jwt/                  → JWT generate + validate
  hash/                 → Bcrypt hash
  payment_gateway/      → Payment gateway abstraction
```

## License

MIT
