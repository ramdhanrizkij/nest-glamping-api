# Glamping Booking System — Technical Specification

## 1. Overview

Sistem pemesanan glamping online untuk mengelola reservasi tenda. Berbeda dari hotel biasa, unit fisik tenda terikat langsung pada tipe tenda. Sistem mendukung RBAC, dynamic pricing, dan pencegahan overbooking.

**Stack:** Go, Fiber v3, GORM, PostgreSQL, JWT (access + refresh token)

---

## 2. Arsitektur

Feature-based layered architecture. Setiap feature punya layer sendiri:

```
delivery (HTTP handler) → usecase (business logic) → domain (entity + interface) → repository (DB access)
```

Dependency injection via `bootstrap/dependency.go`. Setiap feature di-Expose sebagai `Module` yang me-register route.

---

## 3. Database Schema

### 3.1 RBAC

- **roles** — admin, customer
- **permissions** — granular action (create_booking, view_reports, dll)
- **role_permissions** — junction table many-to-many

### 3.2 User

- **users** — punya `role_id` FK ke roles. Password di-hash (bcrypt).

### 3.3 Tent Master Data

| Tabel | Fungsi |
|---|---|
| `tent_types` | Kategori tenda (Safari, Deluxe, dll). Punya `base_price` dan `capacity` |
| `tent_type_rates` | Dynamic pricing. Override `base_price` untuk rentang tanggal tertentu (high season, promo, dll). Ada `is_active` flag |
| `tent_type_images` | Galeri foto per tipe. Satu image bisa jadi `is_primary` (cover) |
| `amenities` | Fasilitas (WiFi, BBQ, dll) |
| `tent_type_amenities` | Junction many-to-many tipe tenda ↔ amenities |
| `tents` | Unit fisik tenda. Terikat ke `tent_type_id`. Punya `status` (available, maintenance) |

### 3.4 Booking & Payment

| Tabel | Fungsi |
|---|---|
| `bookings` | Header booking. Code unik, tanggal check-in/out, total amount, status lifecycle |
| `booking_tents` | Detail tenda yang dipesan. `price_per_night` di-snapshot saat booking dibuat (bukan ambil dari master, supaya harga tidak berubah) |
| `payments` | Record pembayaran per booking. Status: unpaid → paid / failed / refunded |

---

## 4. Business Rules

### 4.1 Dynamic Pricing

1. Saat customer pilih tipe tenda + tanggal, sistem cek `tent_type_rates` yang `is_active=true` dan range-nya overlap dengan tanggal booking
2. Jika ada rate yang cocok → pakai `price_per_night` dari rate tersebut
3. Jika tidak ada → pakai `base_price` dari `tent_types`
4. Harga ini di-snapshot ke `booking_tents.price_per_night` saat booking dibuat

### 4.2 Pencegahan Overbooking

Sistem tidak boleh mengizinkan booking yang overlap pada tenda yang sama di tanggal yang sama. Algoritma:

```
1. User pilih tent_type + check_in + check_out
2. Cari semua tent_id dari tipe tersebut yang status='available'
3. Dari daftar tent_id, exclude yang punya booking overlap:
   SELECT tent_id FROM booking_tents bt
   JOIN bookings b ON b.id = bt.booking_id
   WHERE bt.tent_id IN (list_tersedia)
   AND b.status NOT IN ('cancelled', 'completed')
   AND b.check_in_date < :check_out
   AND b.check_out_date > :check_in
4. Sisa tent_id = ketersediaan aktual
5. Jika kosong → tolak booking
```

### 4.3 Booking Lifecycle

```
pending → confirmed → completed
    └→ cancelled
```

- **pending**: Booking dibuat, menunggu pembayaran
- **confirmed**: Pembayaran berhasil (payment_status = paid)
- **completed**: Setelah check-out selesai
- **cancelled**: Dibatalkan oleh user atau admin, atau payment gagal/expired

### 4.4 Harga Snapshot

Saat booking dibuat, `total_amount` dihitung:

```
total_amount = SUM(hari × price_per_night) untuk setiap booking_tent
```

`price_per_night` di-snapshot dari perhitungan dynamic pricing, bukan dari FK. Ini memastikan harga tidak berubah setelah booking dibuat.

---

## 5. API Endpoints

