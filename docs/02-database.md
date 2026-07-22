# Database — Schema & Migrasi

## ERD (Entity Relationship Diagram)

```
┌─────────────┐       ┌──────────────────┐       ┌─────────────┐
│    roles     │       │ role_permissions  │       │ permissions  │
├─────────────┤       ├──────────────────┤       ├─────────────┤
│ id (UUID PK)│◄──┐   │ id (UUID PK)     │   ┌──►│ id (UUID PK)│
│ name        │   ├───│ role_id (UUID FK) │   │   │ name        │
│ description │   │   │ permission_id     │───┘   │ description │
│ created_at  │   │   │ created_at        │       │ created_at  │
│ updated_at  │   │   │ updated_at        │       │ updated_at  │
│ deleted_at  │   │   │ deleted_at        │       │ deleted_at  │
└─────────────┘   │   └──────────────────┘       └─────────────┘
                  │
                  │   ┌─────────────┐
                  │   │    users     │
                  │   ├─────────────┤
                  └───│ role_id     │
                      │ id (UUID PK)│
                      │ name        │
                      │ email       │
                      │ password_hash│
                      │ phone_number│
                      │ created_at  │
                      │ updated_at  │
                      │ deleted_at  │
                      └──────┬──────┘
                             │
                ┌────────────┼────────────┐
                │            │            │
                ▼            ▼            ▼
    ┌───────────────┐ ┌─────────────┐ ┌──────────────────┐
    │refresh_tokens │ │  bookings   │ │                   │
    ├───────────────┤ ├─────────────┤ │                   │
    │ id (UUID PK)  │ │ id (UUID PK)│ │                   │
    │ user_id (FK)  │ │ user_id(FK) │ │                   │
    │ token         │ │ booking_code│ │                   │
    │ expires_at    │ │ check_in    │ │                   │
    │ created_at    │ │ check_out   │ │                   │
    └───────────────┘ │ total_amount│ │                   │
                      │ status      │ │                   │
                      │ guest_count │ │                   │
                      │ special_req │ │                   │
                      │ created_at  │ │                   │
                      │ updated_at  │ │                   │
                      │ deleted_at  │ │                   │
                      └──────┬──────┘ │                   │
                             │        │                   │
                ┌────────────┤        │                   │
                │            │        │                   │
                ▼            ▼        │                   │
    ┌────────────────┐ ┌─────────────┐│                   │
    │ booking_tents  │ │   payments  ││                   │
    ├────────────────┤ ├─────────────┤│                   │
    │ id (UUID PK)   │ │ id (UUID PK)│                   │
    │ booking_id(FK) │ │ booking_id  │                   │
    │ tent_id (FK)   │ │ amount      │                   │
    │ price_per_night│ │ payment_method                  │
    │ created_at     │ │ payment_status                  │
    │ updated_at     │ │ payment_date│                   │
    │ deleted_at     │ │ gateway_ref │                   │
    └───────┬────────┘ │ created_at  │                   │
            │          │ updated_at  │                   │
            │          │ deleted_at  │                   │
            │          └─────────────┘                   │
            │                                            │
            ▼                                            │
    ┌─────────────┐       ┌──────────────────┐          │
    │    tents     │       │   tent_types      │          │
    ├─────────────┤       ├──────────────────┤          │
    │ id (UUID PK)│       │ id (UUID PK)     │          │
    │ tent_type_id│──────►│ name             │          │
    │ code        │       │ description      │          │
    │ status      │       │ capacity         │          │
    │ created_at  │       │ base_price       │          │
    │ updated_at  │       │ created_at       │          │
    │ deleted_at  │       │ updated_at       │          │
    └─────────────┘       │ deleted_at       │          │
                          └────────┬─────────┘          │
                                   │                    │
                    ┌──────────────┼──────────┐         │
                    │              │          │         │
                    ▼              ▼          ▼         │
        ┌────────────────┐ ┌─────────────┐ ┌──────────────────┐
        │tent_type_images│ │tent_type_   │ │tent_type_rates   │
        ├────────────────┤ │amenities    │ ├──────────────────┤
        │ id (UUID PK)   │ ├─────────────┤ │ id (UUID PK)     │
        │ tent_type_id   │ │ id (UUID PK)│ │ tent_type_id(FK) │
        │ image_url      │ │ tent_type_id│ │ start_date       │
        │ is_primary     │ │ amenity_id  │ │ end_date         │
        │ created_at     │ │ created_at  │ │ price_per_night  │
        │ updated_at     │ │ updated_at  │ │ description      │
        │ deleted_at     │ │ deleted_at  │ │ is_active        │
        └────────────────┘ └──────┬──────┘ │ created_at       │
                                  │        │ updated_at       │
                                  ▼        │ deleted_at       │
                          ┌─────────────┐  └──────────────────┘
                          │  amenities   │
                          ├─────────────┤
                          │ id (UUID PK)│
                          │ name        │
                          │ icon_url    │
                          │ description │
                          │ created_at  │
                          │ updated_at  │
                          │ deleted_at  │
                          └─────────────┘
```

