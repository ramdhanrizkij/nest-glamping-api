# Master Data — Amenities & Tent Types

## Amenities

### Overview

Amenitas adalah fasilitas yang tersedia di setiap tipe tenda. Setiap amenity memiliki nama, icon URL, dan deskripsi. Amenitas dihubungkan ke tent type many-to-many melalui tabel `tent_type_amenities`.

### Flow CRUD

#### Create Amenity (Admin)

```
Admin                           Server
  │                               │
  │  POST /api/v1/amenities       │
  │  {name, icon_url,description} │
  │──────────────────────────────►│
  │                               │─► Validasi input (name required)
  │                               │─► Generate UUID
  │                               │─► Simpan ke DB
  │  201 {amenity}                │
  │◄──────────────────────────────│
```

#### List Amenities (Public)

```
Client                          Server
  │                               │
  │  GET /api/v1/amenities        │
  │──────────────────────────────►│
  │                               │─► Fetch semua amenities
  │  200 [{amenity}, ...]         │
  │◄──────────────────────────────│
```

#### Update Amenity (Admin)

```
Admin                           Server
  │                               │
  │  PUT /api/v1/amenities/:id    │
  │  {name?, icon_url?, desc?}    │
  │──────────────────────────────►│
  │                               │─► Validasi ID
  │                               │─► Partial update (hanya field non-kosong)
  │  200 {amenity}                │
  │◄──────────────────────────────│
```

#### Delete Amenity (Admin)

```
Admin                           Server
  │                               │
  │  DELETE /api/v1/amenities/:id │
  │──────────────────────────────►│
  │                               │─► Cek keberadaan
  │                               │─► Soft delete
  │  200 {message: "deleted"}     │
  │◄──────────────────────────────│
```

### Logic

- **Create**: Generate UUID baru, map DTO ke domain entity, simpan ke database.
- **List**: Fetch semua amenities tanpa filter.
- **Update**: Find-then-patch — cari amenity by ID, update hanya field yang dikirim (partial update).
- **Delete**: Cek keberadaan dulu, lalu soft delete.

### Request/Response

**Create:**
```json
POST /api/v1/amenities
{
  "name": "WiFi Gratis",
  "icon_url": "/icons/wifi.png",
  "description": "Akses internet nirkabel"
}
```

**Response (201):**
```json
{
  "success": true,
  "message": "amenity created",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "WiFi Gratis",
    "icon_url": "/icons/wifi.png",
    "description": "Akses internet nirkabel"
  }
}
```

---

## Tent Types

### Overview

Tipe tenda merepresentasikan kategori/kelas tenda (misal: Safari Tent, Deluxe Camp). Setiap tipe memiliki nama, deskripsi, kapasitas tamu, dan base price. Tipe tenda juga memiliki:
- **Images** — gambar-gambar (dengan flag primary image)
- **Amenities** — fasilitas (many-to-many)
- **Rates** — harga dinamis berdasarkan periode tanggal

### Flow CRUD

#### Create Tent Type (Admin)

```
Admin                           Server
  │                               │
  │  POST /api/v1/tent-types      │
  │  {name, description,          │
  │   capacity, base_price,       │
  │   amenity_ids[]}              │
  │──────────────────────────────►│
  │                               │─► Validasi input
  │                               │─► Generate UUID
  │                               │─► Simpan tent type
  │                               │─► Jika amenity_ids ada:
  │                               │   Simpan ke tent_type_amenities
  │  201 {tent_type}              │
  │◄──────────────────────────────│
```

#### List Tent Types (Public)

```
Client                          Server
  │                               │
  │  GET /api/v1/tent-types       │
  │──────────────────────────────►│
  │                               │─► Fetch semua tent types
  │  200 [{tent_type}, ...]       │
  │◄──────────────────────────────│
```

#### Get Tent Type Detail (Public)

```
Client                          Server
  │                               │
  │  GET /api/v1/tent-types/:id   │
  │──────────────────────────────►│
  │                               │─► Fetch tent type by ID
  │                               │─► Load images
  │                               │─► Load amenities
  │                               │─► Load rates
  │  200 {tent_type_detail}       │
  │◄──────────────────────────────│
```

