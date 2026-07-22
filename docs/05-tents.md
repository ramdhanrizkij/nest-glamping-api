# Tents — Unit Tenda & Ketersediaan

## Overview

Tent (unit) adalah unit fisik tenda yang bisa dipesan. Setiap unit terikat ke satu tent type dan memiliki kode unik (misal: `Safari-01`) serta status (`available` atau `maintenance`).

## Flow CRUD

### Create Tent Unit (Admin)

```
Admin                           Server
  │                               │
  │  POST /api/v1/tents/units     │
  │  {tent_type_id,               │
  │   name_or_number, status}     │
  │──────────────────────────────►│
  │                               │─► Validasi tent_type_id exists
  │                               │─► Default status = "available"
  │                               │─► Create tent record
  │  201 {tent}                   │
  │◄──────────────────────────────│
```

### List Tent Units (Admin)

```
Admin                           Server
  │                               │
  │  GET /api/v1/tents/units      │
  │──────────────────────────────►│
  │                               │─► Fetch semua tents
  │  200 [{tent}, ...]            │
  │◄──────────────────────────────│
```

### List by Tent Type (Admin)

```
Admin                                       Server
  │                                           │
  │  GET /api/v1/tents/tent-types/:id/units   │
  │──────────────────────────────────────────►│
  │                                           │─► Fetch tents by tent_type_id
  │  200 [{tent}, ...]                        │
  │◄──────────────────────────────────────────│
```

### Update Tent Unit (Admin)

```
Admin                           Server
  │                               │
  │  PUT /api/v1/tents/units/:id  │
  │  {name_or_number?, status?}   │
  │──────────────────────────────►│
  │                               │─► Validasi tent exists
  │                               │─► Partial update
  │  200 {tent}                   │
  │◄──────────────────────────────│
```

### Delete Tent Unit (Admin)

```
Admin                           Server
  │                               │
  │  DELETE /api/v1/tents/units/:id
  │──────────────────────────────►│
  │                               │─► Cek keberadaan
  │                               │─► Soft delete
  │  200 {message: "deleted"}     │
  │◄──────────────────────────────│
```

## Logic

- **Create**: Validasi tent type ID harus ada di database. Status default `available` jika tidak disediakan.
- **Update**: Partial update — hanya field yang dikirim yang diubah.
- **Delete**: Soft delete.

## Status Tenda

| Status | Keterangan |
|--------|------------|
| `available` | Tersedia untuk dipesan |
| `maintenance` | Sedang dalam perawatan, tidak bisa dipesan |

---

## Availability Check

### Overview

Pengecekan ketersediaan menentukan unit tenda mana yang tersedia untuk tanggal tertentu, beserta harganya. Ini adalah endpoint public (tidak perlu auth).

### Flow

```
Client                          Server
  │                               │
  │  GET /api/v1/tents/           │
  │     tent-types/:id/           │
  │     availability              │
  │     ?check_in=2025-07-01      │
  │     &check_out=2025-07-03     │
  │──────────────────────────────►│
  │                               │─► Validasi tent type exists
  │                               │─► Query available tents:
  │                               │   1. Find tents dengan tent_type_id
  │                               │   2. Filter status = "available"
  │                               │   3. Exclude yang punya booking
  │                               │      overlap dengan tanggal request
  │                               │─► Untuk setiap available tent:
  │                               │   Hitung harga per malam
  │                               │   (dynamic pricing)
  │  200 [{available_tent}, ...]  │
  │◄──────────────────────────────│
```

### Query Logic — Availability

Sistem menggunakan raw SQL untuk mengecualikan tenda yang sudah terbooking:

```sql
SELECT t.*
FROM tents t
WHERE t.tent_type_id = ?
  AND t.status = 'available'
  AND t.deleted_at IS NULL
  AND t.id NOT IN (
    SELECT bt.tent_id
    FROM booking_tents bt
    JOIN bookings b ON b.id = bt.booking_id
    WHERE b.status IN ('pending', 'confirmed', 'paid')
      AND b.deleted_at IS NULL
      AND bt.deleted_at IS NULL
      AND b.check_in_date < ?   -- check_out request
      AND b.check_out_date > ?  -- check_in request
  )
```

**Penjelasan:**
- `b.check_in_date < check_out` — booking yang check-in sebelum request check-out
- `b.check_out_date > check_in` — booking yang check-out setelah request check-in
- Kondisi ini mendeteksi **overlap** antara tanggal booking existing dengan tanggal request

### Pricing Calculation

Untuk setiap tenda yang tersedia, sistem menghitung harga per malam:

```go
// Untuk setiap malam dalam range:
for date := checkIn; date < checkOut; date = date.AddDate(0,0,1) {
    // Cek apakah ada active rate yang cover tanggal ini
    for _, rate := range rates {
        if rate.StartDate <= date && rate.EndDate >= date && rate.IsActive {
            pricePerNight = rate.PricePerNight  // gunakan rate
            break
        }
    }
    if tidak ada rate {
        pricePerNight = tentType.BasePrice  // gunakan base price
    }
    total += pricePerNight
}

// Rata-rata per malam
avgPricePerNight = total / numberOfNights
```

### Response

```json
GET /api/v1/tents/tent-types/:id/availability?check_in=2025-07-01&check_out=2025-07-03

{
  "success": true,
  "message": "availability checked",
  "data": [
    {
      "id": "uuid-tent-1",
      "name_or_number": "Safari-01",
      "price_per_night": 850000
    },
    {
      "id": "uuid-tent-2",
      "name_or_number": "Safari-02",
      "price_per_night": 850000
    }
  ]
}
```

**Keterangan:**
- `price_per_night`: Rata-rata harga per malam untuk periode tersebut (sudah mempertimbangkan dynamic pricing)

## Endpoint Summary

| Method | Endpoint | Auth | Role | Deskripsi |
|--------|----------|------|------|-----------|
| `GET` | `/api/v1/tents/tent-types/:id/availability` | Tidak | - | Cek ketersediaan + harga |
| `GET` | `/api/v1/tents/units` | Ya | admin | List semua unit |
| `POST` | `/api/v1/tents/units` | Ya | admin | Buat unit |
| `GET` | `/api/v1/tents/units/:id` | Ya | admin | Detail unit |
| `PUT` | `/api/v1/tents/units/:id` | Ya | admin | Update unit |
| `DELETE` | `/api/v1/tents/units/:id` | Ya | admin | Hapus unit |
| `GET` | `/api/v1/tents/tent-types/:id/units` | Ya | admin | List unit by tipe |