## Tabel & Kolom

### 1. `roles`

| Kolom | Tipe | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK, default `gen_random_uuid()` | — |
| `name` | varchar(50) | UNIQUE, NOT NULL | `admin`, `customer` |
| `description` | text | — | — |
| `created_at` | timestamp | default `CURRENT_TIMESTAMP` | — |
| `updated_at` | timestamp | default `CURRENT_TIMESTAMP` | — |
| `deleted_at` | timestamp | NULL, indexed | Soft delete |

### 2. `permissions`

| Kolom | Tipe | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | — |
| `name` | varchar(100) | UNIQUE, NOT NULL | Contoh: `create_booking` |
| `description` | text | — | — |
| `created_at` | timestamp | default `CURRENT_TIMESTAMP` | — |
| `updated_at` | timestamp | default `CURRENT_TIMESTAMP` | — |
| `deleted_at` | timestamp | NULL, indexed | Soft delete |

### 3. `role_permissions`

| Kolom | Tipe | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | — |
| `role_id` | UUID | FK → `roles.id` | — |
| `permission_id` | UUID | FK → `permissions.id` | — |
| `created_at` | timestamp | default `CURRENT_TIMESTAMP` | — |
| `updated_at` | timestamp | default `CURRENT_TIMESTAMP` | — |
| `deleted_at` | timestamp | NULL, indexed | Soft delete |

### 4. `users`

| Kolom | Tipe | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | — |
| `role_id` | UUID | FK → `roles.id`, NOT NULL | — |
| `name` | varchar(255) | NOT NULL | — |
| `email` | varchar(255) | UNIQUE, NOT NULL | — |
| `password_hash` | varchar(255) | NOT NULL | Bcrypt hash |
| `phone_number` | varchar(20) | — | — |
| `created_at` | timestamp | default `CURRENT_TIMESTAMP` | — |
| `updated_at` | timestamp | default `CURRENT_TIMESTAMP` | — |
| `deleted_at` | timestamp | NULL, indexed | Soft delete |

### 5. `refresh_tokens`

| Kolom | Tipe | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | — |
| `user_id` | UUID | FK → `users.id`, indexed | — |
| `token` | varchar(500) | UNIQUE, NOT NULL | JWT refresh token |
| `expires_at` | timestamp | NOT NULL | — |
| `created_at` | timestamp | default `CURRENT_TIMESTAMP` | — |

### 6. `tent_types`

| Kolom | Tipe | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | — |
| `name` | varchar(255) | NOT NULL | — |
| `description` | text | — | — |
| `capacity` | int | NOT NULL, default 2 | Maks tamu |
| `base_price` | decimal(12,2) | NOT NULL | Harga reguler/malam |
| `created_at` | timestamp | default `CURRENT_TIMESTAMP` | — |
| `updated_at` | timestamp | default `CURRENT_TIMESTAMP` | — |
| `deleted_at` | timestamp | NULL, indexed | Soft delete |

### 7. `tent_type_rates`

| Kolom | Tipe | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | — |
| `tent_type_id` | UUID | FK → `tent_types.id`, NOT NULL | — |
| `start_date` | date | NOT NULL | Mulai periode |
| `end_date` | date | NOT NULL | Akhir periode |
| `price_per_night` | decimal(12,2) | NOT NULL | Harga khusus/malam |
| `description` | varchar(255) | — | Contoh: "High Season Lebaran" |
| `is_active` | boolean | default true | — |
| `created_at` | timestamp | default `CURRENT_TIMESTAMP` | — |
| `updated_at` | timestamp | default `CURRENT_TIMESTAMP` | — |
| `deleted_at` | timestamp | NULL, indexed | Soft delete |

### 8. `tent_type_images`

| Kolom | Tipe | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | — |
| `tent_type_id` | UUID | FK → `tent_types.id`, NOT NULL | — |
| `image_url` | varchar(255) | NOT NULL | — |
| `is_primary` | boolean | default false | Cover image |
| `created_at` | timestamp | default `CURRENT_TIMESTAMP` | — |
| `updated_at` | timestamp | default `CURRENT_TIMESTAMP` | — |
| `deleted_at` | timestamp | NULL, indexed | Soft delete |

