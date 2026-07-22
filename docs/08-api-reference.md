# API Reference

Base URL: `http://localhost:3000/api/v1`

Semua response menggunakan format:
```json
{
  "success": true|false,
  "message": "...",
  "data": ...,
  "error": "..."
}
```

Swagger UI: `http://localhost:3000/swagger`

---

## Auth

### POST `/auth/register`

Register akun customer baru.

**Auth:** Tidak diperlukan

**Request Body:**
```json
{
  "name": "string (required, min:2, max:255)",
  "email": "string (required, email format)",
  "password": "string (required, min:6)",
  "phone": "string (optional, max:20)"
}
```

**Response 201:**
```json
{
  "success": true,
  "message": "registration successful",
  "data": {
    "access_token": "string",
    "refresh_token": "string",
    "user": {
      "id": "uuid",
      "name": "string",
      "email": "string",
      "phone_number": "string",
      "role_id": "uuid",
      "role_name": "string"
    }
  }
}
```

---

### POST `/auth/login`

Login dengan email dan password.

**Auth:** Tidak diperlukan

**Request Body:**
```json
{
  "email": "string (required)",
  "password": "string (required)"
}
```

**Response 200:**
```json
{
  "success": true,
  "message": "login successful",
  "data": {
    "access_token": "string",
    "refresh_token": "string",
    "user": { ... }
  }
}
```

---

### POST `/auth/refresh`

Refresh access token menggunakan refresh token.

**Auth:** Tidak diperlukan

**Request Body:**
```json
{
  "refresh_token": "string (required)"
}
```

**Response 200:**
```json
{
  "success": true,
  "message": "token refreshed",
  "data": {
    "access_token": "string (new)",
    "refresh_token": "string (new)",
    "user": { ... }
  }
}
```

---

### POST `/auth/logout`

Invalidate refresh token.

**Auth:** Tidak diperlukan

**Request Body:**
```json
{
  "refresh_token": "string (required)"
}
```

**Response 200:**
```json
{
  "success": true,
  "message": "logged out successfully"
}
```

---

## Users

### GET `/users/profile`

Get profile user yang sedang login.

**Auth:** Bearer Token

**Response 200:**
```json
{
  "success": true,
  "message": "profile retrieved",
  "data": {
    "id": "uuid",
    "name": "string",
    "email": "string",
    "phone_number": "string",
    "role_id": "uuid"
  }
}
```

---

### PUT `/users/profile`

Update profile (name, phone_number).

**Auth:** Bearer Token

**Request Body:**
```json
{
  "name": "string (optional, min:2, max:255)",
  "phone_number": "string (optional, max:20)"
}
```

**Response 200:**
```json
{
  "success": true,
  "message": "profile updated",
  "data": { "id": "...", "name": "...", ... }
}
```

---

## Admin — Users

### GET `/admin/users`

List semua users.

**Auth:** Bearer Token + Role: admin

**Response 200:**
```json
{
  "success": true,
  "message": "users retrieved",
  "data": [
    { "id": "...", "name": "...", "email": "...", "phone_number": "...", "role_id": "..." }
  ]
}
```

---

### GET `/admin/users/:id`

Get user by ID.

**Auth:** Bearer Token + Role: admin

**Params:** `id` (UUID)

**Response 200:** `{ "success": true, "data": { ... } }`

---

### PUT `/admin/users/:id`

Update user (name, phone_number, role_id).

**Auth:** Bearer Token + Role: admin

**Params:** `id` (UUID)

**Request Body:**
```json
{
  "name": "string (optional)",
  "phone_number": "string (optional)",
  "role_id": "uuid (optional)"
}
```

**Response 200:** `{ "success": true, "data": { ... } }`

---

### DELETE `/admin/users/:id`

Soft delete user.

**Auth:** Bearer Token + Role: admin

**Params:** `id` (UUID)

**Response 200:** `{ "success": true, "message": "user deleted" }`

---

## Amenities

### GET `/amenities`

List semua amenities.

**Auth:** Tidak diperlukan

**Response 200:**
```json
{
  "success": true,
  "message": "amenities retrieved",
  "data": [
    { "id": "...", "name": "WiFi Gratis", "icon_url": "/icons/wifi.png", "description": "..." }
  ]
}
```

---

### GET `/amenities/:id`

Detail amenity.

**Auth:** Tidak diperlukan

**Params:** `id` (UUID)

**Response 200:** `{ "success": true, "data": { "id": "...", "name": "...", ... } }`

---

### POST `/amenities`

Buat amenity baru.

**Auth:** Bearer Token + Role: admin