### Logic

- **Create**: Simpan tent type, lalu set amenities jika `amenity_ids` disediakan.
- **FindByID**: Enriched response — memuat images, amenities, dan rates secara bersamaan.
- **Update**: Partial update + replace amenities jika `amenity_ids` dikirim.
- **Delete**: Soft delete.

### Request/Response

**Create:**
```json
POST /api/v1/tent-types
{
  "name": "Safari Tent",
  "description": "Tenda safari premium dengan view pegunungan",
  "capacity": 2,
  "base_price": 750000,
  "amenity_ids": ["uuid-amenity-1", "uuid-amenity-2"]
}
```

**Detail Response (200):**
```json
{
  "success": true,
  "message": "tent type retrieved",
  "data": {
    "id": "...",
    "name": "Safari Tent",
    "description": "Tenda safari premium dengan view pegunungan",
    "capacity": 2,
    "base_price": 750000,
    "images": [
      { "id": "...", "image_url": "/images/safari-1.jpg", "is_primary": true },
      { "id": "...", "image_url": "/images/safari-2.jpg", "is_primary": false }
    ],
    "amenities": [
      { "amenity_id": "...", "name": "WiFi Gratis" },
      { "amenity_id": "...", "name": "AC" }
    ],
    "rates": [
      {
        "id": "...",
        "start_date": "2025-06-15",
        "end_date": "2025-06-30",
        "price_per_night": 950000,
        "description": "High Season Lebaran",
        "is_active": true
      }
    ]
  }
}
```

---

## Tent Type Images

### Flow

#### Add Image (Admin)

```
Admin                                    Server
  │                                        │
  │  POST /api/v1/tent-types/:id/images    │
  │  {image_url, is_primary}               │
  │───────────────────────────────────────►│
  │                                        │─► Validasi tent type exists
  │                                        │─► Jika is_primary=true:
  │                                        │   Clear existing primary images
  │                                        │─► Create image record
  │  201 {image}                           │
  │◄───────────────────────────────────────│
```

#### Set Primary Image (Admin)

```
Admin                                           Server
  │                                               │
  │  PUT /api/v1/tent-types/:id/images/:imgId/primary
  │──────────────────────────────────────────────►│
  │                                               │─► Clear ALL primary for this tent type
  │                                               │─► Set specified image as primary
  │  200 {message: "primary image set"}           │
  │◄──────────────────────────────────────────────│
```

### Logic

- **AddImage**: Jika `is_primary=true`, semua gambar existing yang `is_primary=true` di-clear dulu, lalu gambar baru di-set sebagai primary.
- **SetPrimaryImage**: Clear semua primary images untuk tent type, lalu set gambar yang dipilih sebagai primary.
- **DeleteImage**: Cek keberadaan, lalu soft delete.

---

## Tent Type Rates (Dynamic Pricing)

### Overview

Rates memungkinkan harga dinamis berdasarkan periode tanggal. Jika ada rate yang aktif untuk tanggal tertentu, harga rate tersebut meng-override `base_price`.

### Flow

#### Create Rate (Admin)

```
Admin                                       Server
  │                                           │
  │  POST /api/v1/tent-types/:id/rates        │
  │  {start_date, end_date,                   │
  │   price_per_night, description,           │
  │   is_active}                              │
  │──────────────────────────────────────────►│
  │                                           │─► Parse dates (YYYY-MM-DD)
  │                                           │─► Validasi end_date > start_date
  │                                           │─► Cek overlap dengan rate aktif lain:
  │                                           │   SELECT * FROM tent_type_rates
  │                                           │   WHERE tent_type_id = ?
  │                                           │   AND is_active = true
  │                                           │   AND start_date <= new_end_date
  │                                           │   AND end_date >= new_start_date
  │                                           │─► Jika overlap → reject
  │                                           │─► Jika tidak → create rate
  │  201 {rate}                               │
  │◄──────────────────────────────────────────│
```

