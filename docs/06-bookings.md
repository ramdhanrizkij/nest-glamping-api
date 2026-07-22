# Bookings — Pemesanan

## Overview

Booking adalah inti dari sistem ini. Customer memilih tipe tenda, tanggal menginap, dan unit tenda yang tersedia. Sistem menghitung harga secara dinamis dan mencegah overbooking.

## Status Transitions

```
                    ┌─────────────┐
                    │   pending   │ ◄── Status awal saat booking dibuat
                    └──────┬──────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
              ▼            ▼            ▼
     ┌──────────────┐ ┌─────────┐ ┌───────────┐
     │  confirmed   │ │  paid   │ │ cancelled │
     │ (admin       │ │ (payment│ │ (customer │
     │  confirm)    │ │  callback)│ │  cancel)  │
     └──────┬───────┘ └─────────┘ └───────────┘
            │
            ▼
     ┌─────────────┐
     │  completed  │
     └─────────────┘
```

| Dari | Ke | Trigger | Syarat |
|------|-----|---------|--------|
| `pending` | `confirmed` | Admin konfirmasi | Booking status = pending |
| `pending` | `cancelled` | Customer cancel | Booking status = pending |
| `pending` | `paid` | Payment callback success | Payment status = paid |
| `confirmed` | `paid` | Payment callback success | Payment status = paid |
| `confirmed` | `cancelled` | Payment callback failed/expired | Payment status = failed/expired |

## Flow Create Booking

```
Customer                        Server
  │                               │
  │  POST /api/v1/bookings        │
  │  {tent_type_id,               │
  │   check_in_date,              │
  │   check_out_date,             │
  │   guest_count,                │
  │   special_requests,           │
  │   tent_ids[]}                 │
  │──────────────────────────────►│
  │                               │
  │                               │─── 1. Validasi Input ───┐
  │                               │   - Parse user UUID     │
  │                               │   - Parse tent type UUID│
  │                               │   - Parse dates         │
  │                               │   - check_out > check_in│
  │                               │   - guest_count >= 1    │
  │                               │                         │
  │                               │─── 2. Validasi Tipe ────┘
  │                               │   - Fetch tent type
  │                               │   - Cek exists
  │                               │   - guest_count <= capacity
  │                               │
  │                               │─── 3. Validasi Tenda ───┐
  │                               │   Untuk setiap tent_id: │
  │                               │   - Cek belongs ke type │
  │                               │   - Cek status available│
  │                               │                         │
  │                               │─── 4. Re-check Avail ───┘
  │                               │   (Race condition guard)
  │                               │   - Fetch available tents
  │                               │   - Pastikan semua tent_id
  │                               │     masih available
  │                               │
  │                               │─── 5. Hitung Harga ─────┐
  │                               │   Untuk setiap malam:    │
  │                               │   - Cek active rates     │
  │                               │   - Hitung price/night   │
  │                               │   - Akumulasi total      │
  │                               │                         │
  │                               │─── 6. Simpan Booking ───┘
  │                               │   - Generate booking code
  │                               │     (GLP-XXXXXXXX)
  │                               │   - Create booking record
  │                               │   - Create booking_tents
  │                               │     (dalam 1 transaksi)
  │                               │
  │  201 {booking}                │
  │◄──────────────────────────────│
```

## Logic Detail

### 1. Validasi Input

```go
- userUUID: valid UUID ✅
- tentTypeUUID: valid UUID ✅
- checkIn: parseable to YYYY-MM-DD ✅
- checkOut: after checkIn ✅
- guestCount: >= 1 ✅
- tentIDs: at least 1 ✅
```

### 2. Validasi Tipe Tenda

```go
tentType := findTentType(tentTypeUUID)
if tentType == nil → error "tent type not found"
if guestCount > tentType.Capacity → error "guest count exceeds capacity"
```

### 3. Validasi Tenda

```go
for each tentID in tentIDs:
    tent := findTent(tentID)
    if tent == nil → error "tent not found"
    if tent.TentTypeID != tentTypeUUID → error "tent does not belong to this type"
    if tent.Status != "available" → error "tent is not available"
```

### 4. Race Condition Guard

Setelah validasi awal, sistem **sekali lagi** mengecek ketersediaan:

```go
availableTents := findAvailableTents(tentTypeUUID, checkIn, checkOut)
for each tentID in tentIDs:
    if tentID not in availableTents → error "tent no longer available"
```

Ini mencegah skenario di mana 2 request bersamaan lolos validasi awal.

### 5. Dynamic Pricing

```go
totalAmount := 0.0
for date := checkIn; date < checkOut; date = date.AddDate(0,0,1) {
    pricePerNight := tentType.BasePrice

    // Cari active rate yang cover tanggal ini
    for _, rate := range rates {
        if rate.StartDate <= date && rate.EndDate >= date && rate.IsActive {
            pricePerNight = rate.PricePerNight
            break
        }
    }

    // Harga per malam × jumlah tenda
    totalAmount += pricePerNight * len(tentIDs)
}

// Simpan average price per night di booking_tents
avgPricePerNight := totalAmount / nights / len(tentIDs)
```