### 5.1 Auth (✅ Done)

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| POST | `/api/v1/auth/register` | - | Register customer baru |
| POST | `/api/v1/auth/login` | - | Login, return access + refresh token |
| POST | `/api/v1/auth/refresh` | - | Refresh access token |
| POST | `/api/v1/auth/logout` | - | Hapus refresh token |

### 5.2 Users (✅ Done)

| Method | Endpoint | Auth | Description |
|---|---|---|---|
| GET | `/api/v1/users/profile` | JWT | Get profile user login |
| PUT | `/api/v1/users/profile` | JWT | Update profile |

### 5.3 Tents (Master Data)

| Method | Endpoint | Auth | Role | Description |
|---|---|---|---|---|
| GET | `/api/v1/tent-types` | - | - | List semua tipe tenda (public) |
| GET | `/api/v1/tent-types/:id` | - | - | Detail tipe + images + amenities + rates |
| POST | `/api/v1/tent-types` | JWT | admin | Buat tipe tenda baru |
| PUT | `/api/v1/tent-types/:id` | JWT | admin | Update tipe tenda |
| DELETE | `/api/v1/tent-types/:id` | JWT | admin | Soft delete tipe tenda |
| POST | `/api/v1/tent-types/:id/images` | JWT | admin | Upload gambar tipe |
| DELETE | `/api/v1/tent-types/:id/images/:imageId` | JWT | admin | Hapus gambar |
| POST | `/api/v1/tent-types/:id/rates` | JWT | admin | Tambah rate (dynamic pricing) |
| PUT | `/api/v1/tent-types/:id/rates/:rateId` | JWT | admin | Update rate |
| DELETE | `/api/v1/tent-types/:id/rates/:rateId` | JWT | admin | Hapus rate |
| GET | `/api/v1/amenities` | - | - | List amenities |
| POST | `/api/v1/amenities` | JWT | admin | Buat amenity |
| GET | `/api/v1/tents` | JWT | admin | List semua unit tenda |
| POST | `/api/v1/tents` | JWT | admin | Tambah unit tenda |
| PUT | `/api/v1/tents/:id` | JWT | admin | Update unit (status, nama) |
| GET | `/api/v1/tent-types/:id/availability` | - | - | Cek ketersediaan tipe untuk tanggal tertentu |

### 5.4 Bookings

| Method | Endpoint | Auth | Role | Description |
|---|---|---|---|---|
| POST | `/api/v1/bookings` | JWT | customer | Buat booking baru |
| GET | `/api/v1/bookings` | JWT | customer | List booking user login |
| GET | `/api/v1/bookings/:id` | JWT | customer/admin | Detail booking |
| PATCH | `/api/v1/bookings/:id/cancel` | JWT | customer | Cancel booking (hanya pending) |
| GET | `/api/v1/admin/bookings` | JWT | admin | List semua booking |
| PATCH | `/api/v1/admin/bookings/:id/confirm` | JWT | admin | Manual confirm booking |

### 5.5 Payments

| Method | Endpoint | Auth | Role | Description |
|---|---|---|---|---|
| POST | `/api/v1/bookings/:id/pay` | JWT | customer | Proses pembayaran |
| GET | `/api/v1/bookings/:id/payment` | JWT | customer | Status pembayaran |
| POST | `/api/v1/payments/:id/callback` | - | - | Webhook dari payment gateway |

---

## 6. Flow Detail

### 6.1 Flow Register & Login

```
1. POST /auth/register {name, email, password, phone}
   → Hash password (bcrypt)
   → Cari role "customer" di tabel roles
   → INSERT ke users
   → Generate access_token (15m) + refresh_token (7d)
   → Simpan refresh_token ke tabel refresh_tokens
   → Return {access_token, refresh_token, user}

2. POST /auth/login {email, password}
   → Cari user by email
   → Verify password (bcrypt compare)
   → Generate tokens
   → Return {access_token, refresh_token, user}

3. POST /auth/refresh {refresh_token}
   → Validasi refresh_token di DB (belum expired)
   → Validasi JWT signature
   → Hapus refresh_token lama
   → Generate pair baru
   → Return {access_token, refresh_token}
```

### 6.2 Flow Booking (Core)