### Logic

- **Overlap Validation**: Saat create/update rate, sistem mengecek apakah ada rate aktif lain yang periode tanggalnya tumpang tindih. Jika ya, operasi ditolak.
- **Date Parsing**: Format tanggal harus `YYYY-MM-DD`.
- **Update**: Partial update dengan overlap check (exclude rate yang sedang diupdate).
- **Delete**: Cek keberadaan, lalu soft delete.

### Dynamic Pricing Calculation

Ketika ada booking atau availability check, sistem menghitung harga per malam:

```
Untuk setiap malam dalam range check_in → check_out:
  1. Fetch semua active rates untuk tent type
  2. Cek apakah ada rate yang cover tanggal tersebut:
     - rate.start_date <= tanggal AND rate.end_date >= tanggal
  3. Jika ada → gunakan rate.price_per_night
  4. Jika tidak ada → gunakan tent_type.base_price

Total = Σ (price_per_night per malam × jumlah tenda)
```

**Contoh:**
```
Tipe: Safari Tent (base_price: 750.000)
Rate: High Season (15-30 Jun 2025, price: 950.000)

Booking: 14 Jun → 16 Jun (2 malam)
  Malam 1 (14 Jun): tidak ada rate → 750.000
  Malam 2 (15 Jun): high season rate → 950.000
  Total: 1.700.000

Booking: 15 Jun → 17 Jun (2 malam)
  Malam 1 (15 Jun): high season rate → 950.000
  Malam 2 (16 Jun): high season rate → 950.000
  Total: 1.900.000
```

### Request/Response

**Create Rate:**
```json
POST /api/v1/tent-types/:id/rates
{
  "start_date": "2025-06-15",
  "end_date": "2025-06-30",
  "price_per_night": 950000,
  "description": "High Season Lebaran",
  "is_active": true
}
```

**List Rates Response (200):**
```json
{
  "success": true,
  "message": "rates retrieved",
  "data": [
    {
      "id": "...",
      "start_date": "2025-06-15",
      "end_date": "2025-06-30",
      "price_per_night": 950000,
      "description": "High Season Lebaran",
      "is_active": true
    }
  ]
}
```

## Endpoint Summary

| Method | Endpoint | Auth | Role | Deskripsi |
|--------|----------|------|------|-----------|
| `GET` | `/api/v1/amenities` | Tidak | - | List amenities |
| `GET` | `/api/v1/amenities/:id` | Tidak | - | Detail amenity |
| `POST` | `/api/v1/amenities` | Ya | admin | Buat amenity |
| `PUT` | `/api/v1/amenities/:id` | Ya | admin | Update amenity |
| `DELETE` | `/api/v1/amenities/:id` | Ya | admin | Hapus amenity |
| `GET` | `/api/v1/tent-types` | Tidak | - | List tipe tenda |
| `GET` | `/api/v1/tent-types/:id` | Tidak | - | Detail tipe (dengan images, amenities, rates) |
| `POST` | `/api/v1/tent-types` | Ya | admin | Buat tipe tenda |
| `PUT` | `/api/v1/tent-types/:id` | Ya | admin | Update tipe tenda |
| `DELETE` | `/api/v1/tent-types/:id` | Ya | admin | Hapus tipe tenda |
| `POST` | `/api/v1/tent-types/:id/images` | Ya | admin | Tambah gambar |
| `DELETE` | `/api/v1/tent-types/:id/images/:imageId` | Ya | admin | Hapus gambar |
| `PUT` | `/api/v1/tent-types/:id/images/:imageId/primary` | Ya | admin | Set primary image |
| `GET` | `/api/v1/tent-types/:id/rates` | Tidak | - | List rates |
| `POST` | `/api/v1/tent-types/:id/rates` | Ya | admin | Buat rate |
| `PUT` | `/api/v1/tent-types/:id/rates/:rateId` | Ya | admin | Update rate |
| `DELETE` | `/api/v1/tent-types/:id/rates/:rateId` | Ya | admin | Hapus rate |
