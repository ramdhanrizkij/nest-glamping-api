# Glamping Booking System — Development Tasks

## Phase 1: Foundation (✅ Done)

- [x] Project structure (feature-based layered architecture)
- [x] Config (env, app, database, logger)
- [x] Shared packages (middleware, response, validator, errors, pagination, auth context)
- [x] JWT package (generate + validate + refresh token)
- [x] Hash package (bcrypt)
- [x] Bootstrap (fiber, routes, dependency injection)
- [x] GORM migrate with gormigrate
- [x] Auth module (register, login, refresh, logout)
- [x] Users module (get profile, update profile)

---

## Phase 2: Tent Master Data (✅ Done)

### 2.1 Amenities

- [x] Entity `Amenity` (domain/entity.go)
- [x] Interface `AmenityRepository` (domain/repository.go)
- [x] Interface `AmenityService` (domain/service.go)
- [x] DTO: `CreateAmenityRequest`, `UpdateAmenityRequest`, `AmenityResponse` (dto/amenity.go)
- [x] Repository GORM (repository/postgres.go)
- [x] Usecase: Create, List, Update, Delete (usecase/amenity.go)
- [x] Handler + Routes (delivery/http/handler.go, routes.go)
- [x] Module wiring (module.go)
- [x] Register di dependency.go + routes.go

### 2.2 Tent Types

- [x] Entity `TentType` (domain/entity.go)
- [x] Entity `TentTypeImage` (domain/entity.go)
- [x] Entity `TentTypeAmenity` (domain/entity.go)
- [x] Interface `TentTypeRepository` (domain/repository.go)
- [x] Interface `TentTypeService` (domain/service.go)
- [x] DTO: `CreateTentTypeRequest`, `UpdateTentTypeRequest`, `TentTypeResponse` (dto/tent_type.go)
- [x] DTO: `TentTypeDetailResponse` (include images, amenities, rates)
- [x] Repository GORM (repository/postgres.go)
- [x] Usecase: Create, List, GetByID (detail), Update, Delete (usecase/tent_type.go)
- [x] Handler + Routes — admin endpoints (delivery/http/handler.go, routes.go)
- [x] Handler + Routes — public endpoints (list, detail)
- [x] Module wiring (module.go)

### 2.3 Tent Type Images

- [x] DTO: `AddImageRequest` (dto/image.go)
- [x] Usecase: Add image, Delete image, Set primary
- [x] Handler di tent_types handler (nested resource)
- [ ] Upload storage setup (pkg/storage/ — local or S3)

### 2.4 Tent Type Rates (Dynamic Pricing)

- [x] Entity `TentTypeRate` (sudah di schema)
- [x] DTO: `CreateRateRequest`, `UpdateRateRequest`, `RateResponse` (dto/rate.go)
- [x] Usecase: Create rate, List rates, Update, Delete, GetActiveRate(typeID, date)
- [x] Business rule: validasi overlap date range (tidak boleh ada 2 active rate yang overlap untuk type yang sama)
- [x] Handler di tent_types handler (nested resource)

### 2.5 Tents (Unit Fisik)

- [x] Entity `Tent` (domain/entity.go)
- [x] Interface `TentRepository` (domain/repository.go)
- [x] DTO: `CreateTentRequest`, `UpdateTentRequest`, `TentResponse` (dto/tent.go)
- [x] Repository GORM (repository/postgres.go)
- [x] Usecase: Create unit, List by type, Update status, Update name
- [x] Handler + Routes — admin (delivery/http/)
- [x] Module wiring (module.go)

### 2.6 Availability Check

- [x] Endpoint: `GET /tent-types/:id/availability?check_in=&check_out=`
- [x] Logic: cari unit available, exclude yang overlap booking
- [x] Logic: hitung harga per malam (dynamic pricing)
- [x] Response: `{available_tents: [{id, name, price_per_night}], total}`

### 2.7 Migrations

- [x] Migration: create_tent_types
- [x] Migration: create_tent_type_rates
- [x] Migration: create_tent_type_images
- [x] Migration: create_amenities
- [x] Migration: create_tent_type_amenities
- [x] Migration: create_tents

---

## Phase 3: Booking

### 3.1 Bookings

- [ ] Entity `Booking` (domain/entity.go)
- [ ] Entity `BookingTent` (domain/entity.go)
- [ ] Interface `BookingRepository` (domain/repository.go)
- [ ] Interface `BookingService` (domain/service.go)
- [ ] DTO: `CreateBookingRequest`, `BookingResponse`, `BookingDetailResponse` (dto/booking.go)
- [ ] Repository GORM (repository/postgres.go)
- [ ] Usecase: CreateBooking (dengan overbooking prevention + price snapshot + total calc)
- [ ] Usecase: ListMyBookings, GetBookingDetail
- [ ] Usecase: CancelBooking (hanya jika status=pending)
- [ ] Handler + Routes — customer (delivery/http/)
- [ ] Module wiring (module.go)

### 3.2 Booking Code Generator

- [ ] Helper di pkg/utils atau shared/utils: generate `GLP-XXXXXXXX`
- [ ] Pastikan unique (retry jika collision, walau kemungkinan kecil)

### 3.3 Overbooking Prevention

