# Booking System API Documentation

## Overview
Sistem booking tiket yang telah diimplementasi dengan fitur-fitur:

1. **Create Booking** - User dapat membuat booking dengan multiple penumpang
2. **Upload Payment Proof** - User dapat mengunggah bukti pembayaran
3. **Payment Verification** - Staff dapat memverifikasi pembayaran
4. **Auto Expiration** - Booking otomatis expired setelah 30 menit
5. **Status Management** - Tracking status booking dari pending hingga success/rejected

## Booking Flow

### 1. Check Available Seats (Optional)
```
GET /api/protected/schedules/{schedule_id}/seats
Authorization: Bearer <jwt_token>
```

**Response:**
```json
{
  "code": "SUCCESS",
  "message": "Available seats retrieved successfully",
  "data": {
    "schedule_id": 1,
    "total_seats": 40,
    "available_seats": 36,
    "booked_seats": 4,
    "seats": [
      {
        "id": 1,
        "seat_number": "A1",
        "is_booked": false
      },
      {
        "id": 2,
        "seat_number": "A2",
        "is_booked": true
      }
    ]
  }
}
```

### 2. User Creates Booking
```
POST /api/protected/bookings
Content-Type: application/json
Authorization: Bearer <jwt_token>

{
  "schedule_id": 1,
  "passengers": [
    {
      "passenger_name": "Asep",
      "seat_id": 3
    },
    {
      "passenger_name": "Bambang", 
      "seat_id": 4
    }
  ]
}
```

**Response:**
```json
{
  "code": "SUCCESS",
  "message": "Booking created successfully",
  "data": {
    "id": 1,
    "status": "pending",
    "expires_at": "2025-06-14T10:30:00Z",
    "total_amount": 200000,
    "schedule": {
      "id": 1,
      "origin": "Jakarta",
      "destination": "Bandung",
      "departure_time": "2025-06-15 08:00",
      "arrival_time": "2025-06-15 11:00",
      "price": 100000
    },
    "passengers": [
      {
        "passenger_name": "Asep",
        "seat_number": "A3"
      },
      {
        "passenger_name": "Bambang",
        "seat_number": "A4"
      }
    ],
    "created_at": "2025-06-14T10:00:00Z"
  }
}
```

### 2. User Uploads Payment Proof
```
POST /api/protected/bookings/{id}/payment
Content-Type: multipart/form-data
Authorization: Bearer <jwt_token>

payment_method: "Bank Transfer"
proof_image: <file>
```

**Response:**
```json
{
  "code": "SUCCESS",
  "message": "Payment proof uploaded successfully"
}
```

### 3. Staff Verifies Payment
```
PUT /api/admin/bookings/{id}/status
Content-Type: application/json
Authorization: Bearer <admin_jwt_token>

{
  "status": "success", // or "rejected"
  "notes": "Payment verified successfully"
}
```

### Download Payment Proof (Staff Only)
```
GET /api/admin/bookings/{id}/payment/download
Authorization: Bearer <admin_jwt_token>
```

Parameters:
- `id`: Booking ID (required)

Response:
- **Success**: File download (binary)
- **Headers**: 
  - `Content-Disposition: attachment; filename=payment_proof_{booking_id}.jpg`
  - `Content-Type: application/octet-stream`

**Errors:**
- `404`: Payment proof not found or file missing on disk
- `403`: Not staff/admin
- `400`: Invalid booking ID

## API Endpoints

### User Endpoints (Authenticated)

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/protected/schedules/{id}/seats` | Get available seats for schedule |
| POST | `/api/protected/bookings` | Create new booking |
| GET | `/api/protected/bookings` | Get user's bookings |
| GET | `/api/protected/bookings/{id}` | Get booking details |
| POST | `/api/protected/bookings/{id}/payment` | Upload payment proof |

### Admin/Staff Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/admin/bookings` | Get all bookings |
| GET | `/api/admin/bookings/{id}` | Get booking details |
| GET | `/api/admin/bookings/{id}/payment/download` | Download payment proof |
| PUT | `/api/admin/bookings/{id}/status` | Update booking status |

## Booking Status Flow

```
pending → waiting_verification → success/rejected
   ↓                                    ↓
expired (after 30 minutes)        seats freed
```

### Status Descriptions:
- **pending**: Booking created, waiting for payment
- **waiting_verification**: Payment proof uploaded, waiting for staff verification
- **success**: Payment verified and confirmed
- **rejected**: Payment rejected by staff, **seats freed for rebooking**
- **expired**: Booking expired (30 minutes timeout), **seats freed for rebooking**
- **cancelled**: Booking cancelled by user/system

## Query Parameters

### Get Bookings
```
GET /api/protected/bookings?page=1&limit=10&status=pending,waiting_verification
```

Parameters:
- `page`: Page number (default: 1)
- `limit`: Items per page (default: 10, max: 100)
- `status`: Filter by status (comma-separated)

## Validation Rules

### Create Booking
- `schedule_id`: Required, must exist
- `passengers`: Required, min 1, max 10 passengers
- `passenger_name`: Required, 2-100 characters
- `seat_id`: Required, must be valid and available

### Upload Payment Proof
- `payment_method`: Required, 2-50 characters
- `proof_image`: Required, max 5MB, only JPEG/PNG

### Update Status (Staff)
- `status`: Required, only "success" or "rejected"
- `notes`: Optional, max 500 characters

## Error Handling

Common error responses:
```json
{
  "code": "ERROR",
  "message": "Error description",
  "data": "Detailed error information"
}
```

### Common Errors:
- `400 Bad Request`: Validation errors, invalid data
- `401 Unauthorized`: Missing or invalid JWT token
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource conflict (e.g., seat already booked)
- `500 Internal Server Error`: Server error

## Security Features

1. **JWT Authentication**: All endpoints require valid JWT token
2. **Role-Based Access**: Admin endpoints require admin/staff role
3. **User Isolation**: Users can only access their own bookings
4. **File Validation**: Payment proof files are validated for type and size
5. **Seat Locking**: Database transactions prevent double booking

## Seat Management

### Seat Booking Process
1. **Seat Validation**: Sistem mengecek apakah kursi sudah dibooking (`is_booked = true`)
2. **Seat Locking**: Ketika booking dibuat, kursi di-mark sebagai `is_booked = true`
3. **Seat Release**: Ketika booking expired/cancelled, kursi dikembalikan (`is_booked = false`)

### Database Transaction
- Menggunakan database transaction untuk mencegah race condition
- Seat status diupdate bersamaan dengan pembuatan booking
- Jika ada error, semua perubahan di-rollback

## Background Processes (Cron Job)

### Auto Expiration Purpose
Cron job berjalan setiap menit untuk:
1. **Mencari booking expired**: Booking dengan status `pending` yang sudah > 30 menit
2. **Update status**: Mengubah status dari `pending` ke `expired`  
3. **Free seats**: Mengubah `is_booked = false` untuk kursi yang terkait
4. **Enable rebooking**: Kursi bisa dibooking lagi oleh user lain

### Why This is Important
- **Prevent seat hoarding**: User tidak bisa "menahan" kursi tanpa bayar
- **Optimize availability**: Kursi yang tidak dibayar dikembalikan ke pool
- **Fair booking**: User lain mendapat kesempatan booking kursi yang sama

## File Upload System

### File Storage
- **Directory**: `uploads/payments/` (otomatis dibuat)
- **Filename Format**: `payment_{booking_id}_{timestamp}.{extension}`
- **File Types**: JPEG, JPG, PNG only
- **Max Size**: 5MB
- **Storage**: Local filesystem (bisa dipindah ke cloud storage)