**Request Body:**
```json
{
  "name": "string (required, max:100)",
  "icon_url": "string (optional, max:255)",
  "description": "string (optional)"
}
```

**Response 201:** `{ "success": true, "message": "amenity created", "data": { ... } }`

---

### PUT `/amenities/:id`

Update amenity.

**Auth:** Bearer Token + Role: admin

**Params:** `id` (UUID)

**Request Body:**
```json
{
  "name": "string (optional)",
  "icon_url": "string (optional)",
  "description": "string (optional)"
}
```

**Response 200:** `{ "success": true, "message": "amenity updated", "data": { ... } }`

---

### DELETE `/amenities/:id`

Hapus amenity.

**Auth:** Bearer Token + Role: admin

**Params:** `id` (UUID)

**Response 200:** `{ "success": true, "message": "amenity deleted" }`

---

## Tent Types

### GET `/tent-types`

List semua tipe tenda.

**Auth:** Tidak diperlukan

**Response 200:**
```json
{
  "success": true,
  "message": "tent types retrieved",
  "data": [
    { "id": "...", "name": "Safari Tent", "description": "...", "capacity": 2, "base_price": 750000 }
  ]
}
```

---

### GET `/tent-types/:id`

Detail tipe tenda dengan images, amenities, dan rates.

**Auth:** Tidak diperlukan

**Params:** `id` (UUID)

**Response 200:**
```json
{
  "success": true,
  "message": "tent type retrieved",
  "data": {
    "id": "...",
    "name": "Safari Tent",
    "description": "...",
    "capacity": 2,
    "base_price": 750000,
    "images": [
      { "id": "...", "image_url": "/images/safari-1.jpg", "is_primary": true }
    ],
    "amenities": [
      { "amenity_id": "...", "name": "WiFi Gratis" }
    ],
    "rates": [
      { "id": "...", "start_date": "2025-06-15", "end_date": "2025-06-30", "price_per_night": 950000, "description": "High Season", "is_active": true }
    ]
  }
}
```

---

### POST `/tent-types`

Buat tipe tenda baru.

**Auth:** Bearer Token + Role: admin

**Request Body:**
```json
{
  "name": "string (required, max:255)",
  "description": "string (optional)",
  "capacity": "integer (required, min:1)",
  "base_price": "float (required, >0)",
  "amenity_ids": ["uuid", "..."] (optional)
}
```

**Response 201:** `{ "success": true, "message": "tent type created", "data": { ... } }`

---

### PUT `/tent-types/:id`

Update tipe tenda.

**Auth:** Bearer Token + Role: admin

**Params:** `id` (UUID)

**Request Body:** Semua field optional (partial update). `amenity_ids` akan replace jika dikirim.

**Response 200:** `{ "success": true, "message": "tent type updated", "data": { ... } }`

---

### DELETE `/tent-types/:id`

Hapus tipe tenda.

**Auth:** Bearer Token + Role: admin

**Params:** `id` (UUID)

**Response 200:** `{ "success": true, "message": "tent type deleted" }`

---

### POST `/tent-types/:id/images`

Tambah gambar ke tipe tenda.

**Auth:** Bearer Token + Role: admin

**Params:** `id` (UUID — tent type ID)

**Request Body:**
```json
{
  "image_url": "string (required, max:255)",
  "is_primary": "boolean (optional, default: false)"
}
```

**Response 201:** `{ "success": true, "message": "image added", "data": { "id": "...", "image_url": "...", "is_primary": true } }`

---

### DELETE `/tent-types/:id/images/:imageId`

Hapus gambar.

**Auth:** Bearer Token + Role: admin

**Params:** `id` (UUID), `imageId` (UUID)

**Response 200:** `{ "success": true, "message": "image deleted" }`

---

### PUT `/tent-types/:id/images/:imageId/primary`

Set gambar sebagai primary.

**Auth:** Bearer Token + Role: admin

**Params:** `id` (UUID), `imageId` (UUID)

**Response 200:** `{ "success": true, "message": "primary image set" }`

---

### GET `/tent-types/:id/rates`

List rates untuk tipe tenda.

**Auth:** Tidak diperlukan

**Params:** `id` (UUID)

**Response 200:**
```json
{
  "success": true,
  "message": "rates retrieved",
  "data": [
    { "id": "...", "start_date": "2025-06-15", "end_date": "2025-06-30", "price_per_night": 950000, "description": "High Season", "is_active": true }
  ]
}
```

---

### POST `/tent-types/:id/rates`

Buat dynamic rate.

**Auth:** Bearer Token + Role: admin

**Params:** `id` (UUID)