### 9. `amenities`

| Kolom | Tipe | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | — |
| `name` | varchar(100) | NOT NULL | — |
| `icon_url` | varchar(255) | — | — |
| `description` | text | — | — |
| `created_at` | timestamp | default `CURRENT_TIMESTAMP` | — |
| `updated_at` | timestamp | default `CURRENT_TIMESTAMP` | — |
| `deleted_at` | timestamp | NULL, indexed | Soft delete |

### 10. `tent_type_amenities`

| Kolom | Tipe | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | — |
| `tent_type_id` | UUID | FK → `tent_types.id`, NOT NULL | — |
| `amenity_id` | UUID | FK → `amenities.id`, NOT NULL | — |
| `created_at` | timestamp | default `CURRENT_TIMESTAMP` | — |
| `updated_at` | timestamp | default `CURRENT_TIMESTAMP` | — |
| `deleted_at` | timestamp | NULL, indexed | Soft delete |

### 11. `tents`

| Kolom | Tipe | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | — |
| `tent_type_id` | UUID | FK → `tent_types.id`, NOT NULL | — |
| `code` | varchar(100) | NOT NULL | Contoh: `Safari-01` |
| `status` | varchar(50) | default `available` | `available`, `maintenance` |
| `created_at` | timestamp | default `CURRENT_TIMESTAMP` | — |
| `updated_at` | timestamp | default `CURRENT_TIMESTAMP` | — |
| `deleted_at` | timestamp | NULL, indexed | Soft delete |

### 12. `bookings`

| Kolom | Tipe | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | — |
| `user_id` | UUID | FK → `users.id`, NOT NULL | — |
| `booking_code` | varchar(20) | UNIQUE, NOT NULL | Format: `GLP-XXXXXXXX` |
| `check_in_date` | date | NOT NULL | — |
| `check_out_date` | date | NOT NULL | — |
| `total_amount` | decimal(12,2) | NOT NULL | Total harga |
| `status` | varchar(50) | default `pending` | `pending`, `confirmed`, `paid`, `cancelled`, `completed` |
| `guest_count` | int | NOT NULL, default 1 | — |
| `special_requests` | text | — | — |
| `created_at` | timestamp | default `CURRENT_TIMESTAMP` | — |
| `updated_at` | timestamp | default `CURRENT_TIMESTAMP` | — |
| `deleted_at` | timestamp | NULL, indexed | Soft delete |

### 13. `booking_tents`

| Kolom | Tipe | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | — |
| `booking_id` | UUID | FK → `bookings.id`, NOT NULL | — |
| `tent_id` | UUID | FK → `tents.id`, NOT NULL | — |
| `price_per_night` | decimal(12,2) | NOT NULL | Harga yang dikunci saat pesan |
| `created_at` | timestamp | default `CURRENT_TIMESTAMP` | — |
| `updated_at` | timestamp | default `CURRENT_TIMESTAMP` | — |
| `deleted_at` | timestamp | NULL, indexed | Soft delete |

### 14. `payments`

| Kolom | Tipe | Constraint | Keterangan |
|-------|------|------------|------------|
| `id` | UUID | PK | — |
| `booking_id` | UUID | FK → `bookings.id`, NOT NULL | — |
| `amount` | decimal(12,2) | NOT NULL | — |
| `payment_method` | varchar(100) | — | — |
| `payment_status` | varchar(50) | default `unpaid` | `unpaid`, `paid`, `failed` |
| `payment_date` | timestamp | NULL | — |
| `gateway_ref` | varchar(255) | — | Referensi dari payment gateway |
| `created_at` | timestamp | default `CURRENT_TIMESTAMP` | — |
| `updated_at` | timestamp | default `CURRENT_TIMESTAMP` | — |
| `deleted_at` | timestamp | NULL, indexed | Soft delete |

## Relasi

```
roles            1 ──── ∞  users
roles            1 ──── ∞  role_permissions     ∞ ──── 1  permissions
users            1 ──── ∞  refresh_tokens
users            1 ──── ∞  bookings
tent_types       1 ──── ∞  tent_type_rates
tent_types       1 ──── ∞  tent_type_images
tent_types       1 ──── ∞  tent_type_amenities  ∞ ──── 1  amenities
tent_types       1 ──── ∞  tents
bookings         1 ──── ∞  booking_tents        ∞ ──── 1  tents
bookings         1 ──── ∞  payments
```