### 6. Simpan Booking

```go
// Dalam 1 transaksi database:
booking := Booking{
    BookingCode:    generateBookingCode(), // "GLP-XXXXXXXX"
    UserID:         userUUID,
    TentTypeID:     tentTypeUUID,
    CheckInDate:    checkIn,
    CheckOutDate:   checkOut,
    TotalAmount:    totalAmount,
    Status:         "pending",
    GuestCount:     guestCount,
    SpecialRequests: specialRequests,
}

bookingTents := []BookingTent{}
for _, tentID := range tentIDs {
    bookingTents = append(bookingTents, BookingTent{
        TentID:        tentID,
        PricePerNight: avgPricePerNight,
    })
}

db.CreateBooking(booking, bookingTents) // atomic
```

## Booking Code Format

Format: `GLP-XXXXXXXX` (8 karakter alphanumeric uppercase)

Contoh: `GLP-A3F8K2M1`

Di-generate menggunakan `utils.GenerateBookingCode()`.

## Flow Cancel Booking

```
Customer                        Server
  │                               │
  │  PATCH /api/v1/bookings/:id/cancel
  │──────────────────────────────►│
  │                               │─► Validasi booking exists
  │                               │─► Ownership check (user_id match)
  │                               │─► Cek status = "pending"
  │                               │─► Update status ke "cancelled"
  │  200 {message: "cancelled"}   │
  │◄──────────────────────────────│
```

**Syarat Cancel:**
- Booking harus milik user yang sama (ownership check)
- Status booking harus `pending`

## Flow Admin Confirm Booking

```
Admin                           Server
  │                               │
  │  PATCH /api/v1/admin/         │
  │     bookings/:id/confirm      │
  │──────────────────────────────►│
  │                               │─► Validasi booking exists
  │                               │─► Cek status = "pending"
  │                               │─► Update status ke "confirmed"
  │  200 {message: "confirmed"}   │
  │◄──────────────────────────────│
```

## Pagination

Booking list menggunakan **in-memory pagination**:

```go
// 1. Fetch semua booking (user atau all)
bookings := fetchBookings(userID) // atau fetchAll()

// 2. Hitung total
total := len(bookings)

// 3. Slice untuk page tertentu
start := (page - 1) * perPage
end := start + perPage
if end > total { end = total }
paginatedBookings := bookings[start:end]

// 4. Fetch tents untuk setiap booking di page
for i := range paginatedBookings {
    paginatedBookings[i].Tents = fetchBookingTents(paginatedBookings[i].ID)
}
```

## Response

**Create Booking (201):**
```json
{
  "success": true,
  "message": "booking created",
  "data": {
    "id": "...",
    "booking_code": "GLP-A3F8K2M1",
    "check_in_date": "2025-07-01",
    "check_out_date": "2025-07-03",
    "total_amount": 1700000,
    "status": "pending",
    "guest_count": 2,
    "special_requests": "Mohon view yang bagus",
    "tents": [
      {
        "id": "...",
        "tent_id": "...",
        "tent_name": "Safari-01",
        "price_per_night": 850000
      }
    ]
  }
}
```

**List Bookings (200):**
```json
{
  "success": true,
  "message": "bookings retrieved",
  "data": [
    {
      "id": "...",
      "booking_code": "GLP-A3F8K2M1",
      "check_in_date": "2025-07-01",
      "check_out_date": "2025-07-03",
      "total_amount": 1700000,
      "status": "pending",
      "guest_count": 2,
      "special_requests": "Mohon view yang bagus",
      "tents": [...]
    }
  ],
  "total": 5,
  "page": 1,
  "per_page": 10,
  "total_pages": 1
}
```

## Endpoint Summary

| Method | Endpoint | Auth | Role | Deskripsi |
|--------|----------|------|------|-----------|
| `POST` | `/api/v1/bookings` | Ya | any | Buat booking |
| `GET` | `/api/v1/bookings` | Ya | any | List booking saya (paginated) |
| `GET` | `/api/v1/bookings/:id` | Ya | any | Detail booking |
| `PATCH` | `/api/v1/bookings/:id/cancel` | Ya | any | Cancel booking |
| `GET` | `/api/v1/admin/bookings` | Ya | admin | List semua booking (paginated) |
| `GET` | `/api/v1/admin/bookings/:id` | Ya | admin | Detail booking (any) |
| `PATCH` | `/api/v1/admin/bookings/:id/confirm` | Ya | admin | Konfirmasi booking |

## Error Responses

| Kondisi | Status | Error |
|---------|--------|-------|
| Tent type tidak ditemukan | 404 | `tent type not found` |
| Guest count melebihi kapasitas | 400 | `guest count exceeds capacity` |
| Tent tidak milik tipe ini | 400 | `tent does not belong to this tent type` |
| Tent tidak available | 400 | `tent is not available` |
| Booking tidak ditemukan | 404 | `booking not found` |
| Bukan booking sendiri | 403 | `unauthorized` |
| Status bukan pending | 400 | `booking cannot be cancelled` |
| Tanggal tidak valid | 400 | `check_out must be after check_in` |