- [ ] Repository method: `FindAvailableTents(tentTypeID, checkIn, checkOut) → []uuid`
- [ ] SQL: exclude tent_id yang punya booking overlap
- [ ] Dipanggil di CreateBooking usecase (re-check saat submit, bukan hanya saat browse)

### 3.4 Price Calculation

- [ ] Service/hitung harga per malam berdasarkan dynamic pricing
- [ ] Loop setiap malam, cek active rate → kalau ada pakai rate, kalau tidak pakai base_price
- [ ] Snapshot ke `booking_tents.price_per_night`
- [ ] Total = SUM(malam × price_per_night)

### 3.5 Admin Booking Management

- [ ] Endpoint: `GET /admin/bookings` — list semua booking (filter by status)
- [ ] Endpoint: `PATCH /admin/bookings/:id/confirm` — manual confirm
- [ ] Middleware: RoleAllowed("admin")

### 3.6 Migrations

- [ ] Migration: create_bookings
- [ ] Migration: create_booking_tents

---

## Phase 4: Payments

### 4.1 Payments

- [ ] Entity `Payment` (domain/entity.go)
- [ ] Interface `PaymentRepository` (domain/repository.go)
- [ ] Interface `PaymentService` (domain/service.go)
- [ ] DTO: `PayRequest`, `PaymentResponse` (dto/payment.go)
- [ ] Repository GORM (repository/postgres.go)
- [ ] Usecase: InitiatePayment (create/update payment record)
- [ ] Usecase: HandleCallback (update payment + booking status)
- [ ] Handler + Routes (delivery/http/)
- [ ] Module wiring (module.go)

### 4.2 Payment Gateway Integration

- [ ] Pilih provider (Midtrans / Xendit / manual)
- [ ] Buat abstraction di pkg/payment_gateway/ (interface + implementation)
- [ ] Generate payment URL / token
- [ ] Verify callback signature

### 4.3 Booking Status Sync

- [ ] Payment paid → booking status = confirmed
- [ ] Payment failed → booking status = cancelled, release tent
- [ ] Payment refunded → booking status = cancelled

### 4.4 Migrations

- [ ] Migration: create_payments

---

## Phase 5: RBAC Enhancement

### 5.1 Permission-Based Middleware

- [ ] Ubah middleware `RoleAllowed` → `PermissionRequired(permissionName)`
- [ ] Query role_permissions untuk cek apakah role user punya permission tersebut
- [ ] Cache permission per role (in-memory atau Redis)

### 5.2 Seed Permissions

- [ ] Seeder: create_booking, cancel_booking, view_bookings, manage_tents, manage_rates, view_reports, manage_users, manage_payments
- [ ] Seeder: assign permissions ke role admin dan customer

### 5.3 Admin User Management

- [ ] Endpoint: `GET /admin/users` — list semua user
- [ ] Endpoint: `GET /admin/users/:id` — detail user
- [ ] Endpoint: `PUT /admin/users/:id` — update user (termasuk role)
- [ ] Endpoint: `DELETE /admin/users/:id` — soft delete user

---

## Phase 6: Polish & Production Readiness

### 6.1 Pagination

- [ ] Implementasi query params: `?page=1&per_page=10`
- [ ] Apply di semua list endpoints
- [ ] Response format: `{data: [], meta: {page, per_page, total, total_pages}}`

### 6.2 Filtering & Sorting

- [ ] Tent types: filter by capacity, price range
- [ ] Bookings: filter by status, date range
- [ ] Sorting: `?sort=created_at&order=desc`

### 6.3 Input Validation

- [ ] Review semua DTO, pastikan validate tag lengkap
- [ ] Custom validator: date format, date range (check_out > check_in), positive number

### 6.4 Error Messages

- [ ] Bahasa Indonesia untuk customer-facing errors
- [ ] Bahasa Inggris untuk internal/technical errors

### 6.5 Logging

- [ ] Structured logging (JSON)
- [ ] Log: request/response di level debug
- [ ] Log: error dengan stack trace
- [ ] Log: business events (booking_created, payment_success, dll)

### 6.6 Testing

- [ ] Unit test: usecase layer (mock repository)
- [ ] Integration test: repository layer (test DB)
- [ ] E2E test: HTTP endpoints

### 6.7 Documentation

- [ ] Swagger/OpenAPI spec
- [ ] Postman collection
- [ ] README dengan setup instructions

---

## Task Execution Order

```
Phase 2 (Tent Master Data):
  2.7 Migrations → 2.1 Amenities → 2.2 Tent Types → 2.3 Images →
  2.4 Rates → 2.5 Tents (units) → 2.6 Availability Check

Phase 3 (Booking):
  3.6 Migrations → 3.2 Booking Code → 3.3 Overbooking →
  3.4 Price Calculation → 3.1 Bookings → 3.5 Admin Management

Phase 4 (Payments):
  4.4 Migrations → 4.2 Payment Gateway → 4.1 Payments → 4.3 Status Sync

Phase 5 (RBAC):
  5.2 Seed → 5.1 Permission Middleware → 5.3 Admin Users

Phase 6 (Polish):
  Bisa dikerjakan paralel dengan Phase 3-5
```
