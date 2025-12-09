# Business Trip Histories Feature

## Overview
Fitur `business_trip_histories` digunakan untuk mencatat semua perubahan penting pada business trip, khususnya fokus pada:
1. Perubahan status business trip
2. Proses verifikasi (approved, rejected, pending)

## Architecture

### Database Schema
Table: `business_trip_histories`

| Column | Type | Description |
|--------|------|-------------|
| id | UUID | Primary key |
| business_trip_id | UUID | Reference ke business_trips |
| change_type | VARCHAR(50) | Jenis perubahan (status_change, verification_approved, verification_rejected, verification_pending) |
| field_name | VARCHAR(100) | Nama field yang berubah |
| old_value | TEXT | Nilai lama |
| new_value | TEXT | Nilai baru |
| changed_by_user_id | VARCHAR(100) | User ID yang melakukan perubahan |
| changed_by_user_name | VARCHAR(255) | Nama user yang melakukan perubahan |
| notes | TEXT | Catatan tambahan |
| created_at | TIMESTAMP | Waktu perubahan dicatat |

### Domain Layer

#### Entity
- `BusinessTripHistory` (`internal/domain/entity/business_trip_history.go`)
  - Factory function: `NewBusinessTripHistory`
  - Setter methods: `SetFieldName`, `SetOldValue`, `SetNewValue`, `SetChangedBy`, `SetNotes`
  - Getter methods: untuk semua field

#### Repository Interface
- `BusinessTripHistoryRepository` (`internal/domain/repository/business_trip_history_repository.go`)
  - `Create()` - Membuat record history baru
  - `FindByBusinessTripID()` - Mendapatkan semua history untuk business trip tertentu
  - `FindByID()` - Mendapatkan history berdasarkan ID
  - `FindByBusinessTripIDAndChangeType()` - Filter history berdasarkan tipe perubahan

### Infrastructure Layer

#### PostgreSQL Implementation
- `businessTripHistoryRepository` (`internal/infrastructure/postgres/business_trip_history_repository.go`)
  - Implementasi dari `BusinessTripHistoryRepository` interface
  - Menggunakan sqlx untuk database operations

### Use Case Layer

#### Record History Use Case
- `RecordHistoryUseCase` (`internal/usecase/business_trip/record_history.go`)
  - Input: `RecordHistoryInput` struct
  - Fungsi: Mencatat perubahan ke database
  - Digunakan oleh use case lain untuk tracking changes

#### Get Histories Use Case
- `GetHistoriesUseCase` (`internal/usecase/business_trip/get_histories.go`)
  - Input: `businessTripID` string
  - Output: `[]BusinessTripHistoryResponse`
  - Fungsi: Mengambil semua history untuk business trip tertentu

### Integration

#### Integrated into Update Business Trip
File: `internal/usecase/business_trip/update_business_trip.go`
- Mencatat perubahan status secara otomatis
- History dicatat setelah status berubah

#### Integrated into Update Business Trip With Assignees
File: `internal/usecase/business_trip/update_business_trip_with_assignees.go`
- Mencatat perubahan status secara otomatis
- History dicatat setelah status berubah

#### Integrated into Create Business Trip
File: `internal/usecase/business_trip/create_business_trip.go`
- Mencatat status awal business trip saat pembuatan

#### Integrated into Verify Business Trip
File: `internal/usecase/business_trip/verify_business_trip.go`
- Mencatat setiap aksi verifikasi (approved, rejected)
- Termasuk user information dan notes dari verificator

### Delivery Layer

#### Handler
- `GetBusinessTripHistories()` di `BusinessTripHandler` (`internal/delivery/http/handler/business_trip_handler.go`)
  - Method: GET
  - Path: `/api/v1/business-trips/:tripId/histories`
  - Auth: Required (AuthMiddleware)
  - Response: List of history records

## API Endpoint

### Get Business Trip Histories
```
GET /api/v1/business-trips/:tripId/histories
Authorization: Bearer <token>
```

**Response:**
```json
{
  "message": "Business trip histories retrieved successfully",
  "data": [
    {
      "id": "uuid",
      "business_trip_id": "uuid",
      "change_type": "status_change",
      "field_name": "status",
      "old_value": "draft",
      "new_value": "ready_to_verify",
      "changed_by_user_id": "user-id",
      "changed_by_user_name": "John Doe",
      "notes": "",
      "created_at": "2025-12-03T15:00:00Z"
    },
    {
      "id": "uuid",
      "business_trip_id": "uuid",
      "change_type": "verification_approved",
      "field_name": "verificator",
      "old_value": "",
      "new_value": "approved",
      "changed_by_user_id": "user-id",
      "changed_by_user_name": "Jane Smith",
      "notes": "Approved after review",
      "created_at": "2025-12-03T16:00:00Z"
    }
  ]
}
```

## Migration

### Up Migration
File: `migrations/025_create_business_trip_histories.up.sql`
- Membuat tabel `business_trip_histories`
- Menambahkan indexes untuk performa
- Menambahkan constraint untuk change_type

### Down Migration
File: `migrations/025_create_business_trip_histories.down.sql`
- Drop indexes
- Drop table

## Usage Example

### Automatic Recording (Internal)

#### Status Change
```go
// Otomatis tercatat saat status berubah
updateRequest := UpdateBusinessTripRequest{
    BusinessTripID: "trip-id",
    Status: nullable.NullString{String: "ready_to_verify", Valid: true},
}

result, err := updateBusinessTripUseCase.Execute(ctx, updateRequest)
// History akan otomatis tercatat
```

#### Verification Action
```go
// Otomatis tercatat saat verifikasi
verifyRequest := VerifyBusinessTripRequest{
    BusinessTripID: "trip-id",
    VerificationStatus: "approved",
    VerificationNotes: "All documents are complete",
}

result, err := verifyBusinessTripUseCase.Execute(ctx, verifyRequest, authenticatedUser)
// History akan otomatis tercatat dengan user information
```

### Manual Recording (If Needed)
```go
// Bisa juga digunakan secara manual untuk tracking perubahan lain
err := recordHistoryUseCase.Execute(ctx, RecordHistoryInput{
    BusinessTripID: "trip-id",
    ChangeType: entity.HistoryChangeTypeStatusChange,
    FieldName: "status",
    OldValue: "draft",
    NewValue: "ready_to_verify",
    ChangedByUserID: "user-id",
    ChangedByUserName: "John Doe",
    Notes: "Updated via API",
})
```

## Change Types

1. **status_change** - Untuk tracking perubahan status business trip
2. **verification_approved** - Saat verificator approve
3. **verification_rejected** - Saat verificator reject
4. **verification_pending** - Saat verificator status kembali ke pending

## Future Enhancements

Fitur ini dapat dikembangkan lebih lanjut untuk mencatat:
- Perubahan assignee
- Perubahan transaction
- Perubahan document link
- Perubahan verificator list
- Dan perubahan field penting lainnya
