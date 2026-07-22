# Auth — Autentikasi & Otorisasi

## Flow Autentikasi

### Register

```
Client                          Server
  │                               │
  │  POST /api/v1/auth/register   │
  │  {name, email, password,phone}│
  │──────────────────────────────►│
  │                               │─► Validasi input
  │                               │─► Cek email uniqueness
  │                               │─► Hash password (bcrypt)
  │                               │─► Fetch role "customer"
  │                               │─► Create user
  │                               │─► Generate access token (15m)
  │                               │─► Generate refresh token (7d)
  │                               │─► Simpan refresh token ke DB
  │  201 {access_token,           │
  │        refresh_token, user}   │
  │◄──────────────────────────────│
```

### Login

```
Client                          Server
  │                               │
  │  POST /api/v1/auth/login      │
  │  {email, password}            │
  │──────────────────────────────►│
  │                               │─► Validasi input
  │                               │─► Find user by email
  │                               │─► Verify password hash
  │                               │─► Fetch role
  │                               │─► Generate access token
  │                               │─► Generate refresh token
  │                               │─► Simpan refresh token ke DB
  │  200 {access_token,           │
  │        refresh_token, user}   │
  │◄──────────────────────────────│
```

### Refresh Token (Token Rotation)

```
Client                          Server
  │                               │
  │  POST /api/v1/auth/refresh    │
  │  {refresh_token}              │
  │──────────────────────────────►│
  │                               │─► Find refresh token di DB
  │                               │─► Validasi JWT claims
  │                               │─► Cek expiry
  │                               │─► Hapus refresh token lama (rotation)
  │                               │─► Generate access token baru
  │                               │─► Generate refresh token baru
  │                               │─► Simpan refresh token baru ke DB
  │  200 {access_token,           │
  │        refresh_token, user}   │
  │◄──────────────────────────────│
```

### Logout

```
Client                          Server
  │                               │
  │  POST /api/v1/auth/logout     │
  │  {refresh_token}              │
  │──────────────────────────────►│
  │                               │─► Hapus refresh token dari DB
  │  200 {message: "logged out"}  │
  │◄──────────────────────────────│
```

## JWT Structure

### Access Token

```
Header:  { "alg": "HS256", "typ": "JWT" }
Payload: {
  "sub":     "user-uuid",
  "role_id": "role-uuid",
  "exp":     1234567890
}
Signature: HS256(secret)
```

- **Expiry**: 15 menit (default)
- **Disimpan di**: Client-side (localStorage/cookie)
- **Digunakan di**: Header `Authorization: Bearer {token}`

### Refresh Token

```
Header:  { "alg": "HS256", "typ": "JWT" }
Payload: {
  "sub": "user-uuid",
  "exp": 1234567890
}
Signature: HS256(refresh-secret)
```

- **Expiry**: 168 jam / 7 hari (default)
- **Disimpan di**: Database (`refresh_tokens` table) + Client
- **Digunakan di**: Untuk mendapatkan access token baru

## Token Rotation

Sistem menggunakan **refresh token rotation** — setiap kali refresh token digunakan, token lama dihapus dan token baru dibuat. Ini mencegah token yang sudah digunakan kembali (replay attack).

```
Refresh Token Flow:
1. Client kirim refresh_token ke server
2. Server validasi token di database
3. Server HAPUS refresh token lama
4. Server generate access_token + refresh_token BARU
5. Server simpan refresh_token BARU ke database
6. Client simpan token baru

Jika refresh token lama dikirim lagi:
→ Server tidak menemukan token di DB → 401 Unauthorized
```

## RBAC (Role-Based Access Control)

### Roles

| Role | Deskripsi | Permissions |
|------|-----------|-------------|
| `admin` | System administrator | manage_tents, manage_rates, manage_amenities, view_all_bookings, confirm_bookings, manage_users, manage_payments, view_reports |
| `customer` | Regular customer | create_booking, cancel_booking, view_own_bookings |

### Middleware Chain

```
Request
  │
  ▼
┌──────────────┐
│    CORS      │  → Allow cross-origin requests
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   Logger     │  → Log request method, path, status
└──────┬───────┘
       │
       ▼
┌──────────────┐
│  Auth (JWT)  │  → Extract & validate Bearer token
│              │  → Set userID & roleID di Fiber context
└──────┬───────┘
       │
       ▼
┌──────────────┐
│ RoleAllowed  │  → Compare roleID dengan role name
│ ("admin")    │  → Reject jika tidak match
└──────┬───────┘
       │
       ▼
┌──────────────┐
│   Handler    │  → Business logic
└──────────────┘
```

### Auth Middleware

```
1. Extract header "Authorization"
2. Split "Bearer {token}"
3. Validate JWT dengan secret
4. Extract claims (userID, roleID)
5. Set di Fiber context locals:
   - c.Locals("userID", claims.Sub)
   - c.Locals("roleID", claims.RoleID)
```

### RoleAllowed Middleware

```
1. Get roleID dari context ( hasil extract dari JWT)
2. Compare dengan role name parameter (misal "admin")
3. Jika match → lanjut ke handler
4. Jika tidak match → 403 Forbidden
```

## Endpoint Summary

| Endpoint | Auth | Role | Deskripsi |
|----------|------|------|-----------|
| `POST /api/v1/auth/register` | Tidak | - | Register customer baru |
| `POST /api/v1/auth/login` | Tidak | - | Login |
| `POST /api/v1/auth/refresh` | Tidak | - | Refresh access token |
| `POST /api/v1/auth/logout` | Tidak | - | Invalidate refresh token |

## Request/Response

### Register

**Request:**
```json
POST /api/v1/auth/register
{
  "name": "Budi Santoso",
  "email": "budi@example.com",
  "password": "password123",
  "phone": "081234567890"
}
```

**Response (201):**
```json
{
  "success": true,
  "message": "registration successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "Budi Santoso",
      "email": "budi@example.com",
      "phone_number": "081234567890",
      "role_id": "...",
      "role_name": "customer"
    }
  }
}
```

### Login

**Request:**
```json
POST /api/v1/auth/login
{
  "email": "admin@glamping.com",
  "password": "admin123"
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "login successful",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "...",
      "name": "Admin Utama",
      "email": "admin@glamping.com",
      "phone_number": "081234567890",
      "role_id": "...",
      "role_name": "admin"
    }
  }
}
```

### Refresh Token

**Request:**
```json
POST /api/v1/auth/refresh
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "token refreshed",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIs...",
    "user": { ... }
  }
}
```

### Logout

**Request:**
```json
POST /api/v1/auth/logout
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response (200):**
```json
{
  "success": true,
  "message": "logged out successfully"
}
```

## Error Responses

| Kondisi | Status | Error |
|---------|--------|-------|
| Email sudah terdaftar | 400 | `email already registered` |
| Email tidak ditemukan | 400 | `invalid email or password` |
| Password salah | 400 | `invalid email or password` |
| Refresh token tidak valid | 401 | `invalid refresh token` |
| Refresh token expired | 401 | `refresh token expired` |
| Access token expired | 401 | `invalid or expired token` |
| Tidak ada Authorization header | 401 | `missing or malformed token` |