```
1. Customer browse tent types
   GET /tent-types → list tipe dengan base_price

2. Customer pilih tipe + cek ketersediaan
   GET /tent-types/:id/availability?check_in=2025-08-01&check_out=2025-08-03
   → Hitung harga per malam (dynamic pricing)
   → Cari unit tersedia (overbooking check)
   → Return {available_tents: [{id, name, price_per_night}], total}

3. Customer buat booking
   POST /bookings {
     tent_type_id, check_in_date, check_out_date,
     guest_count, special_requests,
     tent_ids: [id1, id2]  // dari hasil availability check
   }

   Server-side:
   a. Validasi input
   b. Re-check ketersediaan (prevent race condition)
   c. Hitung harga per tenda (dynamic pricing → snapshot)
   d. Hitung total_amount = SUM(hari × price_per_night)
   e. Generate booking_code (GLP-XXXXXXXX)
   f. INSERT bookings (status=pending)
   g. INSERT booking_tents (dengan price_per_night snapshot)
   h. INSERT payments (status=unpaid)
   i. Return booking + payment info

4. Customer bayar
   POST /bookings/:id/pay {payment_method}
   → Integrasi payment gateway (Midtrans, Xendit, dll)
   → Return redirect URL / payment token

5. Payment callback (webhook)
   POST /payments/:id/callback {status, transaction_id}
   → Validasi signature
   → Update payment status (paid/failed)
   → Jika paid → update booking status = confirmed
   → Jika failed → update booking status = cancelled, release tent

6. Admin bisa confirm manual jika diperlukan
   PATCH /admin/bookings/:id/confirm
```

### 6.3 Flow Cek Ketersediaan (Anti-Overbooking)

```
Input: tent_type_id, check_in_date, check_out_date

Step 1: Cari semua unit tenda dari tipe tersebut
  SELECT * FROM tents WHERE tent_type_id = ? AND status = 'available'

Step 2: Exclude unit yang sudah dibooking pada rentang tanggal tsb
  SELECT DISTINCT bt.tent_id FROM booking_tents bt
  JOIN bookings b ON b.id = bt.booking_id
  WHERE bt.tent_id IN (dari_step_1)
  AND b.status NOT IN ('cancelled', 'completed')
  AND b.check_in_date < :check_out
  AND b.check_out_date > :check_in

Step 3: Unit tersedia = Step 1 - Step 2

Step 4: Hitung harga untuk setiap unit tersedia
  - Cek tent_type_rates yang overlap dengan tanggal booking
  - Jika ada → pakai rate tsb
  - Jika tidak → pakai base_price
```

### 6.4 Flow Dynamic Pricing

```
Input: tent_type_id, check_in_date, check_out_date

1. Ambil base_price dari tent_types
2. Ambil semua tent_type_rates yang:
   - tent_type_id = input
   - is_active = true
   - (start_date <= check_out AND end_date >= check_in)  // overlap

3. Untuk setiap malam dalam range check_in..check_out-1:
   - Cek apakah malam tsb masuk dalam range rate manapun
   - Jika ya → pakai price_per_night dari rate tersebut
   - Jika tidak → pakai base_price

4. Return harga per malam (bisa beda-beda per malam)
```

---

## 7. Module Dependency Graph

```
auth ─────→ users (butuh userRepo untuk login/register)
users ────→ (standalone)
tents ────→ (standalone)
bookings ─→ tents (butuh cek ketersediaan + harga)
payments ─→ bookings (update status booking)
```

---

## 8. Error Handling

Semua error dikembalikan dalam format konsisten:

```json
{
  "success": false,
  "error": "message"
}
```

HTTP codes: 400 (bad request), 401 (unauthorized), 403 (forbidden), 404 (not found), 409 (conflict), 500 (internal).

---

## 9. Seed Data

Saat pertama kali migrate:
- Role: `admin`, `customer`
- Admin default: admin@glamping.com / admin123 (buat via seeder)

---

## 10. Security

- Password: bcrypt hash
- Token: JWT HS256, access 15m, refresh 7d
- Refresh token disimpan di DB (bisa di-revoke)
- RBAC check via middleware `RoleAllowed(roleIDs...)`
- Soft delete (deleted_at) untuk data integrity
- Input validation di setiap handler

---

## 11. File Mapping (per feature)

Setiap feature mengikuti struktur:

```
internal/features/<feature>/
  domain/
    entity.go      → GORM struct
    repository.go  → interface repository
    service.go     → interface usecase
  repository/
    postgres.go    → implementasi GORM
  usecase/
    <action>.go    → business logic
  dto/
    <action>.go    → request/response struct
  delivery/http/
    handler.go     → HTTP handlers
    routes.go      → route registration
  module.go        → DI wiring
```