## Migrasi

Migrasi menggunakan `gormigrate` dengan 13 migration files berurutan:

| Order | ID | File | Tabel |
|-------|-----|------|-------|
| 1 | `20250721001` | `create_roles` | `roles` |
| 2 | `20250721002` | `create_permissions` | `permissions`, `role_permissions` |
| 3 | `20250721003` | `create_users` | `users` |
| 4 | `20250721004` | `create_refresh_tokens` | `refresh_tokens` |
| 5 | `20250721005` | `create_tent_types` | `tent_types` |
| 6 | `20250721006` | `create_tent_type_rates` | `tent_type_rates` |
| 7 | `20250721007` | `create_tent_type_images` | `tent_type_images` |
| 8 | `20250721008` | `create_amenities` | `amenities` |
| 9 | `20250721009` | `create_tent_type_amenities` | `tent_type_amenities` |
| 10 | `20250721010` | `create_tents` | `tents` |
| 11 | `20250721011` | `create_bookings` | `bookings` |
| 12 | `20250721012` | `create_booking_tents` | `booking_tents` |
| 13 | `20250721013` | `create_payments` | `payments` |

Perintah:

```bash
make migrate              # Jalankan migrasi
make migrate-rollback     # Rollback 1 migration terakhir
make migrate-rollback-all # Rollback semua
```

## Seed Data

### Roles

| Name | Description |
|------|-------------|
| `admin` | System administrator |
| `customer` | Regular customer |

### Permissions (11)

| Permission | Deskripsi |
|------------|-----------|
| `create_booking` | Membuat booking (customer) |
| `cancel_booking` | Membatalkan booking (customer) |
| `view_own_bookings` | Melihat booking sendiri (customer) |
| `manage_tents` | Mengelola tenda (admin) |
| `manage_rates` | Mengelola rates (admin) |
| `manage_amenities` | Mengelola amenities (admin) |
| `view_all_bookings` | Melihat semua booking (admin) |
| `confirm_bookings` | Mengkonfirmasi booking (admin) |
| `manage_users` | Mengelola user (admin) |
| `manage_payments` | Mengelola pembayaran (admin) |
| `view_reports` | Melihat laporan (admin) |

### Default Users

| Role | Name | Email | Password |
|------|------|-------|----------|
| admin | Admin Utama | admin@glamping.com | admin123 |
| customer | Budi Santoso | budi@example.com | password123 |
| customer | Siti Rahayu | siti@example.com | password123 |
| customer | Andi Wijaya | andi@example.com | password123 |

### Amenities (8)

| Nama | Icon |
|------|------|
| WiFi Gratis | /icons/wifi.png |
| AC | /icons/ac.png |
| Kamar Mandi Dalam | /icons/bathroom.png |
| BBQ Area | /icons/bbq.png |
| Spot Foto | /icons/camera.png |
| Breakfast | /icons/breakfast.png |
| Bonfire | /icons/fire.png |
| Parking | /icons/parking.png |

### Tent Types (4)

| Nama | Kapasitas | Base Price |
|------|-----------|------------|
| Safari Tent | 2 | Rp 750.000 |
| Deluxe Camp | 4 | Rp 1.200.000 |
| Family Glamp | 6 | Rp 1.800.000 |
| Couple's Nest | 2 | Rp 1.500.000 |

### Tent Units (11)

| Code | Tipe | Status |
|------|------|--------|
| Safari-01 | Safari Tent | available |
| Safari-02 | Safari Tent | available |
| Safari-03 | Safari Tent | maintenance |
| Deluxe-01 | Deluxe Camp | available |
| Deluxe-02 | Deluxe Camp | available |
| Deluxe-03 | Deluxe Camp | available |
| Deluxe-04 | Deluxe Camp | available |
| Family-01 | Family Glamp | available |
| Family-02 | Family Glamp | available |
| Couple-01 | Couple's Nest | available |
| Couple-02 | Couple's Nest | available |

### Dynamic Rates (contoh)

| Tipe | Periode | Harga/Malam | Keterangan |
|------|---------|-------------|------------|
| Safari Tent | 15-30 Jun 2025 | Rp 950.000 | High Season Lebaran |
| Deluxe Camp | 15-30 Jun 2025 | Rp 1.500.000 | High Season Lebaran |
| Safari Tent | 20 Des 2025 - 5 Jan 2026 | Rp 650.000 | Promo Tahun Baru |
| Family Glamp | (dynamic) | Rp 2.000.000 | Weekend Spesial |
