# Snake Case Field Mapping

## Standardization Summary

All API request and response fields have been standardized to use snake_case format for consistency.

## Field Mapping Changes

### Business Trip
| Before (camelCase) | After (snake_case) |
|-------------------|-------------------|
| `startDate` | `start_date` |
| `endDate` | `end_date` |
| `activityPurpose` | `activity_purpose` |
| `destinationCity` | `destination_city` |
| `spdDate` | `spd_date` |
| `departureDate` | `departure_date` |
| `returnDate` | `return_date` |
| `totalCost` | `total_cost` |
| `createdAt` | `created_at` |
| `updatedAt` | `updated_at` |

### Assignee
| Before (camelCase) | After (snake_case) |
|-------------------|-------------------|
| `spdNumber` | `spd_number` |
| `employeeId` | `employee_id` |
| `employeeName` | `employee_name` |
| `totalCost` | `total_cost` |
| `createdAt` | `created_at` |
| `updatedAt` | `updated_at` |

### Transaction
| Before (camelCase) | After (snake_case) |
|-------------------|-------------------|
| `totalNight` | `total_night` |
| `subtotal` | `subtotal` (unchanged) |
| `description` | `description` (unchanged) |
| `transportDetail` | `transport_detail` |
| `createdAt` | `created_at` |
| `updatedAt` | `updated_at` |

### Response Collections
| Before (camelCase) | After (snake_case) |
|-------------------|-------------------|
| `businessTrips` | `business_trips` |
| `totalPages` | `total_pages` |

### Summary Objects
| Before (camelCase) | After (snake_case) |
|-------------------|-------------------|
| `businessTripId` | `business_trip_id` |
| `totalCost` | `total_cost` |
| `totalAssignees` | `total_assignees` |
| `totalTransactions` | `total_transactions` |
| `costByType` | `cost_by_type` |
| `assigneeId` | `assignee_id` |
| `assigneeName` | `assignee_name` |

## API Endpoints Affected

All endpoints now use consistent snake_case format:

- `POST /api/v1/business-trips`
- `GET /api/v1/business-trips`
- `GET /api/v1/business-trips/{id}`
- `PUT /api/v1/business-trips/{id}`
- `PUT /api/v1/business-trips/{id}/with-assignees`
- `DELETE /api/v1/business-trips/{id}`
- All assignee and transaction endpoints

## Example Request/Response

### Request (snake_case)
```json
{
  "start_date": "2024-01-15",
  "end_date": "2024-01-17",
  "activity_purpose": "Meeting with client",
  "destination_city": "Jakarta",
  "spd_date": "2024-01-10",
  "departure_date": "2024-01-15",
  "return_date": "2024-01-17",
  "assignees": [
    {
      "name": "John Doe",
      "spd_number": "SPD-2024-001",
      "employee_id": "123456789",
      "employee_name": "John Doe",
      "position": "Software Engineer",
      "rank": "III/a",
      "transactions": [
        {
          "name": "Hotel Booking",
          "type": "accommodation",
          "subtype": "hotel",
          "amount": 500000,
          "total_night": 2,
          "description": "Hotel di Jakarta Selatan"
        }
      ]
    }
  ]
}
```

### Response (snake_case)
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "start_date": "2024-01-15",
  "end_date": "2024-01-17",
  "activity_purpose": "Meeting with client",
  "destination_city": "Jakarta",
  "spd_date": "2024-01-10",
  "departure_date": "2024-01-15",
  "return_date": "2024-01-17",
  "total_cost": 2800000,
  "assignees": [
    {
      "id": "6d362476-cfc0-43c9-bc8a-2d0f1e3942f3",
      "name": "John Doe",
      "spd_number": "SPD-2024-001",
      "employee_id": "123456789",
      "employee_name": "John Doe",
      "position": "Software Engineer",
      "rank": "III/a",
      "total_cost": 2800000,
      "transactions": [...],
      "created_at": "2024-11-17T08:27:59Z",
      "updated_at": "2024-11-17T08:27:59Z"
    }
  ],
  "created_at": "2024-11-17T08:27:59Z",
  "updated_at": "2024-11-17T08:27:59Z"
}
```

## Benefits

1. **Consistency**: All fields follow the same naming convention
2. **Database Alignment**: JSON field names match database column names
3. **API Standards**: Follows REST API best practices
4. **Reduced Confusion**: No more camelCase vs snake_case mismatch
5. **Better Readability**: Snake_case is more readable for longer field names