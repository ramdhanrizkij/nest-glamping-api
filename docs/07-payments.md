# Payments — Pembayaran

## Overview

Payment gateway bersifat **pluggable** — interface `PaymentGateway` bisa diimplementasikan untuk berbagai provider (Midtrans, Xendit, manual, dll). Saat ini menggunakan `ManualGateway` sebagai stub.

## Payment Gateway Interface

```go
type PaymentGateway interface {
    CreatePayment(bookingCode string, amount float64, method string) (*PaymentResult, error)
}

type PaymentResult struct {
    ExternalID  string  // ID dari gateway
    PaymentURL  string  // URL halaman pembayaran
    Status      string  // status awal
}
```

## Flow Initiate Payment

```
Customer                        Server
  │                               │
  │  POST /api/v1/bookings/:id/pay│
  │  {payment_method}             │
  │──────────────────────────────►│
  │                               │─► Validasi booking exists
  │                               │─► Ownership check
  │                               │─► Cek booking status = "pending"
  │                               │─► Find or create payment record
  │                               │─► Cek payment belum "paid"
  │                               │─► Gateway.CreatePayment(
  │                               │     bookingCode,
  │                               │     totalAmount,
  │                               │     method)
  │  200 {payment}                │
  │◄──────────────────────────────│
```

## Flow Get Payment Status

```
Customer                        Server
  │                               │
  │  GET /api/v1/bookings/:id/payment
  │──────────────────────────────►│
  │                               │─► Validasi booking exists
  │                               │─► Ownership check
  │                               │─► Fetch payment by booking_id
  │  200 {payment}                │
  │◄──────────────────────────────│
```

## Flow Callback (Webhook)

```
Payment Gateway                  Server
  │                               │
  │  POST /api/v1/payments/:id/callback
  │  {external_id, status,        │
  │   transaction_id}             │
  │──────────────────────────────►│
  │                               │─► Parse callback payload
  │                               │─► Switch on status:
  │                               │
  │                               │   case "paid":
  │                               │     → Update payment status = "paid"
  │                               │     → Update payment gateway_ref
  │                               │     → Update booking status = "paid"
  │                               │
  │                               │   case "failed":
  │                               │     → Update payment status = "failed"
  │                               │     → Update booking status = "cancelled"
  │                               │
  │                               │   case "expired":
  │                               │     → Update payment status = "failed"
  │                               │     → Update booking status = "cancelled"
  │                               │
  │  200 {message: "processed"}   │
  │◄──────────────────────────────│
```

## Logic Detail

### Initiate Payment

```go
// 1. Validasi
booking := findBookingByID(bookingID)
if booking == nil → error "booking not found"
if booking.UserID != userID → error "unauthorized"
if booking.Status != "pending" → error "booking cannot be paid"

// 2. Find or create payment
payment := findPaymentByBookingID(bookingID)
if payment == nil {
    payment = createPayment(bookingID, booking.TotalAmount, method)
}
if payment.PaymentStatus == "paid" → error "already paid"

// 3. Call gateway
result := gateway.CreatePayment(booking.BookingCode, booking.TotalAmount, method)

// 4. Update payment dengan gateway response
updatePayment(payment.ID, result)
```

### Callback Handler

```go
switch callback.Status {
case "paid":
    updatePaymentStatus(callback.ExternalID, "paid", callback.TransactionID)
    updateBookingStatus(bookingID, "paid")

case "failed":
    updatePaymentStatus(callback.ExternalID, "failed", "")
    updateBookingStatus(bookingID, "cancelled")

case "expired":
    updatePaymentStatus(callback.ExternalID, "failed", "")
    updateBookingStatus(bookingID, "cancelled")
}
```

### ManualGateway (Stub)

Saat ini menggunakan stub yang mengembalikan fake data:

```go
func (g *ManualGateway) CreatePayment(bookingCode string, amount float64, method string) (*PaymentResult, error) {
    return &PaymentResult{
        ExternalID:  "MANUAL-" + bookingCode,
        PaymentURL:  "https://payment.example.com/pay/" + bookingCode,
        Status:      "pending",
    }, nil
}
```

Untuk integrasi payment gateway nyata, implementasikan interface `PaymentGateway` dan registrasi di `bootstrap/dependency.go`.

## Response

**Initiate Payment (200):**
```json
{
  "success": true,
  "message": "payment initiated",
  "data": {
    "id": "...",
    "booking_id": "...",
    "amount": 1700000,
    "payment_method": "bank_transfer",
    "payment_status": "unpaid",
    "payment_url": "https://payment.example.com/pay/GLP-A3F8K2M1"
  }
}
```

**Get Payment Status (200):**
```json
{
  "success": true,
  "message": "payment retrieved",
  "data": {
    "id": "...",
    "booking_id": "...",
    "amount": 1700000,
    "payment_method": "bank_transfer",
    "payment_status": "paid",
    "payment_date": "2025-07-01T10:30:00Z",
    "gateway_ref": "TXN-123456"
  }
}
```

## Payment Status

| Status | Keterangan |
|--------|------------|
| `unpaid` | Belum dibayar |
| `paid` | Sudah dibayar |
| `failed` | Pembayaran gagal |

## Endpoint Summary

| Method | Endpoint | Auth | Role | Deskripsi |
|--------|----------|------|------|-----------|
| `POST` | `/api/v1/bookings/:id/pay` | Ya | any | Mulai pembayaran |
| `GET` | `/api/v1/bookings/:id/payment` | Ya | any | Cek status pembayaran |
| `POST` | `/api/v1/payments/:id/callback` | Tidak | - | Webhook dari gateway |

## Error Responses

| Kondisi | Status | Error |
|---------|--------|-------|
| Booking tidak ditemukan | 404 | `booking not found` |
| Bukan booking sendiri | 403 | `unauthorized` |
| Booking bukan pending | 400 | `booking cannot be paid` |
| Sudah dibayar | 400 | `already paid` |
| Payment tidak ditemukan | 404 | `payment not found` |