**Request Body:**
```json
{
  "start_date": "string (required, YYYY-MM-DD)",
  "end_date": "string (required, YYYY-MM-DD)",
  "price_per_night": "float (required, >0)",
  "description": "string (optional, max:255)",
  "is_active": "boolean (optional, default: true)"
}
```

**Response 201:** `{ "success": true, "message": "rate created", "data": { ... } }`

**Error:** 400 jika overlap dengan rate aktif lain.

---

### PUT `/tent-types/:id/rates/:rateId`

Update rate.

**Auth:** Bearer Token + Role: admin

**Params:** `id` (UUID), `rateId` (UUID)

**Request Body:** Semua field optional (partial update). Overlap check tetap berlaku.

**Response 200:** `{ "success": true, "message": "rate updated", "data": { ... } }`

---

### DELETE `/tent-types/:id/rates/:rateId`

Hapus rate.

**Auth:** Bearer Token + Role: admin

**Params:** `id` (UUID), `rateId` (UUID)

**Response 200:** `{ "success": true, "message": "rate deleted" }`

---

## Tents

### GET `/tents/tent-types/:id/availability`

Cek ketersediaan tenda + harga.

**Auth:** Tidak diperlukan

**Params:** `id` (UUID — tent type ID)

**Query:**
- `check_in` — tanggal check-in (YYYY-MM-DD, required)
- `check_out` — tanggal check-out (YYYY-MM-DD, required)

**Response 200:**
```json
{
  "success": true,
  "message": "availability checked",
  "data": [
    { "id": "...", "name_or_number": "Safari-01", "price_per_night": 850000 },
    { "id": "...", "name_or_number": "Safari-02", "price_per_night": 850000 }
  ]
}
```

---

### GET `/tents/units`

List semua unit tenda.

**Auth:** Bearer Token + Role: admin

**Response 200:**
```json
{
  "success": true,
  "message": "tents retrieved",
  "data": [
    { "id": "...", "tent_type_id": "...", "name_or_number": "Safari-01", "status": "available" }
  ]
}
```

---

### POST `/tents/units`

Buat unit tenda baru.

**Auth:** Bearer Token + Role: admin

**Request Body:**
```json
{
  "tent_type_id": "uuid (required)",
  "name_or_number": "string (required, max:100)",
  "status": "string (optional, oneof: available maintenance)"
}
```

**Response 201:** `{ "success": true, "message": "tent created", "data": { ... } }`

---

### GET `/tents/units/:id`

Detail unit tenda.

**Auth:** Bearer Token + Role: admin

**Params:** `id` (UUID)

**Response 200:** `{ "success": true, "data": { ... } }`

---

### PUT `/tents/units/:id`

Update unit tenda.

**Auth:** Bearer Token + Role: admin

**Params:** `id` (UUID)

**Request Body:**
```json
{
  "name_or_number": "string (optional)",
  "status": "string (optional, oneof: available maintenance)"
}
```

**Response 200:** `{ "success": true, "message": "tent updated", "data": { ... } }`

---

### DELETE `/tents/units/:id`

Hapus unit tenda.

**Auth:** Bearer Token + Role: admin

**Params:** `id` (UUID)

**Response 200:** `{ "success": true, "message": "tent deleted" }`

---

### GET `/tents/tent-types/:id/units`

List unit berdasarkan tipe.

**Auth:** Bearer Token + Role: admin

**Params:** `id` (UUID — tent type ID)

**Response 200:** `{ "success": true, "message": "tents retrieved", "data": [ ... ] }`

---

## Bookings

### POST `/bookings`

Buat booking baru.

**Auth:** Bearer Token

**Request Body:**
```json
{
  "tent_type_id": "uuid (required)",
  "check_in_date": "string (required, YYYY-MM-DD)",
  "check_out_date": "string (required, YYYY-MM-DD)",
  "guest_count": "integer (required, min:1)",
  "special_requests": "string (optional)",
  "tent_ids": ["uuid", "..."] (required, min:1)"
}
```

**Response 201:**
```json
{
  "success": true,
  "message": "booking created",
  "data": {
    "id": "uuid",
    "booking_code": "GLP-A3F8K2M1",
    "check_in_date": "2025-07-01",
    "check_out_date": "2025-07-03",
    "total_amount": 1700000,
    "status": "pending",
    "guest_count": 2,
    "special_requests": "...",
    "tents": [
      { "id": "...", "tent_id": "...", "tent_name": "Safari-01", "price_per_night": 850000 }
    ]
  }
}
```

---

### GET `/bookings`

List booking milik user yang login (paginated).

**Auth:** Bearer Token

**Query:**
- `page` — nomor halaman (default: 1)
- `per_page` — item per halaman (default: 10)

**Response 200:**
```json
{
  "success": true,
  "message": "bookings retrieved",
  "data": [ ... ],
  "total": 5,
  "page": 1,
  "per_page": 10,
  "total_pages": 1
}
```

---

### GET `/bookings/:id`

Detail booking.

**Auth:** Bearer Token (ownership check)

**Params:** `id` (UUID)

**Response 200:**
```json
{
  "success": true,
  "message": "booking retrieved",
  "data": {
    "id": "...",
    "booking_code": "GLP-A3F8K2M1",
    "check_in_date": "2025-07-01",
    "check_out_date": "2025-07-03",
    "total_amount": 1700000,
    "status": "pending",
    "guest_count": 2,
    "special_requests": "...",
    "created_at": "2025-07-01T08:00:00Z",
    "tents": [ ... ]
  }
}
```

---

### PATCH `/bookings/:id/cancel`

Cancel booking (hanya status pending).

**Auth:** Bearer Token (ownership check)

**Params:** `id` (UUID)

**Response 200:** `{ "success": true, "message": "booking cancelled" }`

---

### GET `/admin/bookings`

List semua booking (paginated, admin only).

**Auth:** Bearer Token + Role: admin

**Query:** `page`, `per_page`

**Response 200:** `{ "success": true, "data": [ ... ], "total": ..., ... }`

---

### GET `/admin/bookings/:id`

Detail booking (admin, tanpa ownership check).

**Auth:** Bearer Token + Role: admin

**Params:** `id` (UUID)

**Response 200:** `{ "success": true, "data": { ... } }`

---

### PATCH `/admin/bookings/:id/confirm`

Konfirmasi booking.

**Auth:** Bearer Token + Role: admin

**Params:** `id` (UUID)

**Response 200:** `{ "success": true, "message": "booking confirmed" }`

---

## Payments

### POST `/bookings/:id/pay`

Mulai pembayaran untuk booking.

**Auth:** Bearer Token (ownership check)

**Params:** `id` (UUID — booking ID)

**Request Body:**
```json
{
  "payment_method": "string (required)"
}
```

**Response 200:**
```json
{
  "success": true,
  "message": "payment initiated",
  "data": {
    "id": "uuid",
    "booking_id": "uuid",
    "amount": 1700000,
    "payment_method": "bank_transfer",
    "payment_status": "unpaid",
    "payment_url": "https://payment.example.com/pay/GLP-A3F8K2M1"
  }
}
```

---

### GET `/bookings/:id/payment`

Cek status pembayaran.

**Auth:** Bearer Token (ownership check)

**Params:** `id` (UUID — booking ID)

**Response 200:**
```json
{
  "success": true,
  "message": "payment retrieved",
  "data": {
    "id": "uuid",
    "booking_id": "uuid",
    "amount": 1700000,
    "payment_method": "bank_transfer",
    "payment_status": "paid",
    "payment_date": "2025-07-01T10:30:00Z",
    "gateway_ref": "TXN-123456"
  }
}
```

---

### POST `/payments/:id/callback`

Webhook callback dari payment gateway.

**Auth:** Tidak diperlukan (trusted source)

**Params:** `id` (UUID — payment ID)

**Request Body:**
```json
{
  "external_id": "string (required)",
  "status": "string (required, oneof: paid failed expired)",
  "transaction_id": "string (optional)"
}
```

**Response 200:** `{ "success": true, "message": "callback processed" }`

**Efek ke booking status:**
| Callback Status | Payment Status | Booking Status |
|-----------------|----------------|----------------|
| `paid` | `paid` | `paid` |
| `failed` | `failed` | `cancelled` |
| `expired` | `failed` | `cancelled` |

---

## Common Error Responses

| Status | Error | Keterangan |
|--------|-------|------------|
| 400 | `invalid request body` | Request body tidak valid |
| 400 | `email already registered` | Email sudah digunakan |
| 400 | `invalid email or password` | Email atau password salah |
| 400 | `guest count exceeds capacity` | Jumlah tamu melebihi kapasitas |
| 400 | `tent is not available` | Tenda tidak tersedia |
| 400 | `booking cannot be cancelled` | Status booking bukan pending |
| 400 | `booking cannot be paid` | Status booking bukan pending |
| 400 | `already paid` | Sudah dibayar |
| 400 | `overlapping rate exists` | Rate tumpang tindih |
| 401 | `invalid or expired token` | Token tidak valid/expired |
| 401 | `missing or malformed token` | Header Authorization tidak ada |
| 403 | `unauthorized` | Tidak punya akses |
| 404 | `not found` | Resource tidak ditemukan |
| 500 | `internal server error` | Error tak terduga |
